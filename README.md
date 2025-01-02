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

By default, the Overture iterator only emits the id, name and geometry Overture Data properties. For example:

```
$> go run cmd/iterate/main.go -iterator-uri overture://parquet/places ~/data/overture/parquet/*.parquet 
{"properties":{"id":"08ff39bac830c5900361ff7fe23acab8","name":"KK Beauty Shop 2"},"geometry":[-179.13203,-84.5792175]}
{"properties":{"id":"08ff39baeda4a2580336eaa84afac259","name":"Бряг Дуфек"},"geometry":[-179,-84.5]}
{"properties":{"id":"08ff39b25c2a605003370aee0592a959","name":"Capta Art Deals"},"geometry":[-178.3849454,-84.8698703]}
{"properties":{"id":"08ff2a6c8134db4303566a1fd3d15e4d","name":"Gerasimou-Gletscher"},"geometry":[-177.05,-84.7]}
{"properties":{"id":"08ff2a6c6c5008e003dd18787bcac2d1","name":"La Fuente del Negocio"},"geometry":[-174.5272207,-84.5147528]}
{"properties":{"id":"08ff2a6c6c762411031de95c443cb06e","name":"Cape Surprise"},"geometry":[-174.417,-84.5167]}
{"properties":{"id":"08ff2a654811b714035ca8aff0f891d2","name":"Mount Wade"},"geometry":[-174.31667,-84.85]}
{"properties":{"id":"08ff2a672522b46503151e1a75c8a756","name":"Krout Glacier"},"geometry":[-172.2,-84.8833]}
{"properties":{"id":"08ff2a61551312cd03aaa36ce434fae1","name":"Real Reels Motion Pictures"},"geometry":[-171.8789095,-84.5735595]}
{"properties":{"id":"08ff2a674586b8a40349a8ef35a45a92","name":"Mount Hall"},"geometry":[-170.367,-84.9167]}
... and so on
```

To export all the properties associated with a Overture "place" record pass the `?properties=all` query parameter to the iterator URI. For example:

```
$> go run cmd/iterate/main.go -iterator-uri 'overture://parquet/places?properties=all' ~/data/overture/parquet/*.parquet 
{"properties":{"addresses":[{"country":"CL","freeform":"Avenida Chacabuco 417","locality":"Concepción","postcode":"4030000","region":null}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":["spas"],"primary":"beauty_salon"},"confidence":0.4907482,"emails":[],"id":"08fb2d869b99d68b03680d213838a15f","name":"Centro estético integral NovaBelle","names":{"common":null,"primary":"Centro estético integral NovaBelle","rules":null},"phones":["+56413241617"],"socials":["https://www.facebook.com/1015636625139417"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"1015636625139417","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":["http://www.novabelle.cl/"]},"geometry":[-73.051,-36.83169]}
{"properties":{"addresses":[{"country":"CL","freeform":"Avenida Chacabuco 417","locality":"Concepción","postcode":"4030000","region":null}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":null,"primary":null},"confidence":0.26517966,"emails":[],"id":"08fb2d869b999915038f2ded11d425c8","name":"kine__integral","names":{"common":null,"primary":"kine__integral","rules":null},"phones":["+56982334047"],"socials":["https://www.facebook.com/103186637717488"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"103186637717488","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":[]},"geometry":[-73.0508543,-36.8316943]}
{"properties":{"addresses":[{"country":"CL","freeform":"Tijuana Parque Ecuador, Calle Lincoyán 14","locality":"Concepción","postcode":"4030000","region":null}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":["community_services_non_profits"],"primary":"fire_department"},"confidence":0.55927837,"emails":[],"id":"08fb2d869b88c2de0383e64d4313fa71","name":"Séptima Compañía de Bomberos Concepción","names":{"common":null,"primary":"Séptima Compañía de Bomberos Concepción","rules":null},"phones":["+56412253021"],"socials":["https://www.facebook.com/148149962466897"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"148149962466897","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":["http://www.septima.com/"]},"geometry":[-73.0495173,-36.834187]}
{"properties":{"addresses":[{"country":"CL","freeform":null,"locality":"Concepción","postcode":null,"region":null}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":["community_services_non_profits"],"primary":"public_and_government_association"},"confidence":0.26517966,"emails":[],"id":"08fb2d869b8808e00309ca31ac43fd43","name":"Abuelitas Chile","names":{"common":null,"primary":"Abuelitas Chile","rules":null},"phones":["+56946643895"],"socials":["https://www.facebook.com/2350453898575474"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"2350453898575474","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":["http://www.abuelitaschile.cl/"]},"geometry":[-73.05053,-36.83348]}
{"properties":{"addresses":[{"country":"CL","freeform":"Calle Víctor Lamas 361","locality":"Concepción","postcode":"4030000","region":null}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":["southern_restaurant","comfort_food_restaurant"],"primary":"sandwich_shop"},"confidence":0.60211265,"emails":[],"id":"08fb2d869b8808e00377f87f524f9f2d","name":"Fuente Penquista","names":{"common":null,"primary":"Fuente Penquista","rules":null},"phones":["+56954309433"],"socials":["https://www.facebook.com/983736055044613"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"983736055044613","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":[]},"geometry":[-73.05053,-36.83348]}
{"properties":{"addresses":[{"country":"CL","freeform":"Calle Víctor Lamas 371","locality":"Concepción","postcode":"4030000","region":"BI"}],"brand":{"names":null,"wikidata":null},"categories":{"alternate":null,"primary":"shopping"},"confidence":0.55927837,"emails":[],"id":"08fb2d869b880b3403b2f4fc95473f2e","name":"El Mesón","names":{"common":null,"primary":"El Mesón","rules":null},"phones":["+56992195014"],"socials":["https://www.facebook.com/904144773050750"],"sources":[{"confidence":null,"dataset":"meta","property":"","record_id":"904144773050750","update_time":"2024-09-10T00:00:00.000Z"}],"version":0,"websites":[]},"geometry":[-73.05043,-36.83347]}
...and so on
```

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

Iterate through one or more URIs for an external data source and reverse-geocode each record emitting the record ID, Who's On First parent ID and Who's On First ancestry as CSV data to STDOUT.

```
$> ./bin/assign-ancestors -h
Iterate through one or more URIs for an external data source and reverse-geocode each record emitting the record ID, Who's On First parent ID and Who's On First ancestry as CSV data to STDOUT.
Usage:
	 ./bin/assign-ancestors uri(N) uri(N)
  -iterator-uri string
    	A registered whosonfirst/go-whosonfirst-external/iterator.Iterator URI.
  -properties-reader-uri string
    	A registered whosonfirst/go-reader.Reader URI for reading properties from parent records. If '{spatial-database-uri}' the spatial database instance will be used to read those properties. (default "{spatial-database-uri}")
  -spatial-database-uri string
    	A registered whosonfirst/go-whosonfirst-spatial/database/SpatialDatabase URI to use for perforning reverse geocoding tasks.
  -start-after int
    	If > 0 then delay processing for 'start_after' number of records.
  -verbose
    	Enable verbose (debug) logging.
  -workers int
    	The maximum number of workers to process reverse geocoding tasks. (default 5)
```	

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

Where:

* `whosonfirst-point-in-polygon-z13-20241213` is a PMTiles database containing Who's On First data used for performing point-in-polygon operations. See the [" A global point-in-polygon service using a static 8GB data file"](https://millsfield.sfomuseum.org/blog/2022/12/19/pmtiles-pip/) blog post for details.
* `whosonfirst-data-admin-latest.db` is a Who's On First SQLite distribution, published by [Geocode.earth](https://geocode.earth/data/whosonfirst/combined/).

Notes:

* More "workers" is not necessarily better. I _think_ this has something to do with the underlying library used to query the PMTiles database throttling requests (but I am not sure).
* When building the PMTiles database (containing Who's On First data used for performing point-in-polygon operations) make sure to use a current (fall/winter 2024) version of `tippecanoe`.)
* Should this really write CSV files with per-placetype columns instead of a single JSON-encoded `wof:hierarchy` blob? I don't know either but that's what it does, today.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://opensource.foursquare.com/os-places/
* https://docs.overturemaps.org/guides/places/