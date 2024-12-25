# go-reader-database-sql

[database/sql](https://golang.org/pkg/database/sql/) support for the go-reader Reader interface. 

## Important

Work in progress. Documentation to follow.

## Example

```
package main

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/whosonfirst/go-reader"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"	
	sql_reader "github.com/whosonfirst/go-reader-database-sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	
	ctx := context.Background()

	// You can also enable this specific "read function" by passing
	// ?parse-uri=true to the reader URI below.
	
	sql_reader.URI_READFUNC = func(uri string) (string, error) {
		id, _ := wof_uri.IdFromPath(uri)		
		str_id := strconv.FormatInt(id, 10)
		return str_id, nil
	}

	uri := "sql://sqlite3/geojson/id/body?dsn=fr.db"	
	r, _ := reader.NewReader(ctx, uri)

	fh, _ := r.Read(ctx, "102/065/003/102065003.geojson")
	defer fh.Close()
	
	io.Copy(os.Stdout, fh)		
}
```

## See also

* https://github.com/whosonfirst/go-reader
* https://golang.org/pkg/database/sql/
