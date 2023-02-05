package geobabel

// FIXME add a geos.Geom.Flat method and use that instead
// FIXME test non-XY layouts

import (
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geos"
)

func NewGeomTFromGEOSGeom(geosGeom *geos.Geom) geom.T {
	switch geosGeom.TypeID() {
	case geos.TypeIDPoint:
		return NewGeomPointFromGEOSGeom(geosGeom)
	case geos.TypeIDLineString:
		return NewGeomLineStringFromGEOSGeom(geosGeom)
	case geos.TypeIDLinearRing:
		return NewGeomLinearRingFromGEOSGeom(geosGeom)
	case geos.TypeIDPolygon:
		return NewGeomPolygonFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiPoint:
		return NewGeomMultiPointFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiLineString:
		return NewGeomMultiLineStringFromGEOSGeom(geosGeom)
	case geos.TypeIDMultiPolygon:
		return NewGeomMultiPolygonFromGEOSGeom(geosGeom)
	case geos.TypeIDGeometryCollection:
		return NewGeomGeometryCollectionFromGEOSGeom(geosGeom)
	default:
		panic(fmt.Sprintf("%s: unsupported GEOS type", geosGeom.Type()))
	}
}

func NewGeomPointFromGEOSGeom(geosGeom *geos.Geom) *geom.Point {
	geosCoordSeq := geosGeom.CoordSeq()
	geomLayout := GeomLayoutFromGEOSDimensions(geosCoordSeq.Dimensions())
	return geom.NewPointFlat(geomLayout, geosCoordSeq.FlatCoords())
}

func NewGeomLineStringFromGEOSGeom(geosGeom *geos.Geom) *geom.LineString {
	geosCoordSeq := geosGeom.CoordSeq()
	geomLayout := GeomLayoutFromGEOSDimensions(geosCoordSeq.Dimensions())
	return geom.NewLineStringFlat(geomLayout, geosCoordSeq.FlatCoords())
}

func NewGeomLinearRingFromGEOSGeom(geosGeom *geos.Geom) *geom.LinearRing {
	geosCoordSeq := geosGeom.CoordSeq()
	geomLayout := GeomLayoutFromGEOSDimensions(geosCoordSeq.Dimensions())
	return geom.NewLinearRingFlat(geomLayout, geosCoordSeq.FlatCoords())
}

func NewGeomPolygonFromGEOSGeom(geosGeom *geos.Geom) *geom.Polygon {
	geosExteriorRingCoordSeq := geosGeom.ExteriorRing().CoordSeq()
	geomLayout := GeomLayoutFromGEOSDimensions(geosExteriorRingCoordSeq.Dimensions())
	geomFlatCoords := geosExteriorRingCoordSeq.FlatCoords()
	geomEnds := []int{len(geomFlatCoords)}
	for i, n := 0, geosGeom.NumInteriorRings(); i < n; i++ {
		geosInteriorRing := geosGeom.InteriorRing(i)
		geomFlatCoords = append(geomFlatCoords, geosInteriorRing.CoordSeq().FlatCoords()...)
		geomEnds = append(geomEnds, len(geomFlatCoords))
	}
	return geom.NewPolygonFlat(geomLayout, geomFlatCoords, geomEnds)
}

func NewGeomMultiPointFromGEOSGeom(geosGeom *geos.Geom) *geom.MultiPoint {
	geosNumGeometries := geosGeom.NumGeometries()
	geomLayout := geom.XY
	var geomFlatCoords []float64
	for i := 0; i < geosNumGeometries; i++ {
		geosCoordSeq := geosGeom.Geometry(i).CoordSeq()
		geomLayout = GeomLayoutFromGEOSDimensions(geosCoordSeq.Dimensions())
		geomFlatCoords = append(geomFlatCoords, geosCoordSeq.FlatCoords()...)
	}
	return geom.NewMultiPointFlat(geomLayout, geomFlatCoords)
}

func NewGeomMultiLineStringFromGEOSGeom(geosGeom *geos.Geom) *geom.MultiLineString {
	geosNumGeometries := geosGeom.NumGeometries()
	geomLayout := geom.XY
	var geomFlatCoords []float64
	geomEnds := make([]int, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		geosCoordSeq := geosGeom.Geometry(i).CoordSeq()
		geomLayout = GeomLayoutFromGEOSDimensions(geosCoordSeq.Dimensions())
		geomFlatCoords = append(geomFlatCoords, geosCoordSeq.FlatCoords()...)
		geomEnds = append(geomEnds, len(geomFlatCoords))
	}
	return geom.NewMultiLineStringFlat(geomLayout, geomFlatCoords, geomEnds)
}

func NewGeomMultiPolygonFromGEOSGeom(geosGeom *geos.Geom) *geom.MultiPolygon {
	geosNumGeometries := geosGeom.NumGeometries()
	geomLayout := geom.XY
	var geomFlatCoords []float64
	geomEndss := make([][]int, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		geosPolygon := geosGeom.Geometry(i)
		geosExteriorRingCoordSeq := geosPolygon.ExteriorRing().CoordSeq()
		geomLayout = GeomLayoutFromGEOSDimensions(geosExteriorRingCoordSeq.Dimensions())
		geomFlatCoords = append(geomFlatCoords, geosExteriorRingCoordSeq.FlatCoords()...)
		geomEnds := []int{len(geomFlatCoords)}
		for i, n := 0, geosGeom.NumInteriorRings(); i < n; i++ {
			geosInteriorRing := geosGeom.InteriorRing(i)
			geomFlatCoords = append(geomFlatCoords, geosInteriorRing.CoordSeq().FlatCoords()...)
			geomEnds = append(geomEnds, len(geomFlatCoords))
		}
		geomEndss = append(geomEndss, geomEnds)
	}
	return geom.NewMultiPolygonFlat(geomLayout, geomFlatCoords, geomEndss)
}

func NewGeomGeometryCollectionFromGEOSGeom(geosGeom *geos.Geom) *geom.GeometryCollection {
	geosNumGeometries := geosGeom.NumGeometries()
	geomGeomTs := make([]geom.T, 0, geosNumGeometries)
	for i := 0; i < geosNumGeometries; i++ {
		geosGeomT := NewGeomTFromGEOSGeom(geosGeom.Geometry(i))
		geomGeomTs = append(geomGeomTs, geosGeomT)
	}
	return geom.NewGeometryCollection().MustPush(geomGeomTs...)
}

func GeomLayoutFromGEOSDimensions(geosDimensions int) geom.Layout {
	switch {
	case geosDimensions == 2:
		return geom.XY
	case geosDimensions == 3:
		return geom.XYZ
	case geosDimensions == 4:
		return geom.XYZM
	case geosDimensions > 4:
		return geom.Layout(geosDimensions)
	default:
		panic(fmt.Sprintf("%d: unsupported GEOS dimensions", geosDimensions))
	}
}
