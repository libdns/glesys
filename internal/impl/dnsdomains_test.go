// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2017 GleSYS Internet Services AB
package impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsDomainsAddRecord(t *testing.T) {
	c := &mockClient{body: `{"response": { "record":
          {"recordid": 1234569, "domainname": "example.com", "host": "test", "type": "A", "data": "127.0.0.1", "ttl": 3600}
	}}`}

	params := AddRecordParams{
		DomainName: "example.com",
		Host:       "test",
		Data:       "127.0.0.1",
		Type:       "A",
	}

	d := DNSDomainService{client: c}

	record, _ := d.AddRecord(context.Background(), params)

	assert.Equal(t, "POST", c.lastMethod, "method used is correct")
	assert.Equal(t, "domain/addrecord", c.lastPath, "path used is correct")
	assert.Equal(t, "test", (*record).Host, "Record host is correct")
	assert.Equal(t, "127.0.0.1", (*record).Data, "Record data is correct")
}

func TestDnsDomainsListRecords(t *testing.T) {
	c := &mockClient{body: `{"response": { "records": [
	  {"recordid": 1234567, "domainname": "example.com", "host": "www", "type": "A", "data": "127.0.0.1", "ttl": 3600},
          {"recordid": 1234568, "domainname": "example.com", "host": "mail", "type": "A", "data": "127.0.0.3", "ttl": 3600}
	]}}`}

	d := DNSDomainService{client: c}

	records, _ := d.ListRecords(context.Background(), "example.com")

	assert.Equal(t, "POST", c.lastMethod, "method used is correct")
	assert.Equal(t, "domain/listrecords", c.lastPath, "path used is correct")
	assert.Equal(t, "www", (*records)[0].Host, "Record host is correct")
	assert.Equal(t, "127.0.0.3", (*records)[1].Data, "Record data is correct")
}

func TestDnsDomainsUpdateRecord(t *testing.T) {
	c := &mockClient{body: `{"response": { "record":
          {"recordid": 1234567, "domainname": "example.com", "host": "mail", "type": "A", "data": "127.0.0.3", "ttl": 3600}
	}}`}

	params := UpdateRecordParams{
		RecordID: 1234567,
		Data:     "127.0.0.3",
	}

	d := DNSDomainService{client: c}

	record, _ := d.UpdateRecord(context.Background(), params)

	assert.Equal(t, "POST", c.lastMethod, "method used is correct")
	assert.Equal(t, "domain/updaterecord", c.lastPath, "path used is correct")
	assert.Equal(t, "mail", (*record).Host, "Record host is correct")
	assert.Equal(t, "127.0.0.3", (*record).Data, "Record data is correct")
}

func TestDnsDomainsDeleteRecord(t *testing.T) {
	c := &mockClient{}
	d := DNSDomainService{client: c}

	d.DeleteRecord(context.Background(), 1234567)

	assert.Equal(t, "POST", c.lastMethod, "method is used correct")
	assert.Equal(t, "domain/deleterecord", c.lastPath, "path used is correct")
}
