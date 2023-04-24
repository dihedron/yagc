package cache

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestCacheJSON(t *testing.T) {

	files := []string{
		"./test1.json",
		"./test2.json",
	}

	log := slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr))

	cache := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[0]}),
		WithEncoding[string, string](&JSON[string, string]{Pretty: true}),
		WithPolicy[string, string](&Always{}),
	)

	// put some values
	v, ok := cache.Replace("a", "aaa")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("b", "bbb")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("c", "ccc")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("d", "ddd")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	ok = cache.Put("e", "eee")
	assert.Equal(t, ok, true, "The value should not be present in the cache.")

	// now try to put it again and check that iot is NOT overwritten
	ok = cache.Put("e", "xxx")
	assert.Equal(t, ok, false, "The value should not have been set in the cache.")
	v, ok = cache.Get("e")
	assert.Equal(t, v, "eee", "The value should not have been overwritten.")
	assert.Equal(t, ok, true, "The value should not have been reported as present in the cache.")

	// create a replica and load it from file
	exec.Command("cp", "-rf", files[0], files[1]).Run()
	cache2 := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[1]}),
		WithEncoding[string, string](&JSON[string, string]{Pretty: true}),
		WithPolicy[string, string](&Never{}), // do not persist to disk!!!
	)
	cache2.Load()
	keys := cache2.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")
	for _, k := range keys {
		v, ok := cache2.Get(k)
		assert.Equal(t, ok, true, "The value should be present in the cache.")
		assert.Equal(t, v, k+k+k, "The value should be as expected.")
	}
	// now remove a key and check that is does not store it
	cache2.Delete("c")
	cache2.Store()

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
	keys = cache.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")

	// remove a key
	v, ok = cache.Delete("c")
	assert.Equal(t, ok, true, "The value should have been present in the cache.")
	assert.Equal(t, v, "ccc", "The value is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 4, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "d", "e"}, "The key set is invalid.")

	// clear the cache
	cache.Clear()

	// check size
	i = cache.Size()
	assert.Equal(t, i, 0, "The cache size is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 0, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{}, "The key set is invalid.")

}

func TestCacheGOB(t *testing.T) {

	files := []string{
		"./test1.gob",
		"./test2.gob",
	}

	log := slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr))

	cache := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[0]}),
		WithEncoding[string, string](&GOB[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)

	// put some values
	v, ok := cache.Replace("a", "aaa")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("b", "bbb")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("c", "ccc")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("d", "ddd")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	ok = cache.Put("e", "eee")
	assert.Equal(t, ok, true, "The value should not be present in the cache.")

	// now try to put it again and check that it is NOT overwritten
	ok = cache.Put("e", "xxx")
	assert.Equal(t, ok, false, "The value should not have been set in the cache.")
	v, ok = cache.Get("e")
	assert.Equal(t, v, "eee", "The value should not have been overwritten.")
	assert.Equal(t, ok, true, "The value should not have been reported as present in the cache.")

	// create a replica and load it from file
	exec.Command("cp", "-rf", files[0], files[1]).Run()
	cache2 := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[1]}),
		WithEncoding[string, string](&GOB[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)
	cache2.Load()
	keys := cache2.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")
	for _, k := range keys {
		v, ok := cache2.Get(k)
		assert.Equal(t, ok, true, "The value should be present in the cache.")
		assert.Equal(t, v, k+k+k, "The value should be as expected.")
	}

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
	keys = cache.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")

	// remove a key
	v, ok = cache.Delete("c")
	assert.Equal(t, ok, true, "The value should have been present in the cache.")
	assert.Equal(t, v, "ccc", "The value is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 4, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "d", "e"}, "The key set is invalid.")

	// clear the cache
	cache.Clear()

	// check size
	i = cache.Size()
	assert.Equal(t, i, 0, "The cache size is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 0, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{}, "The key set is invalid.")
}

func TestCacheYAML(t *testing.T) {

	files := []string{
		"./test1.yaml",
		"./test2.yaml",
	}

	log := slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr))

	cache := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[0]}),
		WithEncoding[string, string](&YAML[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)

	// put some values
	v, ok := cache.Replace("a", "aaa")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("b", "bbb")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("c", "ccc")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	v, ok = cache.Replace("d", "ddd")
	assert.Equal(t, ok, false, "The value should not be present in the cache.")
	assert.Equal(t, v, "", "The previous value should be empty.")

	ok = cache.Put("e", "eee")
	assert.Equal(t, ok, true, "The value should not be present in the cache.")

	// now try to put it again and check that iot is NOT overwritten
	ok = cache.Put("e", "xxx")
	assert.Equal(t, ok, false, "The value should not have been set in the cache.")
	v, ok = cache.Get("e")
	assert.Equal(t, v, "eee", "The value should not have been overwritten.")
	assert.Equal(t, ok, true, "The value should not have been reported as present in the cache.")

	// create a replica and load it from file
	exec.Command("cp", "-rf", files[0], files[1]).Run()
	cache2 := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: files[1]}),
		WithEncoding[string, string](&YAML[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)
	cache2.Load()
	keys := cache2.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")
	for _, k := range keys {
		v, ok := cache2.Get(k)
		assert.Equal(t, ok, true, "The value should be present in the cache.")
		assert.Equal(t, v, k+k+k, "The value should be as expected.")
	}

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
	keys = cache.Keys()
	assert.Equal(t, len(keys), 5, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "c", "d", "e"}, "The key set is invalid.")

	// remove a key
	v, ok = cache.Delete("c")
	assert.Equal(t, ok, true, "The value should have been present in the cache.")
	assert.Equal(t, v, "ccc", "The value is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 4, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{"a", "b", "d", "e"}, "The key set is invalid.")

	// clear the cache
	cache.Clear()

	// check size
	i = cache.Size()
	assert.Equal(t, i, 0, "The cache size is invalid.")

	// now check the keys again
	keys = cache.Keys()
	assert.Equal(t, len(keys), 0, "The number of keys is invalid.")
	assert.ElementsMatch(t, keys, []string{}, "The key set is invalid.")
}
