# go-whosonfirst-external

Go package for working with external data sources in a Who's On First context.

## Tools

### iterate

For example:

```
$> go run cmd/iterate/main.go -iterator-uri foursquare:// ~/data/foursquare/parquet/*.parquet > /dev/null
2024/12/24 18:07:50 INFO Time to iterate records count=104529230 time=1m40.40977525s
```

```
$> go run cmd/iterate/main.go -iterator-uri overture:// ~/data/overture/parquet/*.parquet > /dev/null
2024/12/24 18:04:48 INFO Time to iterate records count=55527168 time=2m11.528865417s
```
