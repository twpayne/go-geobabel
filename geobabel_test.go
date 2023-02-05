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
		name string
		geom geom.T
		geos *geos.Geom
		orb  orb.Geometry
	}{
		{
			name: "Point",
			geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			geos: geosContext.NewPoint([]float64{1, 2}),
			orb:  orb.Point{1, 2},
		},
		{
			name: "LineString",
			geom: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			geos: geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
			orb:  orb.LineString{{1, 2}, {3, 4}},
		},
		{
			name: "LinearRing",
			geom: geom.NewLinearRing(geom.XY).MustSetCoords([]geom.Coord{{1, 4}, {5, 2}, {3, 6}, {1, 4}}),
			geos: geosContext.NewLinearRing([][]float64{{1, 4}, {5, 2}, {3, 6}, {1, 4}}),
			orb:  orb.Ring{{1, 4}, {5, 2}, {3, 6}, {1, 4}},
		},
		{
			name: "Polygon",
			geom: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			}),
			geos: geosContext.NewPolygon([][][]float64{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			}),
			orb: orb.Polygon{
				{{0, 0}, {4, 0}, {4, 4}, {0, 0}},
				{{2, 1}, {3, 1}, {3, 2}, {2, 1}},
			},
		},
		{
			name: "MultiPoint",
			geom: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			geos: geosContext.NewCollection(
				geos.TypeIDMultiPoint,
				[]*geos.Geom{
					geosContext.NewPoint([]float64{1, 2}),
					geosContext.NewPoint([]float64{3, 4}),
				}),
			orb: orb.MultiPoint{{1, 2}, {3, 4}},
		},
		{
			name: "MultiLineString",
			geom: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{
				{{1, 2}, {3, 4}},
				{{5, 6}, {7, 8}},
			}),
			geos: geosContext.NewCollection(
				geos.TypeIDMultiLineString,
				[]*geos.Geom{
					geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
					geosContext.NewLineString([][]float64{{5, 6}, {7, 8}}),
				},
			),
			orb: orb.MultiLineString{
				{{1, 2}, {3, 4}},
				{{5, 6}, {7, 8}},
			},
		},
		{
			name: "MultiPolygon",
			geom: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{
				{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
				{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}},
			}),
			geos: geosContext.NewCollection(
				geos.TypeIDMultiPolygon,
				[]*geos.Geom{
					geosContext.NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
					geosContext.NewPolygon([][][]float64{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}}),
				},
			),
			orb: orb.MultiPolygon{
				{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
				{{{2, 1}, {3, 1}, {3, 2}, {2, 1}}},
			},
		},
		{
			name: "GeometryCollection",
			geom: geom.NewGeometryCollection().MustPush(
				geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
				geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
				geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
			),
			geos: geosContext.NewCollection(
				geos.TypeIDGeometryCollection,
				[]*geos.Geom{
					geosContext.NewPoint([]float64{1, 2}),
					geosContext.NewLineString([][]float64{{1, 2}, {3, 4}}),
					geosContext.NewPolygon([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}),
				},
			),
			orb: orb.Collection{
				orb.Point{1, 2},
				orb.LineString{{1, 2}, {3, 4}},
				orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.True(t, tc.geos.IsValid())
			require.Equal(t, "Valid Geometry", tc.geos.IsValidReason())

			assert.Equal(t, tc.geom, geobabel.NewGeomTFromOrbGeometry(tc.orb))

			assert.True(t, tc.geos.Equals(geobabel.NewGEOSGeomFromOrbGeometry(geosContext, tc.orb)))

			assert.Equal(t, tc.orb, geobabel.NewOrbGeometryFromGEOSGeom(tc.geos))
			assert.Equal(t, tc.orb, geobabel.NewOrbGeometryFromGeomT(tc.geom))
		})
	}
}
