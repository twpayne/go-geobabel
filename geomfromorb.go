package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geom"
)

func NewGeomGeometryFromOrbGeometry(orbGeometry orb.Geometry) geom.T {
	switch orbGeometry := orbGeometry.(type) {
	case orb.Point:
		return NewGeomPointFromOrbPoint(orbGeometry)
	case orb.LineString:
		return NewGeomLineStringFromOrbLineString(orbGeometry)
	case orb.Ring:
		return NewGeomLinearRingFromOrbRing(orbGeometry)
	case orb.Polygon:
		return NewGeomPolygonFromOrbPolygon(orbGeometry)
	case orb.MultiPoint:
		return NewGeomMultiPointFromOrbMultiPoint(orbGeometry)
	case orb.MultiLineString:
		return NewGeomMultiLineStringFromOrbMultiLineString(orbGeometry)
	case orb.MultiPolygon:
		return NewGeomMultiPolygonFromOrbMultiPolygon(orbGeometry)
	case orb.Collection:
		return NewGeomGeometryCollectionFromOrbCollection(orbGeometry)
	default:
		panic(fmt.Sprintf("%T: unsupported orb type", orbGeometry))
	}
}

func NewGeomPointFromOrbPoint(orbPoint orb.Point) *geom.Point {
	geomFlatCoords := GeomFlatCoordsFromOrbPoint(orbPoint)
	return geom.NewPointFlat(geom.XY, geomFlatCoords)
}

func NewGeomLineStringFromOrbLineString(orbLineString orb.LineString) *geom.LineString {
	geomFlatCoords := GeomFlatCoordsFromOrbLineString(orbLineString)
	return geom.NewLineStringFlat(geom.XY, geomFlatCoords)
}

func NewGeomLinearRingFromOrbRing(orbRing orb.Ring) *geom.LinearRing {
	geomFlatCoords := GeomFlatCoordsFromOrbRing(orbRing)
	return geom.NewLinearRingFlat(geom.XY, geomFlatCoords)
}

func NewGeomPolygonFromOrbPolygon(orbPolygon orb.Polygon) *geom.Polygon {
	geomFlatCoords, geomEnds := GeomFlatCoordsFromOrbPolygon(orbPolygon)
	return geom.NewPolygonFlat(geom.XY, geomFlatCoords, geomEnds)
}

func NewGeomMultiPointFromOrbMultiPoint(orbMultiPoint orb.MultiPoint) *geom.MultiPoint {
	geomFlatCoords := GeomFlatCoordsFromOrbMultiPoint(orbMultiPoint)
	return geom.NewMultiPointFlat(geom.XY, geomFlatCoords)
}

func NewGeomMultiLineStringFromOrbMultiLineString(orbMultiLineString orb.MultiLineString) *geom.MultiLineString {
	geomFlatCoords, geomEnds := GeomFlatCoordsFromOrbMultiLineString(orbMultiLineString)
	return geom.NewMultiLineStringFlat(geom.XY, geomFlatCoords, geomEnds)
}

func NewGeomMultiPolygonFromOrbMultiPolygon(orbMultiPolygon orb.MultiPolygon) *geom.MultiPolygon {
	geomFlatCoords, geomEndss := GeomFlatCoordsFromOrbMultiPolygon(orbMultiPolygon)
	return geom.NewMultiPolygonFlat(geom.XY, geomFlatCoords, geomEndss)
}

func NewGeomGeometryCollectionFromOrbCollection(orbCollection orb.Collection) *geom.GeometryCollection {
	geomGeometries := make([]geom.T, 0, len(orbCollection))
	for _, orbGeometery := range orbCollection {
		geomGeometry := NewGeomGeometryFromOrbGeometry(orbGeometery)
		geomGeometries = append(geomGeometries, geomGeometry)
	}
	return geom.NewGeometryCollection().MustPush(geomGeometries...)
}

func GeomFlatCoordsFromOrbPoint(orbPoint orb.Point) []float64 {
	return []float64{orbPoint.X(), orbPoint.Y()}
}

func GeomFlatCoordsFromOrbLineString(orbLineString orb.LineString) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbLineString))
	for _, orbPoint := range orbLineString {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func GeomFlatCoordsFromOrbRing(orbRing orb.Ring) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbRing))
	for _, orbPoint := range orbRing {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func GeomFlatCoordsFromOrbPolygon(orbPolygon orb.Polygon) ([]float64, []int) {
	geomFlatCoordsLen := 0
	geomEnds := make([]int, 0, len(orbPolygon))
	for _, orbRing := range orbPolygon {
		geomFlatCoordsLen += 2 * len(orbRing)
		geomEnds = append(geomEnds, geomFlatCoordsLen)
	}
	geomFlatCoords := make([]float64, 0, geomFlatCoordsLen)
	for _, orbRing := range orbPolygon {
		for _, orbPoint := range orbRing {
			geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
		}
	}
	return geomFlatCoords, geomEnds
}

func GeomFlatCoordsFromOrbMultiPoint(orbMultiPoint orb.MultiPoint) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbMultiPoint))
	for _, orbPoint := range orbMultiPoint {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func GeomFlatCoordsFromOrbMultiLineString(orbMultiLineString orb.MultiLineString) ([]float64, []int) {
	geomFlatCoordsLen := 0
	geomEnds := make([]int, 0, len(orbMultiLineString))
	for _, orbLineString := range orbMultiLineString {
		geomFlatCoordsLen += 2 * len(orbLineString)
		geomEnds = append(geomEnds, geomFlatCoordsLen)
	}
	geomFlatCoords := make([]float64, 0, geomFlatCoordsLen)
	for _, orbLineString := range orbMultiLineString {
		for _, orbPoint := range orbLineString {
			geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
		}
	}
	return geomFlatCoords, geomEnds
}

func GeomFlatCoordsFromOrbMultiPolygon(orbMultiPolygon orb.MultiPolygon) ([]float64, [][]int) {
	geomFlatCoordsLen := 0
	geomEndss := make([][]int, 0, len(orbMultiPolygon))
	for _, orbPolygon := range orbMultiPolygon {
		geomEnds := make([]int, 0, len(orbPolygon))
		for _, orbRing := range orbPolygon {
			geomFlatCoordsLen += 2 * len(orbRing)
			geomEnds = append(geomEnds, geomFlatCoordsLen)
		}
		geomEndss = append(geomEndss, geomEnds)
	}
	geomFlatCoords := make([]float64, 0, geomFlatCoordsLen)
	for _, orbPolygon := range orbMultiPolygon {
		for _, orbRing := range orbPolygon {
			for _, orbPoint := range orbRing {
				geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
			}
		}
	}
	return geomFlatCoords, geomEndss
}
