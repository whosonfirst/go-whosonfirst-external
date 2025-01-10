package walk

import (
	"compress/bzip2"
	"crypto/sha256"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	
	"github.com/sfomuseum/go-csvdict/v2"
)

const mode_all string = "all"
const mode_any string = "any"

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	ancestors_cache := new(sync.Map)
	
	var csv_wr *csvdict.Writer
	
	for _, path := range opts.Sources {

		to_walk := os.DirFS(path)

		for row, err := range walkFS(ctx, to_walk) {

			if err != nil {
				return err
			}

			has_match := true
			
			if opts.GeoHash != "" {

				if strings.HasPrefix(row["geohash"], opts.GeoHash) {
					has_match = true
				} else {
					
					has_match = false

					if opts.Mode != mode_any {
						continue
					}
				}
			}

			if len(opts.ParentIds) > 0 {

				if row["wof:parent_id"] == "" {

					has_match = false

					if opts.Mode != mode_any {
						continue
					}
					
				} else {

					parent_id, err := strconv.ParseInt(row["wof:parent_id"], 10, 64)
					
					if err != nil {
						
						slog.Warn("Failed to parse parent ID, skipping", "parent id", row["wof:parent_id"], "error", err)

						has_match = false
						
						if opts.Mode != mode_any {
							continue
						}
						
					} else {
					
						if slices.Contains(opts.ParentIds, parent_id) {
							has_match = true
						} else {
							
							has_match = false
							
							if opts.Mode != mode_any {
								continue
							}					
						}
					}
				}
				
			}

			if len(opts.AncestorIds) > 0 {

				if row["wof:hierarchies"] == "" {
					
					has_match = false
					
					if opts.Mode != mode_any {
						continue
					}					
					
				} else {

					// START OF cache me...

					hierarchies_sum := sha256.Sum256([]byte(row["wof:hierarchies"]))
					hierarchies_k := fmt.Sprintf("%s", hierarchies_sum)

					possible := make([]int64, 0)
					
					v, ok := ancestors_cache.Load(hierarchies_k)

					if ok {
						possible = v.([]int64)
					} else {
						
						var hierarchies []map[string]int64
						
						err := json.Unmarshal([]byte(row["wof:hierarchies"]), &hierarchies)
						
						if err != nil {
							slog.Warn("Failed to unmarshal hierarchies, skipping", "error", err)
							
							has_match = false
							
							if opts.Mode != mode_any {
								continue
							}					
							
						} else {
							
							for _, h := range hierarchies {
								
								for _, id := range h {
									
									if !slices.Contains(possible, id) {
										possible = append(possible, id)
									}
								}
							}
						}

						ancestors_cache.Store(hierarchies_k, possible)
					}
					
					// END OF cache me...
						
					has_ancestor := false
					
					for _, id := range possible {
						
						if slices.Contains(opts.AncestorIds, id) {
								has_ancestor = true
								break
							}
						}
						
						if has_ancestor {
							has_match = true
						} else {
							
							has_match = false
							
							if opts.Mode != mode_any {
								continue
							}					
						}
				}
			}
			
			if !has_match {
				continue
			}
			
			if csv_wr == nil {
				
				wr, err := csvdict.NewWriter(os.Stdout)
				
				if err != nil {
					return fmt.Errorf("Failed to create new CSV writer, %w", err)
				}
				
				csv_wr = wr
			}
			
			csv_wr.WriteRow(row)
			csv_wr.Flush()
		}		
	}

	return nil	
}

func walkFS(ctx context.Context, to_walk fs.FS) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		err := fs.WalkDir(to_walk, ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)

			if ext != ".bz2" {
				return nil
			}

			r, err := to_walk.Open(path)

			if err != nil {
				return err
			}

			defer r.Close()

			br := bzip2.NewReader(r)

			csv_r, err := csvdict.NewReader(br)

			if err != nil {
				return err
			}

			for row, err := range csv_r.Iterate() {

				if err != nil {
					return err
				}

				if !yield(row, nil) {
					break
				}
			}

			return nil
		})

		if err != nil {
			yield(nil, err)
			return
		}
	}
}
