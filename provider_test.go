// SPDX-FileCopyrightText: 2022 Peter Magnusson <me@kmpm.se>
// SPDX-License-Identifier: MIT

package glesys

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/libdns/glesys/internal/impl"
	"github.com/libdns/libdns"
)

func TestProvider_GetRecordsIntegration(t *testing.T) {
	skipUnauth(t)
	type fields struct {
		Project string
		APIKey  string
	}
	type args struct {
		ctx  context.Context
		zone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"first_test",
			fields{os.Getenv("GLESYS_PROJECT"), os.Getenv("GLESYS_KEY")},
			args{context.TODO(), os.Getenv("GLESYS_ZONE")},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				Project: tt.fields.Project,
				APIKey:  tt.fields.APIKey,
			}
			got, err := p.GetRecords(tt.args.ctx, tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < 1 {
				t.Errorf("Expected > 0 records. Got %v", len(got))
			}
			for i, r := range got {
				rr := r.RR()
				if rr.Name == "" {
					t.Errorf("Expected record %d name. Got '%v'", i, rr.Name)
				}
				if rr.TTL < 1*time.Second {
					t.Errorf("Record %d (%v) hand < 1 second TTL. Got %v", i, r, rr.TTL)
				}
				switch rr.Type {
				case "MX":
					x, err := rr.Parse()
					if err != nil {
						t.Errorf("Failed to parse record %d: %v", i, err)
					}
					if mx, ok := x.(libdns.MX); !ok {
						t.Errorf("Expected record %d to be of type MX. Got %T", i, r)
					} else if mx.Preference < 1 {
						t.Errorf("Expected record %d of type %s to have a priority. Got %v", i, rr.Type, mx.Preference)
					}

				}
			}
			// t.Errorf("Got %+v", got)
		})
	}
}

func mustRRParse(t *testing.T, r libdns.RR) libdns.Record {
	rr, err := r.Parse()
	if err != nil {
		t.Fatalf("Failed to parse RR: %v", err)
	}
	return rr
}

func TestProvider_AppendAndDeleteRecordsIntegration(t *testing.T) {
	skipUnauth(t)
	type fields struct {
		Project string
		APIKey  string
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
		wantErr bool
	}{
		{"create_txt_record",
			fields{os.Getenv("GLESYS_PROJECT"), os.Getenv("GLESYS_KEY")},
			args{
				context.TODO(),
				os.Getenv("GLESYS_ZONE"),
				[]libdns.Record{
					mustRRParse(t, libdns.RR{Type: "TXT", TTL: time.Minute * 5, Name: "_libdns-test", Data: time.Now().Local().String()}),
				},
			},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				Project: tt.fields.Project,
				APIKey:  tt.fields.APIKey,
			}
			got, err := p.AppendRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.AppendRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.args.records) {
				t.Errorf("Record count requested is not same as returned. Got %v, want %v", len(got), len(tt.args.records))
			}
			for i, gr := range got {
				wrr := tt.args.records[i].RR()
				grr := gr.RR()
				// if gr.ID == "" {
				// 	t.Errorf("Expected a record ID but got nothing: %+v", gr)
				// }

				if grr.Name != wrr.Name {
					t.Errorf("Mismatching name. Got %v, Want %v", grr.Name, wrr.Name)
				}
				if grr.TTL != wrr.TTL {
					t.Errorf("Mismatching TTL. Got %v, Want %v", grr.TTL, wrr.TTL)
				}
				if grr.Data != wrr.Data {
					t.Errorf("Mismatching Value. Got %v, Want %v", grr.Data, wrr.Data)
				}
			}

			// Delete the same
			gotDeleted, err := p.DeleteRecords(context.TODO(), tt.args.zone, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.DeleteRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotDeleted) != len(tt.args.records) {
				t.Errorf("Record count requested is not same as returned. Got %v, want %v", len(gotDeleted), len(tt.args.records))
			}
		})
	}
}

func TestProvider_getMatchingRecords(t *testing.T) {
	skipUnauth(t)
	type fields struct {
		Project string
		APIKey  string
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
		want    []recordWithMatchingGlesys
		wantErr bool
	}{
		{"getting_matching_records",
			fields{os.Getenv("GLESYS_PROJECT"), os.Getenv("GLESYS_KEY")},
			args{
				context.TODO(),
				os.Getenv("GLESYS_ZONE"),
				[]libdns.Record{
					mustRRParse(t, libdns.RR{Type: "TXT", Name: "_libdns-test"}),
				},
			},
			[]recordWithMatchingGlesys{
				{
					Record:  mustRRParse(t, libdns.RR{Type: "TXT", Name: "_libdns-test"}),
					Matches: []impl.DNSDomainRecord{},
				}},
			false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				Project: tt.fields.Project,
				APIKey:  tt.fields.APIKey,
			}
			got, err := p.getMatchingRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.getMatchingRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.getMatchingRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
