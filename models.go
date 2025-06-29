package gname

import (
	"strconv"
	"time"

	"github.com/libdns/libdns"
)

type CommonResponse struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type ResolutionList struct {
	Code     int                      `json:"code,omitempty"`
	Msg      string                   `json:"msg,omitempty"`
	Data     []DomainResolutionRecord `json:"data,omitempty"`
	Count    int                      `json:"count,omitempty"`
	Page     int                      `json:"page,omitempty"`
	PageSize int                      `json:"pagesize,omitempty"`
}

type DomainResolutionRecord struct {
	ID   string `json:"id,omitempty"`
	Ym   string `json:"ym,omitempty"`
	Zjt  string `json:"zjt,omitempty"`
	Lx   string `json:"lx,omitempty"`
	Jxz  string `json:"jxz,omitempty"`
	Mx   string `json:"mx,omitempty"`
	Xlid int    `json:"xlid,omitempty"`
	Zt   string `json:"zt,omitempty"`
	TTL  string `json:"ttl,omitempty"`
}

type AddDomainRecord struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data int    `json:"data,omitempty"`
}

type UpdateDomainRecord struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data,omitempty"`
}

type DeleteDomainRecord struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (record DomainResolutionRecord) toLibdnsRecord(zone string) libdns.RR {
	// Parse TTL from API response, default to 600 seconds if not provided or invalid
	ttl := time.Second * 600
	if record.TTL != "" {
		if parsedTTL, err := strconv.Atoi(record.TTL); err == nil && parsedTTL > 0 {
			ttl = time.Duration(parsedTTL) * time.Second
		}
	}

	return libdns.RR{
		Name: record.Zjt,
		Type: record.Lx,
		Data: record.Jxz,
		TTL:  ttl,
	}
}

// RR implements the libdns.Record interface
func (record DomainResolutionRecord) RR() libdns.RR {
	return record.toLibdnsRecord("")
}
