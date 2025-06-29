package gname

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/libdns/libdns"
)

var (
	// Ensure Provider implements the libdns interfaces
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)

// Provider facilitates DNS record manipulation with GNAME.
// It implements libdns interfaces for managing DNS records.
type Provider struct {
	// APPID is the application ID for GNAME API authentication.
	// This is required for all API operations.
	APPID string `json:"app_id,omitempty"`

	// APPKey is the application key for GNAME API authentication.
	// This is required for all API operations and should be kept secure.
	APPKey string `json:"app_key,omitempty"`

	// HTTPClient allows you to specify a custom HTTP client for API requests.
	// If not specified, a sensible default will be used with appropriate timeouts.
	HTTPClient *http.Client `json:"-"`

	// mutex for protecting the initialization of the HTTP client
	mu sync.RWMutex
}

// getHTTPClient returns the HTTP client, initializing it if necessary.
func (p *Provider) getHTTPClient() *http.Client {
	p.mu.RLock()
	if p.HTTPClient != nil {
		defer p.mu.RUnlock()
		return p.HTTPClient
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if p.HTTPClient == nil {
		p.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return p.HTTPClient
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	params := fmt.Sprintf("appid=%s&ym=%s", p.APPID, trimmedZone)

	response, err := MakeApiRequestWithClient(p.getHTTPClient(), "POST", "/api/resolution/list", params, p.APPKey, ResolutionList{})
	if err != nil {
		return nil, fmt.Errorf("failed to get records for zone %s: %w", zone, err)
	}

	recs := make([]libdns.Record, 0, len(response.Data))
	for _, rec := range response.Data {
		rr := rec.toLibdnsRecord(trimmedZone)
		recs = append(recs, rr)
	}
	return recs, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var successfullyAppendedRecords []libdns.Record
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	for _, record := range records {
		// Convert record to RR to access its fields
		rr := record.RR()

		params := fmt.Sprintf("appid=%s&ym=%s&lx=%s&zj=%s&jlz=%s&ttl=%.0f",
			p.APPID, trimmedZone, rr.Type, rr.Name, rr.Data, rr.TTL.Seconds())

		_, err := MakeApiRequestWithClient(p.getHTTPClient(), "POST", "/api/resolution/add", params, p.APPKey, AddDomainRecord{})
		if err != nil {
			return successfullyAppendedRecords, fmt.Errorf("failed to append record %s.%s: %w", rr.Name, zone, err)
		}

		appendedRecord := libdns.RR{
			Name: rr.Name,
			Type: rr.Type,
			Data: rr.Data,
			TTL:  rr.TTL,
		}

		successfullyAppendedRecords = append(successfullyAppendedRecords, appendedRecord)
	}

	return successfullyAppendedRecords, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var successfullyUpdatedRecords []libdns.Record
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	recs, err := p.GetRecords(ctx, zone)
	if err != nil {
		return successfullyUpdatedRecords, fmt.Errorf("failed to get existing records: %w", err)
	}

	for _, record := range records {
		rr := record.RR()
		hasRecord := false
		recordId := ""

		for _, rec := range recs {
			recRR := rec.RR()
			if recRR.Name == rr.Name && recRR.Type == rr.Type {
				hasRecord = true
				// Try to extract ID from the record if it's our custom type
				if dnslaRec, ok := rec.(DomainResolutionRecord); ok {
					recordId = dnslaRec.ID
				}
				break
			}
		}

		if !hasRecord {
			appendedRecords, err := p.AppendRecords(ctx, zone, []libdns.Record{record})
			if err != nil {
				return successfullyUpdatedRecords, fmt.Errorf("failed to create new record: %w", err)
			}

			successfullyUpdatedRecords = append(successfullyUpdatedRecords, appendedRecords...)
			continue
		}

		if recordId == "" {
			// Skip if we can't get the record ID
			continue
		}

		params := fmt.Sprintf("appid=%s&ym=%s&lx=%s&zj=%s&jlz=%s&ttl=%.0f&jxid=%s",
			p.APPID, trimmedZone, rr.Type, rr.Name, rr.Data, rr.TTL.Seconds(), recordId)

		response, err := MakeApiRequestWithClient(p.getHTTPClient(), "POST", "/api/resolution/edit", params, p.APPKey, UpdateDomainRecord{})
		if err != nil {
			return successfullyUpdatedRecords, fmt.Errorf("failed to update record %s.%s: %w", rr.Name, zone, err)
		}

		if response.Code == 1 {
			successfullyUpdatedRecords = append(successfullyUpdatedRecords, libdns.RR{
				Name: rr.Name,
				Type: rr.Type,
				Data: rr.Data,
				TTL:  rr.TTL,
			})
		}
	}

	return successfullyUpdatedRecords, nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var successfullyDeletedRecords []libdns.Record
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	recs, err := p.GetRecords(ctx, zone)
	if err != nil {
		return successfullyDeletedRecords, fmt.Errorf("failed to get existing records: %w", err)
	}

	for _, record := range records {
		rr := record.RR()
		recordId := ""

		for _, rec := range recs {
			recRR := rec.RR()
			if recRR.Name == rr.Name && recRR.Type == rr.Type {
				// Try to extract ID from the record if it's our custom type
				if dnslaRec, ok := rec.(DomainResolutionRecord); ok {
					recordId = dnslaRec.ID
				}
				break
			}
		}

		if recordId == "" {
			// Skip if we can't find the record ID
			continue
		}

		params := fmt.Sprintf("appid=%s&ym=%s&jxid=%s", p.APPID, trimmedZone, recordId)

		response, err := MakeApiRequestWithClient(p.getHTTPClient(), "POST", "/api/resolution/delete", params, p.APPKey, DeleteDomainRecord{})
		if err != nil {
			return successfullyDeletedRecords, fmt.Errorf("failed to delete record %s.%s: %w", rr.Name, zone, err)
		}
		if response.Code == 1 {
			successfullyDeletedRecords = append(successfullyDeletedRecords, rr)
		}
	}

	return successfullyDeletedRecords, nil
}
