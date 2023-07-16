package convert

import (
	"io"
	"net"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/oschwald/maxminddb-golang"
	"github.com/sagernet/sing/common"
)

type geoip2Enterprise struct {
	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"represented_country"`
}

func MaxMindIPToSing(binary []byte, output io.Writer) error {
	return maxmindToMMDB(binary, output, "sing-geoip")
}

func MaxMindIPToMetaV0(binary []byte, output io.Writer) error {
	return maxmindToMMDB(binary, output, "Meta-geoip0")
}

func maxmindToMMDB(binary []byte, output io.Writer, databaseType string) error {
	database, err := maxminddb.FromBytes(binary)
	if err != nil {
		return err
	}
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            databaseType,
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
	var country geoip2Enterprise
	var ipNet *net.IPNet
	for networks.Next() {
		ipNet, err = networks.Network(&country)
		if err != nil {
			return err
		}
		var code string
		if country.Country.IsoCode != "" {
			code = strings.ToLower(country.Country.IsoCode)
		} else if country.RegisteredCountry.IsoCode != "" {
			code = strings.ToLower(country.RegisteredCountry.IsoCode)
		} else if country.RepresentedCountry.IsoCode != "" {
			code = strings.ToLower(country.RepresentedCountry.IsoCode)
		} else if country.Continent.Code != "" {
			code = strings.ToLower(country.Continent.Code)
		} else {
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
