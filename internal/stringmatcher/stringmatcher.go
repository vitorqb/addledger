package stringmatcher

import (
	"github.com/adrg/strutil/metrics"
)

//go:generate $MOCKGEN --source=stringmatcher.go --destination=../../mocks/stringmatcher/stringmatcher_mock.go

// Options contains the options for a new StringMatcher
type Options struct {
	Cache *Cache
}

// Cache implements a cache for the distance between two strings
type Cache struct {
	data map[string]map[string]int
}

func (c Cache) Set(a string, b string, distance int) {
	if _, ok := c.data[a]; !ok {
		c.data[a] = map[string]int{}
	}
	c.data[a][b] = distance
}

func (c Cache) Get(a string, b string) (distance int, success bool) {
	if distFromA, ok := c.data[a]; ok {
		if value, ok := distFromA[b]; ok {
			return value, true
		}
	}
	if distFromB, ok := c.data[b]; ok {
		if value, ok := distFromB[a]; ok {
			return value, true
		}
	}
	return 0, false
}

// IStringMatcher is an interface for a matching with cache between wo strings
type IStringMatcher interface {
	Distance(a string, b string) int
}

// StringMatcher implements IStringMatcher
type StringMatcher struct {
	cache *Cache
}

var _ IStringMatcher = &StringMatcher{}

// Distance implements IStringMatcher.
func (s *StringMatcher) Distance(a string, b string) int {
	if distance, ok := s.cache.Get(a, b); ok {
		return distance
	}
	distance := metrics.NewLevenshtein().Distance(a, b)
	s.cache.Set(a, b, distance)
	return distance
}

// New creates a new StringMatcher
func New(options *Options) (*StringMatcher, error) {
	cache := options.Cache
	if cache == nil {
		cache = NewCache()
	}
	return &StringMatcher{cache}, nil
}

// NewCache creates a new Cache
func NewCache() *Cache {
	return &Cache{map[string]map[string]int{}}
}
