# geobabel

[![PkgGoDev](https://pkg.go.dev/badge/github.com/twpayne/go-geobabel)](https://pkg.go.dev/github.com/twpayne/go-geobabel)

Package `geobabel` converts geometry types between popular geometry libraries.

Supported geometry types are:

* [`orb.Geometry`](https://pkg.go.dev/github.com/paulmach/orb#Geometry) from
  [`github.com/paulmach/orb`](https://github.com/paulmach/orb)

* [`geom.T`](https://pkg.go.dev/github.com/twpayne/go-geom#T) from
  [`github.com/twpayne/go-geom`](https://github.com/twpayne/go-geom)

* [`*geos.Geom`](https://pkg.go.dev/github.com/twpayne/go-geos#Geom) from
  [`github.com/twpayne/go-geos`](https://github.com/twpayne/go-geos)

`geobabel` exists because no single geometry library is perfect. For example:

* `github.com/paulmach/orb` is a pure Go library with friendly API, excellent
  [support for OpenStreetMap data](https://github.com/paulmach/orb) and MapBox
  Vector Tiles, but only supports 2D geometries, and has a heavy focus on the
  commonly-used EPSG:4326 and EPSG:3857 projections.

* `github.com/twpayne/go-geom` is extremely high performance pure Go library
  that supports multi-dimensional geometries, multiple encodings, and is
  projection-agnostic, at the expense of a more complex API and a very limited
  set of geometric operations.

* `github.com/twpayne/go-geos` provides an idiomatic Go API to the huge set of
  battle-tested geometric operations in the industry-standard [GEOS
  library](https://libgeos.org/), at the expense of cgo overhead, C-style memory
  management, and a focus on 2D geometries.

With `geobabel` you can combine the best aspects of each library. For example,
you can use a pure Go library like `go-geom` or `orb` to represent your
geometries and use a geometric operation that is only available in `go-geos`.

## Supported conversions

|                     | To `geom.T` | To `*geos.Geom` | To `orb.Geometry` |
| ------------------- | ----------- | --------------- | ----------------- |
| From `geom.T`       | n/a         | no              | yes               |
| From `*geos.Geom`   | no          | n/a             | yes               |
| From `orb.Geometry` | yes         | yes             | n/a               |

## License

MIT