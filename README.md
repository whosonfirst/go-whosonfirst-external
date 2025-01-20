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
go build -mod vendor -ldflags="-s -w" -o bin/sort-ancestors cmd/sort-ancestors/main.go
go build -mod vendor -ldflags="-s -w" -o bin/walk-sorted cmd/walk-sorted/main.go
go build -mod vendor -ldflags="-s -w" -o bin/compile-area cmd/compile-area/main.go
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

For example, imagine a shell script called `assign-ancestors.sh`:

```
#!/bin/sh

DATA=${HOME}/data

ITERATOR=foursquare://parquet
WORKERS=3

PMTILES_DATABASE=whosonfirst-point-in-polygon-z13-20241213
PMTILES_CACHE_SIZE=10000

SPATIAL_DB="pmtiles://?tiles=file://${DATA}/whosonfirst/&database=${PMTILES_DATABASE}&enable-cache=false&pmtiles-cache-size=${PMTILES_CACHE_SIZE}&zoom=13&layer=whosonfirst"

PROPERTIES_READER="sql://sqlite3/geojson/id/body?dsn=${DATA}/whosonfirst/whosonfirst-data-admin-latest.db&parse-uri=true"
# PROPERTIES_READER="{spatial-database-uri}"

go run cmd/assign-ancestors/main.go \
   -workers ${WORKERS} \
   -iterator-uri ${ITERATOR} \
   -spatial-database-uri "${SPATIAL_DB}" \
   -properties-reader-uri "${PROPERTIES_READER}" \
   $@
```

Which would be invoked like this:

```
$> assign-ancestors.sh ~/data/foursquare/parquet/*.parquet
2025/01/03 09:11:58 INFO Status counter=9038 processed=9034 diff=0 "avg t2p"=2.701129067965464 elaspsed=10.00037225s
2025/01/03 09:12:08 INFO Status counter=20559 processed=20555 diff=11521 "avg t2p"=2.400924349306738 elaspsed=20.000024708s
2025/01/03 09:12:18 INFO Status counter=30606 processed=30602 diff=10047 "avg t2p"=2.445199660152931 elaspsed=30.000002792s
2025/01/03 09:12:19 INFO Time to prune databases total=6702 pruned=5101 time=264.497334ms
2025/01/03 09:12:28 INFO Status counter=42446 processed=42442 diff=11840 "avg t2p"=2.3315819235662785 elaspsed=40.000010458s
2025/01/03 09:12:38 INFO Status counter=55812 processed=55808 diff=13366 "avg t2p"=2.1864607224770642 elaspsed=50.000002958s
2025/01/03 09:12:48 INFO Status counter=69066 processed=69062 diff=13254 "avg t2p"=2.0958558975992587 elaspsed=1m0.000026125s
2025/01/03 09:12:49 INFO Time to prune databases total=3562 pruned=3364 time=132.416792ms
... and so on
```

And the output would look like this:

```
external:geometry,external:id,external:namespace,wof:country,wof:hierarchies,wof:parent_id
POINT(-122.35356862771455 37.95792449892255),4c913ad8b641236a03b97f79,4sq,US,"[{""continent_id"":102191575,""country_id"":85633793,""county_id"":102086225,""locality_id"":85922125,""region_id"":85688637}]",85922125
POINT(0 0),e36f619d4357549cb86ba73c,4sq,XY,,-1
POINT(-122.79330140132753 38.56842541968038),9cad9bae8344473ed53e4b10,4sq,US,"[{""continent_id"":102191575,""country_id"":85633793,""county_id"":102081671,""region_id"":85688637}]",102081671
POINT(0 0),304d3a4bb7464b9154a9efdb,4sq,XY,,-1
POINT(-122.71260890434607 38.438739037633034),4b95c216f964a52034b234e3,4sq,US,"[{""continent_id"":102191575,""country_id"":85633793,""county_id"":102081671,""locality_id"":85922693,""region_id"":85688637}]",85922693
... and so on
```

Where:

* `whosonfirst-point-in-polygon-z13-20241213` is a PMTiles database containing Who's On First data used for performing point-in-polygon operations. See the [" A global point-in-polygon service using a static 8GB data file"](https://millsfield.sfomuseum.org/blog/2022/12/19/pmtiles-pip/) blog post for details.
* `whosonfirst-data-admin-latest.db` is a Who's On First SQLite distribution, published by [Geocode.earth](https://geocode.earth/data/whosonfirst/combined/).
* The use of a SQLite database for a "properties reader" is not _necessary_. If absent then the PMTiles database will be used.

Notes:

* More "workers" is not necessarily better. I _think_ this has something to do with the underlying library used to query the PMTiles database throttling requests (but I am not sure).
* When building the PMTiles database (containing Who's On First data used for performing point-in-polygon operations) make sure to use a current (fall/winter 2024) version of `tippecanoe`.

### walk-sorted

Walk one or more [whosonfirst-data/whosonfirst-external-* repositories](https://github.com/whosonfirst-data/?q=whosonfirst-external-&type=all&language=&sort=), with optional filtering, emitting CSV-encoded rows to STDOUT.

```
$> ./bin/walk-sorted -h
Walk one or more whosonfirst-external-* repositories, with optional filtering, emitting CSV-encoded rows to STDOUT.
Usage:
	 ./bin/walk-sorted path(N) path(N) path(N)
  -ancestor-id value
    	Zero or more "wof:hierarchies" values to match.
  -geohash string
    	An optional geohash to do a prefix-first comparison against.
  -mode string
    	Indicate whether all or any filter criteria must match. Valid options are: any, all. (default "all")
  -parent-id value
    	One or more "wof:parent_id" values to match.
  -verbose
    	Enable verbose (debug) logging.
```

For example, to emit all the Foursquare venues in [San Francisco](https://spelunker.whosonfirst.org/id/102087579) and [Alameda](https://spelunker.whosonfirst.org/id/102086959) counties in the sub-directory for [California](https://spelunker.whosonfirst.org/id/85688637) in the [whosonfirst-external-foursquare-venue-us](https://github.com/whosonfirst-data/whosonfirst-external-foursquare-venue-us) repository:

```
$> ./bin/walk-sorted \
	-mode any \
	-ancestor-id 102087579 \
	-ancestor-id 102086959 \
	/usr/local/data/foursquare/whosonfirst/whosonfirst-external-foursquare-venue-us/data/85688637/ \
	| wc -l
	
  248631
```

### compile-area

Merge GeoParquet data for an external data source with one or more Who's On First ancestry sources in to a new GeoParquet file. Ancestry sources can either be individual whosonfirst-external-* bzip2-compressed CSV files or folders containing one or more bzip2-compressed CSV files.

```
$> ./bin/compile-area -h
Merge GeoParquet data for an external data source with one or more Who's On First ancestry sources in to a new GeoParquet file. Ancestry sources can either be individual whosonfirst-external-* bzip2-compressed CSV files or folders containing one or more bzip2-compressed CSV files.
Usage:
	 ./bin/compile-area uri(N) uri(N)
  -ancestor-id value
    	Zero or more "wof:hierarchies" values to match.
  -external-id-key string
    	The name of the unique identifier key for an external data source. The output is a new GeoParquet file.xs
  -external-source string
    	The string to pass to the DuckDB 'read_parquet' command for reading an external data source.
  -geohash string
    	An optional geohash to do a prefix-first comparison against.
  -mode string
    	Indicate whether all or any filter criteria must match. Valid options are: any, all. (default "all")
  -parent-id value
    	One or more "wof:parent_id" values to match.
  -target string
    	The path where the final GeoParquet file should be written.
  -verbose
    	Enable verbose (debug) logging.
```

For example, to create a new GeoParquet file (called `sfba.parquet`) for all the Foursquare venues in [San Francisco](https://spelunker.whosonfirst.org/id/102087579), [Alameda](https://spelunker.whosonfirst.org/id/102086959) and [San Mateo](https://spelunker.whosonfirst.org/id/102085387) counties, merged with their Who's On First ancestry data, in the sub-directory for [California](https://spelunker.whosonfirst.org/id/85688637) in the [whosonfirst-external-foursquare-venue-us](https://github.com/whosonfirst-data/whosonfirst-external-foursquare-venue-us) repository:

```
$> ./bin/compile-area \
	-external-source "/usr/local/data/foursquare/parquet/*.parquet" \
	-external-id-key fsq_place_id \
	-mode any \
	-ancestor-id 102087579 \
	-ancestor-id 102086959 \
	-ancestor-id 102085387 \
	-target sfba.parquet \
	/usr/local/data/foursquare/whosonfirst/whosonfirst-external-foursquare-venue-us/data/85688637
```

Investigate the `sfba.parquet` file:

```
$> du -h sfba.parquet 
 34M	sfba.parquet
 
$> duckdb
v1.1.3 19864453f7
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
D DESCRIBE (SELECT * FROM read_parquet('sfba.parquet'));
┌─────────────────────┬────────────────────────────────────────────────────────────┬─────────┬─────────┬─────────┬─────────┐
│     column_name     │                        column_type                         │  null   │   key   │ default │  extra  │
│       varchar       │                          varchar                           │ varchar │ varchar │ varchar │ varchar │
├─────────────────────┼────────────────────────────────────────────────────────────┼─────────┼─────────┼─────────┼─────────┤
│ fsq_place_id        │ VARCHAR                                                    │ YES     │         │         │         │
│ name                │ VARCHAR                                                    │ YES     │         │         │         │
│ latitude            │ DOUBLE                                                     │ YES     │         │         │         │
│ longitude           │ DOUBLE                                                     │ YES     │         │         │         │
│ address             │ VARCHAR                                                    │ YES     │         │         │         │
│ locality            │ VARCHAR                                                    │ YES     │         │         │         │
│ region              │ VARCHAR                                                    │ YES     │         │         │         │
│ postcode            │ VARCHAR                                                    │ YES     │         │         │         │
│ admin_region        │ VARCHAR                                                    │ YES     │         │         │         │
│ post_town           │ VARCHAR                                                    │ YES     │         │         │         │
│ po_box              │ VARCHAR                                                    │ YES     │         │         │         │
│ country             │ VARCHAR                                                    │ YES     │         │         │         │
│ date_created        │ VARCHAR                                                    │ YES     │         │         │         │
│ date_refreshed      │ VARCHAR                                                    │ YES     │         │         │         │
│ date_closed         │ VARCHAR                                                    │ YES     │         │         │         │
│ tel                 │ VARCHAR                                                    │ YES     │         │         │         │
│ website             │ VARCHAR                                                    │ YES     │         │         │         │
│ email               │ VARCHAR                                                    │ YES     │         │         │         │
│ facebook_id         │ BIGINT                                                     │ YES     │         │         │         │
│ instagram           │ VARCHAR                                                    │ YES     │         │         │         │
│ twitter             │ VARCHAR                                                    │ YES     │         │         │         │
│ fsq_category_ids    │ VARCHAR[]                                                  │ YES     │         │         │         │
│ fsq_category_labels │ VARCHAR[]                                                  │ YES     │         │         │         │
│ geom                │ BLOB                                                       │ YES     │         │         │         │
│ bbox                │ STRUCT(xmin DOUBLE, ymin DOUBLE, xmax DOUBLE, ymax DOUBLE) │ YES     │         │         │         │
│ geohash             │ VARCHAR                                                    │ YES     │         │         │         │
│ wof:country         │ VARCHAR                                                    │ YES     │         │         │         │
│ wof:parent_id       │ BIGINT                                                     │ YES     │         │         │         │
│ wof:hierarchies     │ VARCHAR                                                    │ YES     │         │         │         │
├─────────────────────┴────────────────────────────────────────────────────────────┴─────────┴─────────┴─────────┴─────────┤
│ 29 rows                                                                                                        6 columns │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

And then query for all the restaurants in the [Temescal](https://spelunker.whosonfirst.org/id/85872391) neighbourhood:

```
D SELECT fsq_place_id, name, address, JSON("wof:hierarchies")[0].neighbourhood_id AS neighbourhood, latitude, longitude, date_closed FROM read_parquet('sfba.parquet') WHERE neighbourhood=85872391 AND JSON(fsq_category_labels)[0]  LIKE '%Dining and Drinking > Restaurant%';
┌──────────────────────────┬───────────────────────────────────┬──────────────────────────┬───────────────┬────────────────────┬─────────────────────┬─────────────┐
│       fsq_place_id       │               name                │         address          │ neighbourhood │      latitude      │      longitude      │ date_closed │
│         varchar          │              varchar              │         varchar          │     json      │       double       │       double        │   varchar   │
├──────────────────────────┼───────────────────────────────────┼──────────────────────────┼───────────────┼────────────────────┼─────────────────────┼─────────────┤
│ 55663e24498ed3e6077e2282 │ Rosamunde Sausage Grill           │ 4659 Telegraph Ave       │ 85872391      │  37.83422562297274 │ -122.26336365164984 │ 2017-11-14  │
│ 58a90c051e1de51e677d81e9 │ EZ Taqueria                       │ 4013 Telegraph Ave       │ 85872391      │  37.82956442174899 │ -122.26453047653766 │             │
│ 5334870a498e4600b3af150b │ KOREAN WOOD CHARCOAL BBQ.         │ 4390 Telegraph Ave Ste J │ 85872391      │  37.83205171569861 │ -122.26316650636586 │             │
│ 590003af32b61d706e013649 │ Bunaburger                        │ 4901 Telegraph Ave       │ 85872391      │  37.83594921355949 │ -122.26303233878318 │ 2018-12-04  │
│ 574b5f00498e40ec765fba16 │ Azit                              │ 4390 Telegraph Ave       │ 85872391      │  37.83201381681889 │ -122.26323610358418 │             │
│ 4a1b0a82f964a520c27a1fe3 │ Koryo Ja Jang                     │ 4390 Telegraph Ave       │ 85872391      │  37.83219705852519 │ -122.26354821183742 │ 2023-11-25  │
...
│ 5b207871a92d980039bf5bc8 │ Hancook                           │ 4315 Telegraph Ave       │ 85872391      │  37.83181924223889 │ -122.26408627115096 │             │
│ 4acfef69f964a520f9d620e3 │ Chef Yu - Yuyu Za Zang            │ 4871 Telegraph Ave       │ 85872391      │  37.83540380302543 │ -122.26297108197923 │             │
│ 6681fe8a198ecd288f03d609 │ Small Change Oyster Bar           │ 5000 Telegraph ave       │ 85872391      │          37.836433 │         -122.262268 │             │
│ 534dcb8f11d216d6de8875bf │ Koryo Kalbi                       │ 4390 Telegraph Ave Ste J │ 85872391      │  37.83210754394531 │ -122.26309967041016 │             │
├──────────────────────────┴───────────────────────────────────┴──────────────────────────┴───────────────┴────────────────────┴─────────────────────┴─────────────┤
│ 93 rows (40 shown)                                                                                                                                     7 columns │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### area-whosonfirst-properties

Derive Who's On First properties for an "area" parquet file (produced by by the `compile-area` tool).

```
$> ./bin/area-whosonfirst-properties -h
Derive Who's On First properties for an "area" parquet file (produced by by the `compile-area` tool).
Usage:
	 ./bin/area-whosonfirst-properties [options]
  -area-parquet compile-area
    	The URI for the "area" parquet file (produced by by the compile-area tool) from which Who's On First properties will be derived.
  -reader-uri string
    	A registered whosonfirst/go-reader.Reader URI. (default "https://data.whosonfirst.org")
  -verbose
    	Enable verbose (debug) logging.
  -whosonfirst-parquet string
    	The URI for the parquet file where Who's On First properties will be written to.
  -with-spatial-geom
    	Store geometry property as spatial GEOMETRY type (rather than TEXT.	
```

For example:

```
$> area-whosonfirst-properties/main.go \
	-area-parquet sfba.parquet \
	-whosonfirst-parquet whosonfirst.parquet
```

And then:

```
$> duckdb
v1.1.3 19864453f7
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.

D SELECT id, name, geometry FROM read_parquet('whosonfirst.parquet') LIMIT 1;
┌──────────┬─────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│    id    │      name       │                                                                                  geometry                                                                                   │
│  int32   │     varchar     │                                                                                   varchar                                                                                   │
├──────────┼─────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ 85872355 │ Piedmont Avenue │ {"type":"MultiPolygon","coordinates":[[[[-122.243389,37.830281],[-122.24419,37.829481],[-122.244891,37.82878],[-122.245192,37.828481],[-122.245421,37.828318],[-122.24589…  │
└──────────┴─────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

Note: The default behaviour is to store geoemtry information as a JSON-encoded GeoJSON `geometry`. If you need or want native DuckDB spatial types run the tool with the `-with-spatial-geom` flag.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://opensource.foursquare.com/os-places/
* https://docs.overturemaps.org/guides/places/