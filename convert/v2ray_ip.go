package convert

import (
	"io"
	"net"

	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/sagernet/sing/common"
)

func V2RayIPToSing(geoipList []*v2raygeo.GeoIP, output io.Writer) error {
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            "sing-geoip",
		IPVersion:               6,
		RecordSize:              24,
		Inserter:                inserter.ReplaceWith,
		DisableIPv4Aliasing:     true,
		IncludeReservedNetworks: true,
	})
	if err != nil {
		return err
	}

	for _, geoipEntry := range geoipList {
		for _, cidrEntry := range geoipEntry.Cidr {
			ipAddress := net.IP(cidrEntry.Ip)
			if ip4 := ipAddress.To4(); ip4 != nil {
				ipAddress = ip4
			}
			ipNet := &net.IPNet{
				IP:   ipAddress,
				Mask: net.CIDRMask(int(cidrEntry.Prefix), len(ipAddress)*8),
			}
			err = writer.Insert(ipNet, mmdbtype.String(geoipEntry.CountryCode))
			if err != nil {
				return err
			}
		}
	}

	return common.Error(writer.WriteTo(output))
}

func V2RayIPToMetaV0(geoipList []*v2raygeo.GeoIP, output io.Writer) error {
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            "Meta-geoip0",
		IPVersion:               6,
		RecordSize:              24,
		Inserter:                inserter.ReplaceWith,
		DisableIPv4Aliasing:     true,
		IncludeReservedNetworks: true,
	})
	if err != nil {
		return err
	}

	for _, geoipEntry := range geoipList {
		for _, cidrEntry := range geoipEntry.Cidr {
			ipAddress := net.IP(cidrEntry.Ip)
			if ip4 := ipAddress.To4(); ip4 != nil {
				ipAddress = ip4
			}
			ipNet := net.IPNet{
				IP:   ipAddress,
				Mask: net.CIDRMask(int(cidrEntry.Prefix), len(ipAddress)*8),
			}
			_, record := writer.Get(ipAddress)
			switch typedRecord := record.(type) {
			case nil:
				record = mmdbtype.String(geoipEntry.CountryCode)
			case mmdbtype.String:
				record = mmdbtype.Slice{record, mmdbtype.String(geoipEntry.CountryCode)}
			case mmdbtype.Slice:
				record = append(typedRecord, mmdbtype.String(geoipEntry.CountryCode))
			default:
				panic("bad record type")
			}
			err = writer.Insert(&ipNet, record)
			if err != nil {
				return err
			}
		}
	}

	return common.Error(writer.WriteTo(output))
}
