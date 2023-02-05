package geobabel

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geos"
)

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
