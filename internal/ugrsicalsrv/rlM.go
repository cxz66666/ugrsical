package ugrsicalsrv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	nLimit = 20
	nTime  = 1
	pLimit = 20
	pTime  = 6
)

func getIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if h := ctx.Request.Header.Get(_serverConfig.IpHeader); _serverConfig.IpHeader != "" && h != "" {
		ip = h
	}
	return ip
}

func RateLimiterM() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rc != nil {
			path := ctx.Request.URL.Path
			region := "n"
			t := nTime
			l := nLimit
			if path == "/ical" {
				region = "p"
				t = pTime
				l = pLimit
			}

			key := fmt.Sprintf("%s%s", region, getIP(ctx))
			counter, err := rc.Get(ctx, key).Int64()
			if err == redis.Nil {
				err = rc.Set(ctx, key, 1, time.Duration(t)*time.Minute).Err()
				if err != nil {
					ctx.AbortWithError(500, err)
					return
				}
				counter = 1
			} else if err != nil {
				ctx.AbortWithError(500, err)
				return
			} else {
				if counter > int64(l) {
					ctx.String(http.StatusOK, "limit reached for you ip.")
					ctx.Abort()
					return
				}
				_ = rc.Incr(ctx, key)
			}
		}
		ctx.Next()
	}
}
