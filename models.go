package gname

import (
	"github.com/libdns/libdns"
	"strconv"
	"time"
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

func (record DomainResolutionRecord) toLibdnsRecord(zone string) libdns.Record {
	mx, err := strconv.Atoi(record.Mx)
	if err != nil {
		mx = 0
	}

	return libdns.Record{
		ID:       record.ID,
		Name:     record.Zjt,
		Priority: uint(mx),
		TTL:      time.Second * 120,
		Type:     record.Lx,
		Value:    record.Jxz,
	}
}
