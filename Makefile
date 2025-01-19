GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	@make modvendor
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iterate cmd/iterate/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/assign-ancestors cmd/assign-ancestors/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/sort-ancestors cmd/sort-ancestors/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/walk-sorted cmd/walk-sorted/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/compile-area cmd/compile-area/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/area-whosonfirst-properties cmd/area-whosonfirst-properties/main.go

# https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring
modvendor:
	modvendor -copy="**/*.a **/*.h" -v
