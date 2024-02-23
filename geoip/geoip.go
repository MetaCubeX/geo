// Package geoip implements functions for parsing and processing GeoIP resource files.
package geoip

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	"github.com/sagernet/sing/common"
	F "github.com/sagernet/sing/common/format"
)

type Database struct {
	reader     *maxminddb.Reader
	SourceType Type
	MemoryType Type
}

func (db Database) Reader() *maxminddb.Reader {
	return db.reader
}

func (db Database) LookupCode(ip net.IP) []string {
	switch db.MemoryType {
	case TypeMaxmind:
		var country geoip2Country
		_ = db.reader.Lookup(ip, &country)
		if country.Country.IsoCode == "" {
			return nil
		}
		return []string{country.Country.IsoCode}

	case TypeSing:
		var code string
		_ = db.reader.Lookup(ip, &code)
		if code == "" {
			return nil
		}
		return []string{code}

	case TypeMetaV0:
		var codes any
		_ = db.reader.Lookup(ip, &codes)
		switch codes := codes.(type) {
		case string:
			return []string{codes}
		case []any: // lookup returned type of slice is []any
			return common.Map(codes, func(it any) string {
				return it.(string)
			})
		}
		return []string{}

	default:
		panic(F.ToString("unknown GeoIP database type: ", string(db.MemoryType)))
	}
}

func (db Database) Close() error {
	return db.reader.Close()
}
