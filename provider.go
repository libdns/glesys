package glesys

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/glesys/glesys-go/v3"
	"github.com/libdns/libdns"
)

type Provider struct {
	mutex       sync.Mutex
	clientCache *glesys.Client
	Project     string `json:"project,omitempty"`
	ApiKey      string `json:"api_key,omitempty"`
}

func (p *Provider) client() *glesys.Client {
	if p.clientCache == nil {
		p.clientCache = glesys.NewClient(p.Project, p.ApiKey, "libdns-glesys/0.0.1")

	}
	return p.clientCache
}
func cleanZ(z string) string {
	return strings.TrimRight(z, ". ")
}
func gle2lib(dr *glesys.DNSDomainRecord) libdns.Record {
	r := libdns.Record{
		ID:    strconv.Itoa(dr.RecordID),
		Type:  dr.Type,
		Name:  dr.Host,
		Value: dr.Data,
		TTL:   time.Duration(dr.TTL) * time.Second,
	}
	switch dr.Type {
	// extract priority
	case "MX", "SRV", "URI":
		parts := strings.Split(dr.Data, " ")
		if p, err := strconv.Atoi(parts[0]); err == nil {
			r.Priority = p
			r.Value = parts[1]
		}
	}
	return r
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	zone = cleanZ(zone)
	log.Printf("GetRecords zone=%s", zone)
	drs, err := p.client().DNSDomains.ListRecords(ctx, zone)
	if err != nil {
		return nil, err
	}
	records := make([]libdns.Record, len(*drs))
	for i, dr := range *drs {
		if zone != dr.DomainName {
			return records, fmt.Errorf("unexpected domainname in respose: %v", dr.DomainName)
		}
		r := gle2lib(&dr)
		records[i] = r
	}
	return records, nil

}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	zone = cleanZ(zone)
	log.Printf("AppendRecords zone=%s", zone)
	results := []libdns.Record{}
	for _, r := range records {
		param := glesys.AddRecordParams{
			DomainName: zone,
			Host:       r.Name,
			Data:       r.Value,
			TTL:        int(r.TTL / time.Second),
			Type:       strings.ToUpper(r.Type),
		}
		if r.Priority > 0 {
			param.Data = fmt.Sprintf("%d %s", r.Priority, param.Data)
		}

		dr, err := p.client().DNSDomains.AddRecord(ctx, param)
		if err != nil {
			return results, err
		}
		results = append(results, gle2lib(dr))
	}
	return results, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	zone = cleanZ(zone)
	log.Printf("SetRecords zone=%s", zone)
	results := []libdns.Record{}
	for _, r := range records {
		id, err := strconv.Atoi(r.ID)
		if err != nil && r.ID != "" {
			return results, err
		}
		param := glesys.UpdateRecordParams{
			RecordID: id,
			Host:     r.Name,
			Data:     r.Value,
			TTL:      int(r.TTL / time.Second),
			Type:     strings.ToUpper(r.Type),
		}
		if r.Priority > 0 {
			param.Data = fmt.Sprintf("%d %s", r.Priority, param.Data)
		}

		dr, err := p.client().DNSDomains.UpdateRecord(ctx, param)
		if err != nil {
			return results, err
		}
		results = append(results, gle2lib(dr))
	}
	return results, nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	zone = cleanZ(zone)
	log.Printf("DeleteRecords zone=%s", zone)
	results := []libdns.Record{}
	for _, r := range records {
		if r.ID == "" {
			return results, fmt.Errorf("record must have ID to delete: %+v", r)
		}
		id, err := strconv.Atoi(r.ID)
		if err != nil {
			return results, err
		}

		err = p.client().DNSDomains.DeleteRecord(ctx, id)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}
	return results, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
