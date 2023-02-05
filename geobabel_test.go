package geobabel_test

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geos"

	"github.com/twpayne/go-geobabel"
)

func TestAll(t *testing.T) {
	geosContext := geos.NewContext()
	for _, tc := range []struct {
		name        string
		geomT       geom.T
		geosGeom    *geos.Geom
		orbGeometry orb.Geometry
		skipWKB     string
	}{
		{
			name:        "Point",
			geomT:       geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			geosGeom:    geosContext.NewPoint([]float64{1, 2}),
			orbGeometry: orb.Point{1, 2},
		},
		{
			name:        "LineString",
			geomT:       geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			geosGeom:    geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
			orbGeometry: orb.LineString{{1, 2}, {3, 4}},
		},
		{
			name:        "LinearRing",
			geomT:       geom.NewLinearRing(geom.XY).MustSetCoords([]geom.Coord{{1, 4}, {5, 2}, {3, 6}, {1, 4}}),
			geosGeom:    geosContext.NewLinearRing([][]float64{{1, 4}, {5, 2}, {3, 6}, {1, 4}}),
			orbGeometry: orb.Ring{{1, 4}, {5, 2}, {3, 6}, {1, 4}},
			skipWKB:     "WKB does not support LinearRings",
		},
		{
			name: "Polygon",
			geomT: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			}),
			geosGeom: geosContext.NewPolygon([][][]float64{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			}),
			orbGeometry: orb.Polygon{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			},
		},
		{
			name:  "MultiPoint",
			geomT: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			geosGeom: geosContext.NewCollection(
				geos.TypeIDMultiPoint,
				[]*geos.Geom{
					geosContext.NewPoint([]float64{1, 2}),
					geosContext.NewPoint([]float64{3, 4}),
				}),
			orbGeometry: orb.MultiPoint{{1, 2}, {3, 4}},
		},
		{
			name: "MultiLineString",
			geomT: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{
				{{1, 2}, {3, 4}},
				{{5, 6}, {7, 8}},
			}),
			geosGeom: geosContext.NewCollection(
				geos.TypeIDMultiLineString,
				[]*geos.Geom{
					geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
					geosContext.NewLineString([][]float64{{5, 6}, {7, 8}}),
				},
			),
			orbGeometry: orb.MultiLineString{
				{{1, 2}, {3, 4}},
				{{5, 6}, {7, 8}},
			},
		},
		{
			name: "MultiPolygon",
			geomT: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{
				{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
				{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}},
			}),
			geosGeom: geosContext.NewCollection(
				geos.TypeIDMultiPolygon,
				[]*geos.Geom{
					geosContext.NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
					geosContext.NewPolygon([][][]float64{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}}),
				},
			),
			orbGeometry: orb.MultiPolygon{
				{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
				{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}},
			},
		},
		{
			name: "GeometryCollection",
			geomT: geom.NewGeometryCollection().MustPush(
				geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
				geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
				geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
			),
			geosGeom: geosContext.NewCollection(
				geos.TypeIDGeometryCollection,
				[]*geos.Geom{
					geosContext.NewPoint([]float64{1, 2}),
					geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
					geosContext.NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
				},
			),
			orbGeometry: orb.Collection{
				orb.Point{1, 2},
				orb.LineString{{1, 2}, {3, 4}},
				orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.True(t, tc.geosGeom.IsValid())
			require.Equal(t, "Valid Geometry", tc.geosGeom.IsValidReason())

			assert.Equal(t, tc.geomT, geobabel.NewGeomTFromOrbGeometry(tc.orbGeometry))

			assert.True(t, tc.geosGeom.Equals(geobabel.NewGEOSGeomFromOrbGeometry(geosContext, tc.orbGeometry)))

			assert.Equal(t, tc.orbGeometry, geobabel.NewOrbGeometryFromGEOSGeom(tc.geosGeom))
			assert.Equal(t, tc.orbGeometry, geobabel.NewOrbGeometryFromGeomT(tc.geomT))

			if tc.skipWKB == "" {
				geomWKB, err := geobabel.WKBFromGeomT(tc.geomT)
				require.NoError(t, err)
				geosWKB := geobabel.WKBFromGEOSGeom(tc.geosGeom)
				orbWKB := geobabel.WKBFromOrbGeometry(tc.orbGeometry)

				geomT, err := geobabel.NewGeomTFromWKB(geomWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.geomT, geomT)

				geomT, err = geobabel.NewGeomTFromWKB(geosWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.geomT, geomT)

				geomT, err = geobabel.NewGeomTFromWKB(orbWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.geomT, geomT)

				geosGeom, err := geobabel.NewGEOSGeomFromWKB(geosContext, geomWKB)
				assert.NoError(t, err)
				assert.True(t, tc.geosGeom.Equals(geosGeom))

				geosGeom, err = geobabel.NewGEOSGeomFromWKB(geosContext, geosWKB)
				assert.NoError(t, err)
				assert.True(t, tc.geosGeom.Equals(geosGeom))

				geosGeom, err = geobabel.NewGEOSGeomFromWKB(geosContext, orbWKB)
				assert.NoError(t, err)
				assert.True(t, tc.geosGeom.Equals(geosGeom))

				orbGeometry, err := geobabel.NewOrbGeometryFromWKB(geomWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.orbGeometry, orbGeometry)

				orbGeometry, err = geobabel.NewOrbGeometryFromWKB(geosWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.orbGeometry, orbGeometry)

				orbGeometry, err = geobabel.NewOrbGeometryFromWKB(orbWKB)
				assert.NoError(t, err)
				assert.Equal(t, tc.orbGeometry, orbGeometry)
			}
		})
	}
}
