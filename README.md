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

### Overture

Implement the `go-whosonfirst-external` interfaces for [Overture Data's Places dataset](https://docs.overturemaps.org/guides/places/).

```
overture://parquet/places
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
* https://opensource.foursquare.com/os-places/
* https://docs.overturemaps.org/guides/places/