package v2raygeo

import (
	"os"

	E "github.com/sagernet/sing/common/exceptions"
	"google.golang.org/protobuf/proto"
)

func LoadIP(geoipBytes []byte) ([]*GeoIP, error) {
	var geoipList GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}
	return geoipList.Entry, nil
}

func LoadSite(geositeBytes []byte) ([]*GeoSite, error) {
	var geositeList GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}
	return geositeList.Entry, nil
}

func LoadIPFromFile(filename string) ([]*GeoIP, error) {
	geoipBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, E.Cause(err, "failed to load V2Ray GeoIP database")
	}
	return LoadIP(geoipBytes)
}

func LoadSiteFromFile(filename string) ([]*GeoSite, error) {
	geositeBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, E.Cause(err, "failed to load V2Ray GeoSite database")
	}
	return LoadSite(geositeBytes)
}
