// SPDX-FileCopyrightText: 2022 Peter Magnusson <me@kmpm.se>
// SPDX-License-Identifier: MIT
package glesys

import (
	"context"
	"os"
	"testing"
	"time"

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
				if r.Name == "" {
					t.Errorf("Expected record %d name. Got '%v'", i, r.Name)
				}
				if r.TTL < 1*time.Second {
					t.Errorf("Record %d (%v) hand < 1 second TTL. Got %v", i, r, r.TTL)
				}
				switch r.Type {
				case "MX":
					if r.Priority < 1 {
						t.Errorf("Expected record %d of type %s to have a priority. Got %v", i, r.Type, r.Priority)
					}
				}
			}
			// t.Errorf("Got %+v", got)
		})
	}
}

func TestProvider_AppendRecordsIntegration(t *testing.T) {
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
					{Type: "TXT", TTL: time.Minute * 5, Name: "_libdns-test", Value: time.Now().Local().String()},
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
				if gr.ID == "" {
					t.Errorf("Expected a record ID but got nothing: %+v", gr)
				}
				wr := tt.args.records[i]
				if gr.Name != wr.Name {
					t.Errorf("Mismatching name. Got %v, Want %v", gr.Name, wr.Name)
				}
				if gr.TTL != wr.TTL {
					t.Errorf("Mismatching TTL. Got %v, Want %v", gr.TTL, wr.TTL)
				}
				if gr.Value != wr.Value {
					t.Errorf("Mismatching Value. Got %v, Want %v", gr.Value, wr.Value)
				}
			}
		})
	}
}
