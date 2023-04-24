# YAGC - Yet another simple Golang cache.

This is a simple Golang cache.

It adds synchronisation and persistence to Golang maps; it can use any `comparable` key type and any type of value.

The persistence support revolves around three concepts:

. *policy*: when to persist data (e.g. whenever updated, never, in batches...)
. *encoding*: in which format to persist it (e.g. YAML, JSON, GOB...)
. *persistence*: where to persist it (e.g. file, S3, database)

## How to use it

Look at the unit tests.

```golang
	cache := New(
		WithLogger[string, string](log),
		WithPersistence[string, string](&File{Path: "./test.gob"}),
		WithEncoding[string, string](&GOB[string, string]{}),
		WithPolicy[string, string](&Always{}),
	)
    var (
        ok bool
        v string
    )
    // put some value in; ok says whether it was actually put
    // because Put() does not replace and existing value
    ok = cache.Put("a", "aaa")
	// if you want to force-put it, use Replace() instead
	v, ok = cache.Replace("a", "aaa")
    // now get the value back
    v, ok = cache.Get("a")

    // to store the cache, call the method
    cache.Store()


```
