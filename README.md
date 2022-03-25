libdns-glesys for [`libdns`](https://github.com/libdns/libdns)
==============================================================

This package implements the [libdns interfaces](https://github.com/libdns/libdns) for [Glesys](https://glesys.se), allowing you to manage DNS records.
It utilizes [glesys-go](https://github.com/glesys/glesys-go) for API communication.

## Usage
```golang
include (
    glesys "github.com/libdns/glesys"
)
p := &glesys.Provider{
    Project: "your project/username usually clXXXXXX",
    ApiKey: "api-key",
}

zone := "example.org"
records, err := p.GetRecords(ctx, zone)
```
For more examples check the `_examples` folder in the source.

## Noteworthy
To do everything this library can do the Glesys API user needs permissions to the following...

- Domain.addrecord
- Domain.deleterecord
- Domain.listrecords
- Domain.updaterecord
