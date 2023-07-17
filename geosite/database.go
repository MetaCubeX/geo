package geosite

import (
	"errors"
	"os"

	"github.com/metacubex/geo/encoding/singgeo"
)

type Type string

const (
	TypeV2Ray = "V2Ray"
	TypeSing  = "sing-geosite"
)

var ErrInvalidDatabase = errors.New("invalid GeoSite database")

func FromBytes(data []byte) (db *Database, err error) {
	// attempt parsing sing-geosite
	singGeoSite, codes, err := singgeo.LoadSite(data)
	if err == nil {
		db, err = loadSing(singGeoSite, codes)
		return db, nil
	}

	// attempt parsing V2Ray dat
	var ok bool
	db, ok = loadV2RayFromBytes(data)
	if ok {
		return db, nil
	}

	return nil, ErrInvalidDatabase
}

func FromFile(file string) (db *Database, err error) {
	// attempt parsing sing-geosite
	singGeoSite, codes, err := singgeo.LoadSiteFromFile(file)
	if err == nil {
		db, err = loadSing(singGeoSite, codes)
		return db, nil
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

	return nil, ErrInvalidDatabase
}
