# go-whosonfirst-external

Go package for working with external data sources in a Who's On First context.

## Documentation

Documentation is incomplete.

## Records

```
type Record interface {
	Id() string
	Name() string
	Placetype() string
	Namespace() string
	Geometry() orb.Geometry
	Properties() map[string]any
}
```

## Tools

### iterate

For example:

```
$> go run cmd/iterate/main.go -iterator-uri foursquare://parquet ~/data/foursquare/parquet/*.parquet > /dev/null
2024/12/24 18:07:50 INFO Time to iterate records count=104529230 time=1m40.40977525s
```

```
$> go run cmd/iterate/main.go -iterator-uri overture://parquet/places ~/data/overture/parquet/*.parquet > /dev/null
2024/12/24 18:04:48 INFO Time to iterate records count=55527168 time=2m11.528865417s
```

### assign-ancestors

```
#!/bin/sh

TILES=whosonfirst-point-in-polygon-z13-20241213
DATA=${HOME}/data

ITERATOR=foursquare://parquet
WORKERS=3

CACHE_SIZE=10000

SPATIAL_DB="pmtiles://?tiles=file://${DATA}/whosonfirst/&database=${TILES}&enable-cache=false&pmtiles-cache-size=${CACHE_SIZE}&zoom=13&layer=whosonfirst"
		       
go run cmd/assign-ancestors/main.go \
   -workers ${WORKERS} \
   -iterator-uri ${ITERATOR} \
   -spatial-database-uri "${SPATIAL_DB}" \
   -properties-reader-uri "sql://sqlite3/geojson/id/body?dsn=${DATA}/whosonfirst/whosonfirst-data-admin-latest.db&parse-uri=true" \
   $@
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles