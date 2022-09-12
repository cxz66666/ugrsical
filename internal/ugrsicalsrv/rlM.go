package ugrsicalsrv

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

const (
	nLimit = 10
	nTime  = 1
	pLimit = 10
	pTime  = 6
)

func getIP(ctx *fiber.Ctx) string {
	ip := ctx.IP()
	if h, ok := ctx.GetReqHeaders()[_serverConfig.IpHeader]; _serverConfig.IpHeader != "" && ok {
		ip = h
	}
	return ip
}

func RateLimiterM(ctx *fiber.Ctx) error {
	if rc == nil {
		return ctx.Next()
	}
	path := string(ctx.Request().URI().Path())
	region := "n"
	t := nTime
	l := nLimit
	if path == "/ical" {
		region = "p"
		t = pTime
		l = pLimit
	}

	key := fmt.Sprintf("%s%s", region, getIP(ctx))
	c := ctx.UserContext()
	counter, err := rc.Get(c, key).Int64()
	if err == redis.Nil {
		err = rc.Set(c, key, 1, time.Duration(t)*time.Minute).Err()
		if err != nil {
			return err
		}
		counter = 1
	} else if err != nil {
		return err
	} else {
		if counter > int64(l) {
			return ctx.SendString("limit reached for you ip.")
		}
		_ = rc.Incr(c, key)
	}

	return ctx.Next()
}
