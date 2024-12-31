# go-whosonfirst-external

Go package for working with external data sources in a Who's On First context.

## Documentation

Documentation is incomplete.

## Interfaces

This packages uses a handful of interface definitions to hide the details of any one external data source. Each external data source provides its own implementation of these interfaces. 

### Iterators

Iterators are used to walk (crawl, iterate) a collection of records from an external data source.

```
type Iterator interface {
	Iterate(context.Context, ...string) iter.Seq2[Record, error]
	Close() error
}
```

### Records

Records are the common interface for a location defined by an external data source.

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

_If you are reading this then it means this interface stands a good chance of changing still._

## Providers

Implementations of the above-mention interfaces are available for the following external data sources. Provider-specific implementations are identified by a URI which takes the form of:

```
{PROVIDER} + "//" + {DATA_FORMAT} + "/" + {OPTIONAL_PLACETYPE} + "?" + {OPTIONAL_QUERY_PARAMETERS}
```

Where `{PROVIDER}` is the unique label (scheme) for the external data source, `{DATA_FORMAT}` is the data format that the external data is encoded in and `{OPTIONAL_PLACETYPE}` is an additional per-provider placetype filter.

### Foursquare

Implement the `go-whosonfirst-external` interfaces for [Foursquare's Open Source POI dataset](https://opensource.foursquare.com/os-places/).

```
foursquare://parquet
```

By default the Foursquare iterator only emits the id, name and geometry (derived from the latitude and longitude) properties. For example:

```
$> go run cmd/iterate/main.go -iterator-uri 'foursquare://parquet' ~/data/foursquare/parquet/*.parquet

{"properties":{"id":"4878c69f3d6c44dbc55b277e","name":"Huge Brands"},"geometry":[0,0]}
{"properties":{"id":"aef27cd85bb544efb68fa03a","name":"Альфа"},"geometry":[0,0]}
{"properties":{"id":"16c1ee7ff3e046d162b771fc","name":"СТИПЛ"},"geometry":[0,0]}
{"properties":{"id":"14f4601b631548f6bb6a9e56","name":"Advance Auto Parts"},"geometry":[0,0]}
{"properties":{"id":"5a26cf2c9de23b0468340907","name":"Ziemasprieki"},"geometry":[0,0]}
{"properties":{"id":"4c913ad8b641236a03b97f79","name":"Richmond Tires - playaz only"},"geometry":[-122.35356862771455,37.95792449892255]}
{"properties":{"id":"9cad9bae8344473ed53e4b10","name":"Safe Harbor Project"},"geometry":[-122.79330140132753,38.56842541968038]}
{"properties":{"id":"e36f619d4357549cb86ba73c","name":"Tesco Chelmsford 2"},"geometry":[0,0]}
{"properties":{"id":"13e58c07ea76461a291e7bb6","name":"Www.blackbeardesign.com"},"geometry":[0,0]}
{"properties":{"id":"50faaecb460e4ce7e44d5b09","name":"Servicios Marinos del Golfo"},"geometry":[0,0]}
{"properties":{"id":"4b95c216f964a52034b234e3","name":"Cafe Martin"},"geometry":[-122.71260890434607,38.438739037633034]}
... and so on
```

To export all the properties associated with a Foursquare POI pass the `?properties=all` query parameter to the iterator URI. For example:

```
$>  go run cmd/iterate/main.go -iterator-uri 'foursquare://parquet?properties=all' ~/data/foursquare/parquet/*.parquet | less

{"properties":{"address":"4910 W Amelia Earhart Dr","admin_region":"","category_ids":["52f2ab2ebcbc57f1066b8b28"],"category_labels":["Retail \u003e Print Store"],"country":"US","date_closed":"","date_created":"2020-02-03","date_refreshed":"2023-12-02","facebook_id":0,"id":"4878c69f3d6c44dbc55b277e","instagram":"","locality":"Salt Lake City","name":"Huge Brands","po_box":"","post_town":"","postcode":"84116","region":"UT","tel":"(801) 355-0331","twitter":"","website":"https://www.hugebrands.com"},"geometry":[0,0]}
{"properties":{"address":"Фрунзе Ул., д. 5","admin_region":"","category_ids":["63be6904847c3692a84b9b98"],"category_labels":["Business and Professional Services \u003e Wholesaler"],"country":"RU","date_closed":"","date_created":"2013-10-30","date_refreshed":"2023-12-06","facebook_id":0,"id":"aef27cd85bb544efb68fa03a","instagram":"","locality":"Новосибирск","name":"Альфа","po_box":"","post_town":"","postcode":"630091","region":"Новосибирская область","tel":"","twitter":"","website":"http://www.alaba.ru"},"geometry":[0,0]}
{"properties":{"address":"Ленинский Пр., 18","admin_region":"","category_ids":[],"category_labels":[],"country":"RU","date_closed":"","date_created":"2023-12-06","date_refreshed":"2023-12-06","facebook_id":0,"id":"16c1ee7ff3e046d162b771fc","instagram":"","locality":"Москва","name":"СТИПЛ","po_box":"","post_town":"","postcode":"119071","region":"Москва","tel":"","twitter":"","website":""},"geometry":[0,0]}
{"properties":{"address":"485 US Highway 1","admin_region":"","category_ids":["63be6904847c3692a84b9be6"],"category_labels":["Retail \u003e Automotive Retail \u003e Car Parts and Accessories"],"country":"US","date_closed":"","date_created":"2021-12-02","date_refreshed":"2023-11-29","facebook_id":0,"id":"14f4601b631548f6bb6a9e56","instagram":"","locality":"Edison","name":"Advance Auto Parts","po_box":"","post_town":"","postcode":"08817","region":"NJ","tel":"(732) 985-3308","twitter":"","website":"https://stores.advanceautoparts.com/nj/edison/485-route-1-s"},"geometry":[0,0]}
{"properties":{"address":"Aroniju Iela 279","admin_region":"","category_ids":[],"category_labels":[],"country":"","date_closed":"","date_created":"2017-12-05","date_refreshed":"2020-11-22","facebook_id":0,"id":"5a26cf2c9de23b0468340907","instagram":"","locality":"Aizkraukle","name":"Ziemasprieki","po_box":"","post_town":"","postcode":"5101","region":"Aizkraukles novads","tel":"","twitter":"","website":""},"geometry":[0,0]}
{"properties":{"address":"1608 Market Ave","admin_region":"","category_ids":["52f2ab2ebcbc57f1066b8b44"],"category_labels":["Business and Professional Services \u003e Automotive Service \u003e Automotive Repair Shop"],"country":"US","date_closed":"","date_created":"2010-09-15","date_refreshed":"2024-10-15","facebook_id":169088089774148,"id":"4c913ad8b641236a03b97f79","instagram":"","locality":"San Pablo","name":"Richmond Tires - playaz only","po_box":"","post_town":"","postcode":"94806","region":"CA","tel":"(510) 237-2712","twitter":"","website":"http://www.bureaumembers.com/richmondtire"},"geometry":[-122.35356862771455,37.95792449892255]}
... and so on
```

### Overture

Implement the `go-whosonfirst-external` interfaces for [Overture Data's Places dataset](https://docs.overturemaps.org/guides/places/).

```
overture://parquet/places
```

_As of this writing the Overture iterator only emits the id, name and geometry Overture Data properties._

## DuckDB

This package uses the [marcboeker/go-duckdb](https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring) package for working with Parquet files. Because the `go-duckdb` package bundles all the plaform-specific "libduckdb.a" files it is _NOT_ included in this package's `vendor` directory (because it just makes everything too big).

This introduces some obvious compile-time problems. Per the [go-duckdb documentation](https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring) the best way to deal with this is to install and use the [goware/modvendor](https://github.com/goware/modvendor) tool as follows:

```
$> go install github.com/goware/modvendor@latest
$> go mod vendor
$> modvendor -copy="**/*.a **/*.h" -v
```

The is also a handy `modvendor` Makefile target (in this package) to make that last step easier.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/iterate cmd/iterate/main.go
go build -mod vendor -ldflags="-s -w" -o bin/assign-ancestors cmd/assign-ancestors/main.go
```

### iterate

Iterate through one or more URIs for an external data source and emit each record as line-separated JSON.

```
$> ./bin/iterate -h
Iterate through one or more URIs for an external data source and emit each record as line-separated JSON.
Usage:
	 ./bin/iterate uri(N) uri(N)
  -as-geojson
    	Emit records as GeoJSON Features.
  -iterator-uri string
    	A registered whosonfirst/go-whosonfirst-external/iterator.
```

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

_TBW_

For example:

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
* https://opensource.foursquare.com/os-places/
* https://docs.overturemaps.org/guides/places/