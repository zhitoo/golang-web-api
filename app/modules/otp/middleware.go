package otp

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/config"
	"golang.org/x/time/rate"
)

func init() {
	if cfg == nil {
		cfg = config.GetConfig()
	}
}

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	ipLimiters = make(map[string]*entry)
	mu         sync.Mutex
)

func getLimiter(store map[string]*entry, key string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	e, ok := store[key]
	if !ok {
		// 1 request per minute
		if cfg.App.Env == "local" {
			e = &entry{limiter: rate.NewLimiter(rate.Every(time.Second), 1)}
		} else {
			e = &entry{limiter: rate.NewLimiter(rate.Every(time.Minute), 1)}
		}
		store[key] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

func LimitOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !getLimiter(ipLimiters, ip).Allow() {
			response.NewResponse().SetError(fmt.Errorf("too many requests from %s", ip)).SetHttpStatusCode(http.StatusTooManyRequests).Json(c)
			return
		}

		c.Next()
	}
}
