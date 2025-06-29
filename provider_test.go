package gname

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/libdns/libdns"
)

func TestProvider_AppendRecords(t *testing.T) {
	type fields struct {
		APPID  string
		APPKey string
	}
	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []libdns.Record
		wantErr bool
	}{
		{
			name: "Test AppendRecords",
			fields: fields{
				APPID:  "Your_APPID",
				APPKey: "Your_APPKEY",
			},
			args: args{
				ctx:  context.Background(),
				zone: "388vip.com",
				records: []libdns.Record{
					libdns.RR{
						Name: "jump-test",
						Type: "A",
						Data: "8.8.8.8",
						TTL:  time.Second * 120,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				APPID:  tt.fields.APPID,
				APPKey: tt.fields.APPKey,
			}
			got, err := p.AppendRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_DeleteRecords(t *testing.T) {
	type fields struct {
		APPID  string
		APPKey string
	}
	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []libdns.Record
		wantErr bool
	}{
		{
			name: "Test DeleteRecords",
			fields: fields{
				APPID:  "Your_APPID",
				APPKey: "Your_APPKEY",
			},
			args: args{
				ctx:  context.Background(),
				zone: "388vip.com",
				records: []libdns.Record{
					libdns.RR{
						Name: "jump-test",
						Type: "A",
						Data: "8.8.8.8",
						TTL:  time.Second * 120,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				APPID:  tt.fields.APPID,
				APPKey: tt.fields.APPKey,
			}
			got, err := p.DeleteRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_GetRecords(t *testing.T) {
	type fields struct {
		APPID  string
		APPKey string
	}
	type args struct {
		ctx  context.Context
		zone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []libdns.Record
		wantErr bool
	}{
		{
			name: "Test GetRecords",
			fields: fields{
				APPID:  "Your_APPID",
				APPKey: "Your_APPKEY",
			},
			args: args{
				ctx:  context.Background(),
				zone: "388vip.com",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				APPID:  tt.fields.APPID,
				APPKey: tt.fields.APPKey,
			}
			got, err := p.GetRecords(tt.args.ctx, tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_SetRecords(t *testing.T) {
	type fields struct {
		APPID  string
		APPKey string
	}
	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []libdns.Record
		wantErr bool
	}{
		{
			name: "Test SetRecords",
			fields: fields{
				APPID:  "Your_APPID",
				APPKey: "Your_APPKEY",
			},
			args: args{
				ctx:  context.Background(),
				zone: "388vip.com",
				records: []libdns.Record{
					libdns.RR{
						Name: "jump-test",
						Type: "A",
						Data: "8.8.8.8",
						TTL:  time.Second * 120,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				APPID:  tt.fields.APPID,
				APPKey: tt.fields.APPKey,
			}
			got, err := p.SetRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}
