package geobabel

import (
	"encoding/binary"

	"github.com/paulmach/orb"
	orbwkb "github.com/paulmach/orb/encoding/wkb"
	"github.com/twpayne/go-geom"
	geomwkb "github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geos"
)

var wkbByteOrder = binary.LittleEndian

func NewGEOSGeomFromWKB(geosContext *geos.Context, wkb []byte) (*geos.Geom, error) {
	return geosContext.NewGeomFromWKB(wkb)
}

func NewGeomTFromWKB(wkb []byte) (geom.T, error) {
	return geomwkb.Unmarshal(wkb)
}

func NewOrbGeometryFromWKB(wkb []byte) (orb.Geometry, error) {
	return orbwkb.Unmarshal(wkb)
}

func WKBFromGEOSGeom(geosGeom *geos.Geom) []byte {
	return geosGeom.ToWKB()
}

func WKBFromGeomT(geomT geom.T) ([]byte, error) {
	return geomwkb.Marshal(geomT, wkbByteOrder)
}

func WKBFromOrbGeometry(orbGeometry orb.Geometry) []byte {
	return orbwkb.MustMarshal(orbGeometry, wkbByteOrder)
}
