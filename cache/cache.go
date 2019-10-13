package cache

import "time"

// Cache should be able to set and get keys to improve performance
type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}
