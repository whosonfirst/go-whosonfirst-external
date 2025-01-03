package external

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Iterator interface {
	Iterate(context.Context, ...string) iter.Seq2[Record, error]
	Close() error
}

var iterator_roster roster.Roster

// IteratorInitializationFunc is a function defined by individual iterator package and used to create
// an instance of that iterator
type IteratorInitializationFunc func(ctx context.Context, uri string) (Iterator, error)

// RegisterIterator registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Iterator` instances by the `NewIterator` method.
func RegisterIterator(ctx context.Context, scheme string, init_func IteratorInitializationFunc) error {

	err := ensureIteratorRoster()

	if err != nil {
		return err
	}

	return iterator_roster.Register(ctx, scheme, init_func)
}

func ensureIteratorRoster() error {

	if iterator_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		iterator_roster = r
	}

	return nil
}

// NewIterator returns a new `Iterator` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `IteratorInitializationFunc`
// function used to instantiate the new `Iterator`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterIterator` method.
func NewIterator(ctx context.Context, uri string) (Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := iterator_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(IteratorInitializationFunc)
	return init_func(ctx, uri)
}

// IteratorSchemes returns the list of schemes that have been registered.
func IteratorSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureIteratorRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range iterator_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
