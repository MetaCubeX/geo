package convert

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"

	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/sagernet/sing/common"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func V2RayIPToSing(geoipList []*v2raygeo.GeoIP, output io.Writer) error {
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: "sing-geoip",
		Languages: common.Map(geoipList, func(it *v2raygeo.GeoIP) string {
			return strings.ToLower(it.CountryCode)
		}),
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
			err = writer.Insert(&ipNet, mmdbtype.String(strings.ToLower(geoipEntry.CountryCode)))
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

	included := make([]netip.Prefix, 0, 4*len(geoipList))
	codeMap := make(map[netip.Prefix][]string, 4*len(geoipList))
	for _, geoipEntry := range geoipList {
		code := strings.ToLower(geoipEntry.CountryCode)
		for _, cidrEntry := range geoipEntry.Cidr {
			addr, ok := netip.AddrFromSlice(cidrEntry.Ip)
			if !ok {
				return fmt.Errorf("bad IP address: %v", cidrEntry.Ip)
			}
			addr = addr.Unmap()
			prefix := netip.PrefixFrom(addr, int(cidrEntry.Prefix))
			included = append(included, prefix)
			codeMap[prefix] = append(codeMap[prefix], code)
		}
	}
	included = common.Uniq(included)
	slices.SortFunc(included, func(a, b netip.Prefix) int {
		// sort in ascending order
		return cmpCompare(a.Bits(), b.Bits())
	})

	for _, prefix := range included {
		ipAddress := net.IP(prefix.Addr().AsSlice())
		ipNet := net.IPNet{
			IP:   ipAddress,
			Mask: net.CIDRMask(prefix.Bits(), len(ipAddress)*8),
		}
		codes := codeMap[prefix]
		_, record := writer.Get(ipAddress)
		switch typedRecord := record.(type) {
		case nil:
			if len(codes) == 1 {
				record = mmdbtype.String(codes[0])
			} else {
				record = mmdbtype.Slice(common.Map(codes, func(it string) mmdbtype.DataType {
					return mmdbtype.String(it)
				}))
			}
		case mmdbtype.String:
			recordSlice := make(mmdbtype.Slice, 0, 1+len(codes))
			recordSlice = append(recordSlice, typedRecord)
			for _, code := range codes {
				recordSlice = append(recordSlice, mmdbtype.String(code))
			}
			recordSlice = common.Uniq(recordSlice)
			if len(recordSlice) == 1 {
				record = recordSlice[0]
			} else {
				record = recordSlice
			}
		case mmdbtype.Slice:
			recordSlice := typedRecord
			for _, code := range codes {
				recordSlice = append(recordSlice, mmdbtype.String(code))
			}
			recordSlice = common.Uniq(recordSlice)
			record = recordSlice
		default:
			panic("bad record type")
		}
		err = writer.Insert(&ipNet, record)
		if err != nil {
			return err
		}
	}

	return common.Error(writer.WriteTo(output))
}

// cmpCompare is a copy of cmp.Compare from the Go 1.21 release.
// T cannot be float types.
func cmpCompare[T constraints.Ordered](x, y T) int {
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}
