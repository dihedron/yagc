package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestCachePutGetDeleteSizeKeys(t *testing.T) {

	log := slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr))

	cache := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: "./test.json"}),
		WithEncoding[string, string](&JSON[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)

	// put some values
	v, ok := cache.Put("a", "aaa")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Put("b", "bbb")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Put("c", "ccc")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Put("d", "ddd")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Put("e", "eee")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	// get some values
	v, ok = cache.Get("a")
	assert.Equal(t, ok, true, "The value should be present in the cache.")
	assert.Equal(t, v, "aaa", "The value should be as expected.")

	v, ok = cache.Get("<not present>")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The value should be empty.")

	// check size
	i := cache.Size()
	assert.Equal(t, i, 5, "The cache size is invalid.")

	// get cache keys
	k := cache.Keys()
	assert.Equal(t, len(k), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, k, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")

	// remove a key
	v, ok = cache.Delete("c")
	assert.Equal(t, ok, true, "The value should have been present in the cache.")
	assert.Equal(t, v, "ccc", "The value is invalid.")

	// now check the keys again
	k = cache.Keys()
	assert.Equal(t, len(k), 4, "The number of keys is invalid.")
	assert.ElementsMatch(t, k, []string{"a", "b", "d", "e"}, "The key set is invalid.")

	// clear the cache
	cache.Clear()

	// check size
	i = cache.Size()
	assert.Equal(t, i, 0, "The cache size is invalid.")

	// now check the keys again
	k = cache.Keys()
	assert.Equal(t, len(k), 0, "The number of keys is invalid.")
	assert.ElementsMatch(t, k, []string{}, "The key set is invalid.")

}
