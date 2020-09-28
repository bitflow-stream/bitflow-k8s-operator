package bitflow

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ModificationCache struct {
	client           client.Client
	operationTimeout time.Duration
	lock             sync.Mutex

	deletions cacheMap
	creations cacheMap
	updates   cacheMap
}

// objectKind -> name -> timestamp
type cacheMap map[string]map[string]time.Time

func NewModificationCache(cl client.Client, operationTimeout time.Duration) *ModificationCache {
	return &ModificationCache{
		client:           cl,
		operationTimeout: operationTimeout,
		deletions:        make(cacheMap),
		creations:        make(cacheMap),
		updates:          make(cacheMap),
	}
}

func (c *ModificationCache) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()

	reset := func(cache cacheMap) {
		for kind := range cache {
			delete(cache, kind)
		}
	}
	reset(c.creations)
	reset(c.deletions)
	reset(c.updates)
}

func (c *ModificationCache) Create(obj runtime.Object, objKind string, name string) (bool, error) {
	return c.performModification(c.creations, []cacheMap{c.deletions, c.updates}, objKind, name, func() error {
		return c.client.Create(context.TODO(), obj)
	})
}

func (c *ModificationCache) Delete(obj runtime.Object, objKind string, name string, opts ...client.DeleteOptionFunc) (bool, error) {
	return c.performModification(c.deletions, []cacheMap{c.creations, c.updates}, objKind, name, func() error {
		return c.client.Delete(context.TODO(), obj, opts...)
	})
}

func (c *ModificationCache) Update(obj runtime.Object, objKind string, name string) (bool, error) {
	return c.performModification(c.updates, []cacheMap{c.creations, c.deletions}, objKind, name, func() error {
		return c.client.Update(context.TODO(), obj)
	})
}

func (c *ModificationCache) performModification(cache cacheMap, otherCaches []cacheMap, objKind string, name string, modification func() error) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	kindCache := cache[objKind]
	if kindCache == nil {
		kindCache = make(map[string]time.Time)
		cache[objKind] = kindCache
	}

	operationTimestamp := kindCache[name]
	timeSinceOperation := now.Sub(operationTimestamp)
	if c.operationTimeout <= 0 || operationTimestamp.IsZero() || timeSinceOperation > c.operationTimeout {
		// Operation not yet performed or timed out
		if err := modification(); err != nil {
			// Error: do not refresh the operation timestamp
			return true, err
		} else {
			// Operation succeeded: store the timestamp of the operation and remove entry from other caches
			kindCache[name] = now
			for _, otherCache := range otherCaches {
				otherKindCache := otherCache[objKind]
				if _, ok := otherKindCache[name]; ok {
					delete(otherKindCache, name)
				}
			}
			return true, nil
		}
	} else {
		// Operation is repeated, ignore for now, until it times out
		log.WithFields(log.Fields{"kind": objKind, "name": name}).Debugf("Ignoring modification operation, because it occurred %v ago")
		return false, nil
	}
}
