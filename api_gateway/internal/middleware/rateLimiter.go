package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BucketLimiter struct {
	BucketCh chan struct{}
}

func NewBucketLimiter(ctx context.Context, capacity int, period time.Duration) *BucketLimiter {
	limiter := &BucketLimiter{
		BucketCh: make(chan struct{}, capacity),
	}

	for i := 0; i < capacity; i++ {
		limiter.BucketCh <- struct{}{}
	}

	replenishmentTiming := period.Nanoseconds() / int64(capacity)

	go limiter.ReplenishBucket(ctx, time.Duration(replenishmentTiming))

	return limiter

}

func (b *BucketLimiter) ReplenishBucket(ctx context.Context, replenishmentTiming time.Duration) {
	ticker := time.NewTicker(replenishmentTiming)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.BucketCh <- struct{}{}
		}
	}
}

func (b *BucketLimiter) Allow() bool {
	select {
	case <-b.BucketCh:
		return true
	default:
		return false
	}
}

func RateLimitMiddleware(limiter *BucketLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		allow := limiter.Allow()
		if !allow {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}
		c.Next()
	}
}
