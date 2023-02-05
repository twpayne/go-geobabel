package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geos"
)

func NewOrbGeometryFromGEOSGeometry(geosGeometry *geos.Geom) orb.Geometry {
	switch geosGeometry.TypeID() {
	case geos.TypeIDPoint:
		geosCoords := geosGeometry.CoordSeq().ToCoords()
		return orb.Point{geosCoords[0][0], geosCoords[0][1]}
	case geos.TypeIDLineString:
		geosCoords := geosGeometry.CoordSeq().ToCoords()
		orbLineString := make(orb.LineString, 0, len(geosCoords))
		for _, coord := range geosCoords {
			orbPoint := orb.Point{coord[0], coord[1]}
			orbLineString = append(orbLineString, orbPoint)
		}
		return orbLineString
	case geos.TypeIDLinearRing:
		geosCoords := geosGeometry.CoordSeq().ToCoords()
		orbRing := make(orb.Ring, 0, len(geosCoords))
		for _, coord := range geosCoords {
			orbPoint := orb.Point{coord[0], coord[1]}
			orbRing = append(orbRing, orbPoint)
		}
		return orbRing
	case geos.TypeIDPolygon:
		geosNumInteriorRings := geosGeometry.NumInteriorRings()
		orbPolygon := make(orb.Polygon, 0, 1+geosNumInteriorRings)
		orbRing := NewOrbGeometryFromGEOSGeometry(geosGeometry.ExteriorRing()).(orb.Ring)
		orbPolygon = append(orbPolygon, orbRing)
		for i := 0; i < geosNumInteriorRings; i++ {
			orbRing := NewOrbGeometryFromGEOSGeometry(geosGeometry.InteriorRing(i)).(orb.Ring)
			orbPolygon = append(orbPolygon, orbRing)
		}
		return orbPolygon
	case geos.TypeIDMultiPoint:
		geosNumGeometries := geosGeometry.NumGeometries()
		orbMultiPoint := make(orb.MultiPoint, 0, geosNumGeometries)
		for i := 0; i < geosNumGeometries; i++ {
			orbPoint := NewOrbGeometryFromGEOSGeometry(geosGeometry.Geometry(i)).(orb.Point)
			orbMultiPoint = append(orbMultiPoint, orbPoint)
		}
		return orbMultiPoint
	case geos.TypeIDMultiLineString:
		geosNumGeometries := geosGeometry.NumGeometries()
		orbMultiLineString := make(orb.MultiLineString, 0, geosNumGeometries)
		for i := 0; i < geosNumGeometries; i++ {
			orbLineString := NewOrbGeometryFromGEOSGeometry(geosGeometry.Geometry(i)).(orb.LineString)
			orbMultiLineString = append(orbMultiLineString, orbLineString)
		}
		return orbMultiLineString
	case geos.TypeIDMultiPolygon:
		geosNumGeometries := geosGeometry.NumGeometries()
		orbMultiPolygon := make(orb.MultiPolygon, 0, geosNumGeometries)
		for i := 0; i < geosNumGeometries; i++ {
			orbPolygon := NewOrbGeometryFromGEOSGeometry(geosGeometry.Geometry(i)).(orb.Polygon)
			orbMultiPolygon = append(orbMultiPolygon, orbPolygon)
		}
		return orbMultiPolygon
	case geos.TypeIDGeometryCollection:
		geosNumGeometries := geosGeometry.NumGeometries()
		orbCollection := make(orb.Collection, 0, geosNumGeometries)
		for i := 0; i < geosNumGeometries; i++ {
			orbGeometry := NewOrbGeometryFromGEOSGeometry(geosGeometry.Geometry(i))
			orbCollection = append(orbCollection, orbGeometry)
		}
		return orbCollection
	default:
		panic(fmt.Sprintf("%s: unsupported GEOS type", geosGeometry.Type()))
	}
}
