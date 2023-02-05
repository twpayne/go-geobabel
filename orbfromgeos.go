package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geos"
)

func NewOrbGeometryFromGEOSGeom(geosGeom *geos.Geom) orb.Geometry {
	switch geosGeom.TypeID() {
	case geos.TypeIDPoint:
		return NewOrbPointFromGEOSGeom(geosGeom)
	case geos.TypeIDLineString:
		return NewOrbLineStringFromGEOSGeom(geosGeom)
	case geos.TypeIDLinearRing:
		return NewOrbRingFromGEOSGeom(geosGeom)
	case geos.TypeIDPolygon:
		return NewOrbPolygonFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiPoint:
		return NewOrbMultiPointFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiLineString:
		return NewOrbMultiLineStringFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiPolygon:
		return NewOrbMultiPolygonFromGEOSGeom(geosGeom)
	case geos.TypeIDGeometryCollection:
		return NewOrbCollectionFromGEOSGeom(geosGeom)
	default:
		panic(fmt.Sprintf("%s: unsupported GEOS type", geosGeom.Type()))
	}
}

func NewOrbPointFromGEOSGeom(geosGeom *geos.Geom) orb.Point {
	geosCoords := geosGeom.CoordSeq().ToCoords()
	return orb.Point{geosCoords[0][0], geosCoords[0][1]}
}

func NewOrbLineStringFromGEOSGeom(geosGeom *geos.Geom) orb.LineString {
	geosCoords := geosGeom.CoordSeq().ToCoords()
	orbLineString := make(orb.LineString, 0, len(geosCoords))
	for _, coord := range geosCoords {
		orbPoint := orb.Point{coord[0], coord[1]}
		orbLineString = append(orbLineString, orbPoint)
	}
	return orbLineString
}

func NewOrbRingFromGEOSGeom(geosGeom *geos.Geom) orb.Ring {
	geosCoords := geosGeom.CoordSeq().ToCoords()
	orbRing := make(orb.Ring, 0, len(geosCoords))
	for _, coord := range geosCoords {
		orbPoint := orb.Point{coord[0], coord[1]}
		orbRing = append(orbRing, orbPoint)
	}
	return orbRing
}

func NewOrbPolygonFromGEOSGeom(geosGeom *geos.Geom) orb.Polygon {
	geosNumInteriorRings := geosGeom.NumInteriorRings()
	orbPolygon := make(orb.Polygon, 0, 1+geosNumInteriorRings)
	orbRing := NewOrbGeometryFromGEOSGeom(geosGeom.ExteriorRing()).(orb.Ring) //nolint:forcetypeassert
	orbPolygon = append(orbPolygon, orbRing)
	for i := 0; i < geosNumInteriorRings; i++ {
		orbRing := NewOrbGeometryFromGEOSGeom(geosGeom.InteriorRing(i)).(orb.Ring) //nolint:forcetypeassert
		orbPolygon = append(orbPolygon, orbRing)
	}
	return orbPolygon
}

func NewOrbMultiPointFromGEOSGeom(geosGeom *geos.Geom) orb.MultiPoint {
	geosNumGeometries := geosGeom.NumGeometries()
	orbMultiPoint := make(orb.MultiPoint, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		orbPoint := NewOrbGeometryFromGEOSGeom(geosGeom.Geometry(i)).(orb.Point) //nolint:forcetypeassert
		orbMultiPoint = append(orbMultiPoint, orbPoint)
	}
	return orbMultiPoint
}

func NewOrbMultiLineStringFromGEOSGeom(geosGeom *geos.Geom) orb.MultiLineString {
	geosNumGeometries := geosGeom.NumGeometries()
	orbMultiLineString := make(orb.MultiLineString, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		orbLineString := NewOrbGeometryFromGEOSGeom(geosGeom.Geometry(i)).(orb.LineString) //nolint:forcetypeassert
		orbMultiLineString = append(orbMultiLineString, orbLineString)
	}
	return orbMultiLineString
}

func NewOrbMultiPolygonFromGEOSGeom(geosGeom *geos.Geom) orb.MultiPolygon {
	geosNumGeometries := geosGeom.NumGeometries()
	orbMultiPolygon := make(orb.MultiPolygon, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		orbPolygon := NewOrbGeometryFromGEOSGeom(geosGeom.Geometry(i)).(orb.Polygon) //nolint:forcetypeassert
		orbMultiPolygon = append(orbMultiPolygon, orbPolygon)
	}
	return orbMultiPolygon
}

func NewOrbCollectionFromGEOSGeom(geosGeom *geos.Geom) orb.Collection {
	geosNumGeometries := geosGeom.NumGeometries()
	orbCollection := make(orb.Collection, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		orbGeometry := NewOrbGeometryFromGEOSGeom(geosGeom.Geometry(i))
		orbCollection = append(orbCollection, orbGeometry)
	}
	return orbCollection
}
