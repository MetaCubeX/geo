package geoip

import (
	"errors"
	"os"

	"github.com/oschwald/maxminddb-golang"
)

type Type string

const (
	TypeMaxmind Type = "MaxMind"
	TypeV2Ray        = "V2Ray"
	TypeMetaV0       = "Meta-geoip0"
	TypeSing         = "sing-geoip"
)

var ErrInvalidGeoIPDatabase = errors.New("invalid GeoIP database")

func FromBytes(data []byte) (db *Database, err error) {
	var mmdbReader *maxminddb.Reader
	mmdbReader, err = maxminddb.FromBytes(data)
	if err == nil {
		db = &Database{
			reader:     mmdbReader,
			SourceType: Type(mmdbReader.Metadata.DatabaseType),
			MemoryType: Type(mmdbReader.Metadata.DatabaseType),
		}
		switch mmdbReader.Metadata.DatabaseType {
		case TypeSing, TypeMetaV0:
		default:
			db.SourceType = TypeMaxmind
			db.MemoryType = TypeMaxmind
		}
		return
	}
	switch err.(type) {
	case maxminddb.InvalidDatabaseError: // not MMDB, skip
	default:
		return
	}

	// attempt parsing V2Ray dat
	var ok bool
	db, ok = loadV2RayFromBytes(data)
	if ok {
		return db, nil
	}

	return nil, ErrInvalidGeoIPDatabase
}

func FromFile(file string) (db *Database, err error) {
	var mmdbReader *maxminddb.Reader
	mmdbReader, err = maxminddb.Open(file)
	if err == nil {
		db = &Database{
			reader:     mmdbReader,
			SourceType: Type(mmdbReader.Metadata.DatabaseType),
			MemoryType: Type(mmdbReader.Metadata.DatabaseType),
		}
		switch mmdbReader.Metadata.DatabaseType {
		case TypeSing, TypeMetaV0:
		default:
			db.SourceType = TypeMaxmind
			db.MemoryType = TypeMaxmind
		}
		return
	}
	switch err.(type) {
	case maxminddb.InvalidDatabaseError: // not MMDB, skip
	default:
		return
	}

	// attempt parsing V2Ray dat
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var ok bool
	db, ok = loadV2RayFromBytes(fileContent)
	if ok {
		return db, nil
	}

	return nil, ErrInvalidGeoIPDatabase
}
