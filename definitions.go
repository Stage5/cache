package cache

import (
	"sync"
	"time"
)

type cachedValue struct {
	value             interface{}
	expireAtTimestamp int64
}

type LocalCache struct {
	stop chan struct{}

	wg sync.WaitGroup
	mu sync.RWMutex

	cacheTTL time.Duration

	cache map[string]cachedValue
}

func (lc *LocalCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-lc.stop:
			return
		case <-t.C:
			lc.mu.Lock()
			for uid, cu := range lc.cache {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.cache, uid)
				}
			}
			lc.mu.Unlock()
		}
	}
}

func (lc *LocalCache) stopCleanup() {
	close(lc.stop)
	lc.wg.Wait()
}
