package convert

import (
	"io"
	"net"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/oschwald/maxminddb-golang"
	"github.com/sagernet/sing/common"
)

func MetaV0ToSing(binary []byte, output io.Writer) error {
	database, err := maxminddb.FromBytes(binary)
	if err != nil {
		return err
	}
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

	networks := database.Networks(maxminddb.SkipAliasedNetworks)
	var codes any
	var ipNet *net.IPNet
	for networks.Next() {
		ipNet, err = networks.Network(&codes)
		if err != nil {
			return err
		}
		switch codes := codes.(type) {
		case nil:
			continue
		case string:
			if codes != "" {
				err = writer.Insert(ipNet, mmdbtype.String(codes))
			}
		case []any: // network returned type of slice is []any
			err = writer.Insert(ipNet, mmdbtype.String(
				common.MaxBy(codes, func(it any) int {
					return len(it.(string))
				}).(string)))
		}
		if err != nil {
			return err
		}
	}
	err = networks.Err()
	if err != nil {
		return err
	}

	return common.Error(writer.WriteTo(output))
}
