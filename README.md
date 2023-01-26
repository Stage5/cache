# Cache
This is a simple, thread-safe, in-memory cache library written in Go. It allows for the storage of key-value pairs with a defined time-to-live (TTL) value, and regularly cleans up expired values.

## Features
* Insertion of key-value pairs with a defined TTL
* Updating of key-value pairs with a defined TTL
* Retrieval of key-value pairs
* Deletion of key-value pairs
* Retrieval of current cache size

## Usage
```
import (
    "time"
    "github.com/stage5/cache"
)

cacheTTL := 30 * time.Second
cleanupInterval := 10 * time.Second

lc := localcache.NewLocalCache(cacheTTL, cleanupInterval)

// Insert a value into the cache
lc.InsertOne("key1", "value1")

// Update a value in the cache
lc.UpdateOne("key1", "newvalue1")

// Get a value from the cache
value, _ := lc.GetOne("key1")

// Delete a value from the cache  
lc.DeleteOne("key1")

// Get the current cache size
size := lc.GetSize()
```

## Errors
The following errors may be returned by the library:

* `ErrKeyExists`: Returned when attempting to insert a key that already exists in the cache.
* `ErrKeyNotInCache`: Returned when attempting to update or retrieve a key that does not exist in the cache.

## Cleanup
Cache cleanup is performed in a separate goroutine, using the `cleanupLoop()` function. The cleanup interval can be set when initializing the cache with the `NewLocalCache()` function.

## Thread Safety
This library is designed to be thread-safe, using Go's built-in `sync.RWMutex` for safe read and write access to the cache.

## Note
This cache has no limit on maximum size.

##  Closing
To stop the cleanup goroutine, use the `stop` channel in the LocalCache struct.
```
Copy code
lc.stop <- struct{}{}
lc.wg.Wait()
```
