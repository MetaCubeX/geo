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

func SingIPToMetaV0(binary []byte, output io.Writer) error {
	database, err := maxminddb.FromBytes(binary)
	if err != nil {
		return err
	}
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

	networks := database.Networks(maxminddb.SkipAliasedNetworks)
	var code string
	var ipNet *net.IPNet
	for networks.Next() {
		ipNet, err = networks.Network(&code)
		if err != nil {
			return err
		}
		if code == "" {
			continue
		}

		err = writer.Insert(ipNet, mmdbtype.String(code))
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
