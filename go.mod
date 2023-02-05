module github.com/twpayne/go-geobabel

go 1.19

require (
	github.com/paulmach/orb v0.8.0
	github.com/stretchr/testify v1.8.1
	github.com/twpayne/go-geom v1.5.0
	github.com/twpayne/go-geos v0.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/twpayne/go-geos => ../go-geos
