package geoip

import (
	"bytes"

	"github.com/metacubex/geo/convert"
	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/oschwald/maxminddb-golang"
)

func loadV2RayFromBytes(data []byte) (*Database, bool) {
	geoipList, err := v2raygeo.LoadIP(data)
	if err != nil {
		return nil, false
	}

	var buffer bytes.Buffer
	err = convert.V2RayIPToMetaV0(geoipList, &buffer)
	mmdb, err := maxminddb.FromBytes(buffer.Bytes())
	if err != nil {
		return nil, false
	}

	return &Database{
		reader:     mmdb,
		SourceType: TypeV2Ray,
		MemoryType: TypeMetaV0,
	}, true
}
