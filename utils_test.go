package glesys

import (
	"reflect"
	"testing"

	"github.com/libdns/glesys/internal/impl"
	"github.com/libdns/libdns"
)

func Test_checkParamsMatching(t *testing.T) {
	type args struct {
		rr libdns.RR
		dr *impl.DNSDomainRecord
	}
	tests := []struct {
		name string
		args args
		want matchParams
	}{
		{"all", args{libdns.RR{Name: "test", Type: "A", Data: "1.1.1.1"}, &impl.DNSDomainRecord{Host: "test", Type: "A", Data: "1.1.1.1"}},
			matchParams{Name: true, Type: true, Data: true, TTL: true}},
		{"zero_name", args{libdns.RR{Name: "", Type: "A", Data: "1.1.1.1"}, &impl.DNSDomainRecord{Host: "test", Type: "A", Data: "1.1.1.1"}},
			matchParams{Name: true, Type: true, Data: true, TTL: true}},
		{"zero_type", args{libdns.RR{Name: "test", Data: "1.1.1.1"}, &impl.DNSDomainRecord{Host: "test", Type: "A", Data: "1.1.1.1"}},
			matchParams{Name: true, Type: true, Data: true, TTL: true}},
		{"diff_type", args{libdns.RR{Name: "test", Type: "CNAME", Data: "1.1.1.1"}, &impl.DNSDomainRecord{Host: "test", Type: "A", Data: "1.1.1.1"}},
			matchParams{Name: true, Type: false, Data: true, TTL: true}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkParamsMatching(tt.args.rr, tt.args.dr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkParamsMatching() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanZ(t *testing.T) {
	type args struct {
		z string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with_dot", args{"test."}, "test"},
		{"with_space", args{"test "}, "test"},
		{"with_dot_space", args{"test. "}, "test"},
		{"without_any", args{"test"}, "test"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanZ(tt.args.z); got != tt.want {
				t.Errorf("cleanZ() = %v, want %v", got, tt.want)
			}
		})
	}
}
