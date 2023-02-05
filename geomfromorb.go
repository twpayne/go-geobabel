package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geom"
)

func NewGeomGeometryFromOrbGeometry(orbGeometry orb.Geometry) geom.T {
	switch orbGeometry := orbGeometry.(type) {
	case orb.Point:
		geomFlatCoords := geomFlatCoordsFromOrbPoint(orbGeometry)
		return geom.NewPointFlat(geom.XY, geomFlatCoords)
	case orb.LineString:
		geomFlatCoords := geomFlatCoordsFromOrbLineString(orbGeometry)
		return geom.NewLineStringFlat(geom.XY, geomFlatCoords)
	case orb.Ring:
		geomFlatCoords := geomFlatCoordsFromOrbRing(orbGeometry)
		return geom.NewLinearRingFlat(geom.XY, geomFlatCoords)
	case orb.Polygon:
		geomFlatCoords, geomEnds := geomFlatCoordsFromOrbPolygon(orbGeometry)
		return geom.NewPolygonFlat(geom.XY, geomFlatCoords, geomEnds)
	case orb.MultiPoint:
		geomFlatCoords := geomFlatCoordsFromOrbMultiPoint(orbGeometry)
		return geom.NewMultiPointFlat(geom.XY, geomFlatCoords)
	case orb.MultiLineString:
		geomFlatCoords, geomEnds := geomFlatCoordsFromOrbMultiLineString(orbGeometry)
		return geom.NewMultiLineStringFlat(geom.XY, geomFlatCoords, geomEnds)
	case orb.MultiPolygon:
		geomFlatCoords, geomEndss := geomFlatCoordsFromOrbMultiPolygon(orbGeometry)
		return geom.NewMultiPolygonFlat(geom.XY, geomFlatCoords, geomEndss)
	case orb.Collection:
		geomGeometries := make([]geom.T, 0, len(orbGeometry))
		for _, orbGeometery := range orbGeometry {
			geomGeometry := NewGeomGeometryFromOrbGeometry(orbGeometery)
			geomGeometries = append(geomGeometries, geomGeometry)
		}
		return geom.NewGeometryCollection().MustPush(geomGeometries...)
	default:
		panic(fmt.Sprintf("%T: unsupported orb type", orbGeometry))
	}
}

func geomFlatCoordsFromOrbPoint(orbPoint orb.Point) []float64 {
	return []float64{orbPoint.X(), orbPoint.Y()}
}

func geomFlatCoordsFromOrbLineString(orbLineString orb.LineString) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbLineString))
	for _, orbPoint := range orbLineString {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func geomFlatCoordsFromOrbRing(orbRing orb.Ring) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbRing))
	for _, orbPoint := range orbRing {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func geomFlatCoordsFromOrbPolygon(orbPolygon orb.Polygon) ([]float64, []int) {
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

func geomFlatCoordsFromOrbMultiPoint(orbMultiPoint orb.MultiPoint) []float64 {
	geomFlatCoords := make([]float64, 0, 2*len(orbMultiPoint))
	for _, orbPoint := range orbMultiPoint {
		geomFlatCoords = append(geomFlatCoords, orbPoint[0], orbPoint[1])
	}
	return geomFlatCoords
}

func geomFlatCoordsFromOrbMultiLineString(orbMultiLineString orb.MultiLineString) ([]float64, []int) {
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

func geomFlatCoordsFromOrbMultiPolygon(orbMultiPolygon orb.MultiPolygon) ([]float64, [][]int) {
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
