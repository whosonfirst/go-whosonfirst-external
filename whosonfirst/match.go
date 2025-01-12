package whosonfirst

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"
)

var ancestors_cache = new(sync.Map)

const mode_all string = "all"
const mode_any string = "any"

type RowHasMatchOptions struct {
	Geohash     string
	ParentIds   []int64
	AncestorIds []int64
	Mode        string
}

func RowHasMatch(ctx context.Context, row map[string]string, opts *RowHasMatchOptions) (bool, error) {

	has_match := true

	if opts.Geohash != "" {

		if strings.HasPrefix(row["geohash"], opts.Geohash) {
			has_match = true
		} else {

			has_match = false

			if opts.Mode != mode_any {
				return has_match, nil
			}
		}
	}

	if len(opts.ParentIds) > 0 {

		if row["wof:parent_id"] == "" {

			has_match = false

			if opts.Mode != mode_any {
				return has_match, nil
			}

		} else {

			parent_id, err := strconv.ParseInt(row["wof:parent_id"], 10, 64)

			if err != nil {

				slog.Warn("Failed to parse parent ID, skipping", "parent id", row["wof:parent_id"], "error", err)

				has_match = false

				if opts.Mode != mode_any {
					return has_match, nil
				}

			} else {

				if slices.Contains(opts.ParentIds, parent_id) {
					has_match = true
				} else {

					has_match = false

					if opts.Mode != mode_any {
						return has_match, nil
					}
				}
			}
		}

	}

	if len(opts.AncestorIds) > 0 {

		if row["wof:hierarchies"] == "" {

			has_match = false

			if opts.Mode != mode_any {
				return has_match, nil
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
						return has_match, nil
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
					return has_match, nil
				}
			}
		}
	}

	return has_match, nil
}
