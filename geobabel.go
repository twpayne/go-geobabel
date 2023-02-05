package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geos"
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

func NewOrbGeometryFromGeomGeometry(geomGeometry geom.T) orb.Geometry {
	switch geomGeometry := geomGeometry.(type) {
	case *geom.Point:
		return orbPointFromGeomPoint(geomGeometry)
	case *geom.LineString:
		return orbLineStringFromGeomLineString(geomGeometry)
	case *geom.LinearRing:
		return orbRingFromGeomLinearRing(geomGeometry)
	case *geom.Polygon:
		return orbPolygonFromGeomPolygon(geomGeometry)
	case *geom.MultiPoint:
		return orbMultiPointFromGeomMultiPoint(geomGeometry)
	case *geom.MultiLineString:
		return orbMultiLineStringFromGeomMultiLineString(geomGeometry)
	case *geom.MultiPolygon:
		return orbMultiPolygonFromGeomMultiPolygon(geomGeometry)
	case *geom.GeometryCollection:
		return orbCollectionFromGeomGeometryCollection(geomGeometry)
	default:
		panic(fmt.Sprintf("%T: unsupported type", geomGeometry))
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

func orbPointFromGeomPoint(geomPoint *geom.Point) orb.Point {
	geomFlatCoords := geomPoint.FlatCoords()
	return orb.Point{geomFlatCoords[0], geomFlatCoords[1]}
}

func orbLineStringFromGeomLineString(geomLineString *geom.LineString) orb.LineString {
	geomFlatCoords := geomLineString.FlatCoords()
	geomStride := geomLineString.Stride()
	orbLineString := make(orb.LineString, 0, len(geomFlatCoords)/geomStride)
	for i := 0; i < len(geomFlatCoords); i += geomStride {
		orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
		orbLineString = append(orbLineString, orbPoint)
	}
	return orbLineString
}

func orbRingFromGeomLinearRing(geomLinearRing *geom.LinearRing) orb.Ring {
	geomFlatCoords := geomLinearRing.FlatCoords()
	geomStride := geomLinearRing.Stride()
	orbRing := make(orb.Ring, 0, len(geomFlatCoords)/geomStride)
	for i := 0; i < len(geomFlatCoords); i += geomStride {
		orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
		orbRing = append(orbRing, orbPoint)
	}
	return orbRing
}

func orbPolygonFromGeomPolygon(geomPolygon *geom.Polygon) orb.Polygon {
	geomFlatCoords := geomPolygon.FlatCoords()
	geomEnds := geomPolygon.Ends()
	geomStride := geomPolygon.Stride()
	geomStart := 0
	orbPolygon := make(orb.Polygon, 0, len(geomEnds))
	for _, geomEnd := range geomEnds {
		orbRing := make(orb.Ring, 0, (geomEnd-geomStart)/geomStride)
		for i := geomStart; i < geomEnd; i += geomStride {
			orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
			orbRing = append(orbRing, orbPoint)
		}
		orbPolygon = append(orbPolygon, orbRing)
		geomStart = geomEnd
	}
	return orbPolygon
}

func orbMultiPointFromGeomMultiPoint(geomMultiPoint *geom.MultiPoint) orb.MultiPoint {
	geomFlatCoords := geomMultiPoint.FlatCoords()
	geomStride := geomMultiPoint.Stride()
	orbMultiPoint := make(orb.MultiPoint, 0, len(geomFlatCoords)/geomStride)
	for i := 0; i < len(geomFlatCoords); i += geomStride {
		orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
		orbMultiPoint = append(orbMultiPoint, orbPoint)
	}
	return orbMultiPoint
}

func orbMultiLineStringFromGeomMultiLineString(geomMultiLineString *geom.MultiLineString) orb.MultiLineString {
	geomFlatCoords := geomMultiLineString.FlatCoords()
	geomEnds := geomMultiLineString.Ends()
	geomStride := geomMultiLineString.Stride()
	geomStart := 0
	orbMultiLineString := make(orb.MultiLineString, 0, len(geomEnds))
	for _, geomEnd := range geomEnds {
		orbLineString := make(orb.LineString, 0, (geomEnd-geomStart)/geomStride)
		for i := geomStart; i < geomEnd; i += geomStride {
			orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
			orbLineString = append(orbLineString, orbPoint)
		}
		orbMultiLineString = append(orbMultiLineString, orbLineString)
		geomStart = geomEnd
	}
	return orbMultiLineString
}

func orbMultiPolygonFromGeomMultiPolygon(geomMultiPolygon *geom.MultiPolygon) orb.MultiPolygon {
	geomFlatCoords := geomMultiPolygon.FlatCoords()
	geomEndss := geomMultiPolygon.Endss()
	geomStride := geomMultiPolygon.Stride()
	geomStart := 0
	orbMultiPolygon := make(orb.MultiPolygon, 0, len(geomEndss))
	for _, geomEnds := range geomEndss {
		orbPolygon := make(orb.Polygon, 0, len(geomEnds))
		for _, geomEnd := range geomEnds {
			orbRing := make(orb.Ring, 0, (geomEnd-geomStart)/geomStride)
			for i := geomStart; i < geomEnd; i += geomStride {
				orbPoint := orb.Point{geomFlatCoords[i], geomFlatCoords[i+1]}
				orbRing = append(orbRing, orbPoint)
			}
			orbPolygon = append(orbPolygon, orbRing)
			geomStart = geomEnd
		}
		orbMultiPolygon = append(orbMultiPolygon, orbPolygon)
	}
	return orbMultiPolygon
}

func orbCollectionFromGeomGeometryCollection(geomGeometryCollection *geom.GeometryCollection) orb.Collection {
	geomNumGeoms := geomGeometryCollection.NumGeoms()
	orbCollection := make(orb.Collection, 0, geomNumGeoms)
	for i := 0; i < geomNumGeoms; i++ {
		geomGeometry := geomGeometryCollection.Geom(i)
		orbGeometry := NewOrbGeometryFromGeomGeometry(geomGeometry)
		orbCollection = append(orbCollection, orbGeometry)
	}
	return orbCollection
}
