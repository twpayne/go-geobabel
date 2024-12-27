package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geom"
)

func NewOrbGeometryFromGeomT(geomT geom.T) orb.Geometry {
	switch geomT := geomT.(type) {
	case *geom.Point:
		return orbPointFromGeomPoint(geomT)
	case *geom.LineString:
		return orbLineStringFromGeomLineString(geomT)
	case *geom.LinearRing:
		return orbRingFromGeomLinearRing(geomT)
	case *geom.Polygon:
		return orbPolygonFromGeomPolygon(geomT)
	case *geom.MultiPoint:
		return orbMultiPointFromGeomMultiPoint(geomT)
	case *geom.MultiLineString:
		return orbMultiLineStringFromGeomMultiLineString(geomT)
	case *geom.MultiPolygon:
		return orbMultiPolygonFromGeomMultiPolygon(geomT)
	case *geom.GeometryCollection:
		return orbCollectionFromGeomGeometryCollection(geomT)
	default:
		panic(fmt.Sprintf("%T: unsupported type", geomT))
	}
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
	orbPolygon := make(orb.Polygon, len(geomEnds))
	for i, geomEnd := range geomEnds {
		orbRing := make(orb.Ring, 0, (geomEnd-geomStart)/geomStride)
		for j := geomStart; j < geomEnd; j += geomStride {
			orbPoint := orb.Point{geomFlatCoords[j], geomFlatCoords[j+1]}
			orbRing = append(orbRing, orbPoint)
		}
		orbPolygon[i] = orbRing
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
	orbMultiLineString := make(orb.MultiLineString, len(geomEnds))
	for i, geomEnd := range geomEnds {
		orbLineString := make(orb.LineString, 0, (geomEnd-geomStart)/geomStride)
		for j := geomStart; j < geomEnd; j += geomStride {
			orbPoint := orb.Point{geomFlatCoords[j], geomFlatCoords[j+1]}
			orbLineString = append(orbLineString, orbPoint)
		}
		orbMultiLineString[i] = orbLineString
		geomStart = geomEnd
	}
	return orbMultiLineString
}

func orbMultiPolygonFromGeomMultiPolygon(geomMultiPolygon *geom.MultiPolygon) orb.MultiPolygon {
	geomFlatCoords := geomMultiPolygon.FlatCoords()
	geomEndss := geomMultiPolygon.Endss()
	geomStride := geomMultiPolygon.Stride()
	geomStart := 0
	orbMultiPolygon := make(orb.MultiPolygon, len(geomEndss))
	for i, geomEnds := range geomEndss {
		orbPolygon := make(orb.Polygon, len(geomEnds))
		for j, geomEnd := range geomEnds {
			orbRing := make(orb.Ring, 0, (geomEnd-geomStart)/geomStride)
			for k := geomStart; k < geomEnd; k += geomStride {
				orbPoint := orb.Point{geomFlatCoords[k], geomFlatCoords[k+1]}
				orbRing = append(orbRing, orbPoint)
			}
			orbPolygon[j] = orbRing
			geomStart = geomEnd
		}
		orbMultiPolygon[i] = orbPolygon
	}
	return orbMultiPolygon
}

func orbCollectionFromGeomGeometryCollection(geomGeometryCollection *geom.GeometryCollection) orb.Collection {
	geomNumGeoms := geomGeometryCollection.NumGeoms()
	orbCollection := make(orb.Collection, geomNumGeoms)
	for i := range geomNumGeoms {
		geomT := geomGeometryCollection.Geom(i)
		orbGeometry := NewOrbGeometryFromGeomT(geomT)
		orbCollection[i] = orbGeometry
	}
	return orbCollection
}
