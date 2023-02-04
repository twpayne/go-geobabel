package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geos"
)

func NewGEOSGeometryFromOrbGeometry(geosContext *geos.Context, orbGeometry orb.Geometry) *geos.Geom {
	switch orbGeometry := orbGeometry.(type) {
	case orb.Point:
		return geosContext.NewPoint(geosCoordsFromOrbPoint(orbGeometry))
	case orb.LineString:
		return geosContext.NewLineString(geosCoordsFromOrbLineString(orbGeometry))
	case orb.Ring:
		return geosContext.NewLinearRing(geosCoordsFromOrbRing(orbGeometry))
	case orb.Polygon:
		return geosContext.NewPolygon(geosCoordsFromOrbPolygon(orbGeometry))
	case orb.MultiPoint:
		geosPoints := make([]*geos.Geom, 0, len(orbGeometry))
		for _, orbPoint := range orbGeometry {
			geosPoint := NewGEOSGeometryFromOrbGeometry(geosContext, orbPoint)
			geosPoints = append(geosPoints, geosPoint)
		}
		return geosContext.NewCollection(geos.TypeIDMultiPoint, geosPoints)
	case orb.MultiLineString:
		geosLineStrings := make([]*geos.Geom, 0, len(orbGeometry))
		for _, orbLineString := range orbGeometry {
			geosLineString := NewGEOSGeometryFromOrbGeometry(geosContext, orbLineString)
			geosLineStrings = append(geosLineStrings, geosLineString)
		}
		return geosContext.NewCollection(geos.TypeIDMultiLineString, geosLineStrings)
	case orb.MultiPolygon:
		geosPolygons := make([]*geos.Geom, 0, len(orbGeometry))
		for _, orbPolygon := range orbGeometry {
			geosPolygon := NewGEOSGeometryFromOrbGeometry(geosContext, orbPolygon)
			geosPolygons = append(geosPolygons, geosPolygon)
		}
		return geosContext.NewCollection(geos.TypeIDMultiPolygon, geosPolygons)
	case orb.Collection:
		geosGeometries := make([]*geos.Geom, 0, len(orbGeometry))
		for _, orbGeometry := range orbGeometry {
			geosGeometry := NewGEOSGeometryFromOrbGeometry(geosContext, orbGeometry)
			geosGeometries = append(geosGeometries, geosGeometry)
		}
		return geosContext.NewCollection(geos.TypeIDGeometryCollection, geosGeometries)
	default:
		panic(fmt.Sprintf("%T: unsupported type", orbGeometry))
	}
}

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

func geosCoordsFromOrbPoint(orbPoint orb.Point) []float64 {
	return []float64{orbPoint[0], orbPoint[1]}
}

func geosCoordsFromOrbLineString(orbLineString orb.LineString) [][]float64 {
	geosCoords := make([][]float64, 0, len(orbLineString))
	for _, orbPoint := range orbLineString {
		geosCoords = append(geosCoords, geosCoordsFromOrbPoint(orbPoint))
	}
	return geosCoords
}

func geosCoordsFromOrbRing(orbRing orb.Ring) [][]float64 {
	geosCoords := make([][]float64, 0, len(orbRing))
	for _, orbPoint := range orbRing {
		geosCoords = append(geosCoords, geosCoordsFromOrbPoint(orbPoint))
	}
	return geosCoords
}

func geosCoordsFromOrbPolygon(orbPolygon orb.Polygon) [][][]float64 {
	geosCoords := make([][][]float64, 0, len(orbPolygon))
	for _, orbRing := range orbPolygon {
		geosCoords = append(geosCoords, geosCoordsFromOrbRing(orbRing))
	}
	return geosCoords
}
