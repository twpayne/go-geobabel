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
