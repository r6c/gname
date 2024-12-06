package gname

import (
	"context"
	"fmt"
	"github.com/libdns/libdns"
	"strconv"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with <TODO: PROVIDER NAME>.
type Provider struct {
	APPID  string `json:"app_id,omitempty"`
	APPKey string `json:"app_key,omitempty"`
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	params := fmt.Sprintf("appid=%s&ym=%s", p.APPID, trimmedZone)

	response, err := MakeApiRequest("POST", "/api/resolution/list", params, p.APPKey, ResolutionList{})
	if err != nil {
		return nil, err
	}

	recs := make([]libdns.Record, 0, len(response.Data))
	for _, rec := range response.Data {
		recs = append(recs, rec.toLibdnsRecord(trimmedZone))
	}
	return recs, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var successfullyAppendedRecords []libdns.Record
	trimmedZone := libdnsZoneToDnslaDomain(zone)

	for _, record := range records {
		params := fmt.Sprintf("appid=%s&ym=%s&lx=%s&zj=%s&jlz=%s&mx=%d&ttl=%.0f",
			p.APPID, trimmedZone, record.Type, record.Name, record.Value, record.Weight, record.TTL.Seconds())

		response, err := MakeApiRequest("POST", "/api/resolution/add", params, p.APPKey, AddDomainRecord{})
		if err != nil {
			return []libdns.Record{}, err
		}

		appendedRecord := libdns.Record{
			ID:     strconv.Itoa(response.Data),
			Name:   record.Name,
			Type:   record.Type,
			Value:  record.Value,
			Weight: record.Weight,
			TTL:    record.TTL,
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
		return successfullyUpdatedRecords, err
	}

	for _, record := range records {
		hasRecord := false
		recordId := 0
		for _, rec := range recs {
			if rec.Name == record.Name && rec.Type == record.Type {
				hasRecord = true
				recordId, err = strconv.Atoi(rec.ID)
				if err != nil {
					return successfullyUpdatedRecords, err
				}
				break
			}
		}

		if !hasRecord {
			appendedRecords, err := p.AppendRecords(ctx, zone, []libdns.Record{record})
			if err != nil {
				return successfullyUpdatedRecords, err
			}

			successfullyUpdatedRecords = append(successfullyUpdatedRecords, appendedRecords...)
			continue
		}

		params := fmt.Sprintf("appid=%s&ym=%s&lx=%s&zj=%s&jlz=%s&mx=%d&ttl=%.0f&jxid=%d",
			p.APPID, trimmedZone, record.Type, record.Name, record.Value, record.Weight, record.TTL.Seconds(), recordId)

		response, err := MakeApiRequest("POST", "/api/resolution/edit", params, p.APPKey, UpdateDomainRecord{})
		if err != nil {
			return successfullyUpdatedRecords, err
		}

		if response.Code == 1 {
			successfullyUpdatedRecords = append(successfullyUpdatedRecords, libdns.Record{
				ID:     strconv.Itoa(recordId),
				Name:   record.Name,
				Type:   record.Type,
				Value:  record.Value,
				Weight: record.Weight,
				TTL:    record.TTL,
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
		return successfullyDeletedRecords, err
	}

	for _, record := range records {
		for _, rec := range recs {
			if rec.Name == record.Name && rec.Type == record.Type {
				record.ID = rec.ID
				record.Value = rec.Value
				break
			}
		}

		params := fmt.Sprintf("appid=%s&ym=%s&jxid=%s", p.APPID, trimmedZone, record.ID)

		response, err := MakeApiRequest("POST", "/api/resolution/delete", params, p.APPKey, DeleteDomainRecord{})
		if err != nil {
			return successfullyDeletedRecords, err
		}
		if response.Code == 1 {
			successfullyDeletedRecords = append(successfullyDeletedRecords, record)
		}
	}

	return successfullyDeletedRecords, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
