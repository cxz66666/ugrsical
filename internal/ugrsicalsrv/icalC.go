package ugrsicalsrv

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	common2 "ugrs-ical/internal/common"
	"ugrs-ical/pkg/ical"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const DurationIcalCache = time.Hour * 48

func decrypt(b []byte) ([]byte, error) {
	ns := gcm.NonceSize()
	if len(b) < ns {
		return []byte(""), errors.New("invalid data")
	}
	nonce, ct := b[:ns], b[ns:]
	p, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return []byte(""), err
	}
	return p, nil
}

// genIcalKey return key for redis ical data, format "319010xxxx" + "###" + sha256(passwd)[0:16]
func genIcalKey(username, passwd string) string {
	uP := bytes.Repeat([]byte("#"), 12)
	l := 12
	if len(username) < 12 {
		l = len(username)
	}
	for i := 0; i < l; i++ {
		uP[i] = username[i]
	}

	hashPasswd := sha256.Sum256([]byte(passwd))
	b := append(uP, hashPasswd[0:16]...)

	return string(b)
}

func FetchCal(ctx *gin.Context) {
	p := ctx.Query("p")
	if p == "" {
		ctx.String(http.StatusOK, "invalid p")
		return
	}

	exam := ctx.Query("exam")

	b, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		ctx.String(http.StatusOK, "invalid p2")
		return
	}
	unpw, err := decrypt(b)
	if err != nil {
		ctx.String(http.StatusOK, "invalid p2")
		return
	}
	un := unpw[:12]
	pw := unpw[12:]
	for i := 11; i >= 0; i-- {
		if un[i] != '#' {
			un = un[:i+1]
			break
		}
	}

	c := log.With().Str("u", string(un)).Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	if rc != nil {
		data, err := rc.Get(c, genIcalKey(string(un), string(pw)+exam)).Bytes()
		if err == redis.Nil {
			log.Ctx(c).Info().Msgf("don't find cache with id %s, will login and fetch", genIcalKey(string(un), string(pw)+exam))
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache with id %s failed", genIcalKey(string(un), string(pw)+exam))
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find cache with id %s, return data", genIcalKey(string(un), string(pw)+exam))
			ctx.Header("Content-Type", "text/calendar")
			ctx.Data(http.StatusOK, "text/calendar", data)
			return
		}
	}
	var vCal ical.VCalendar
	if exam == "0" {
		vCal, err = common2.GetClassCalendar(c, string(un), string(pw))
	} else {
		vCal, err = common2.GetBothCalendar(c, string(un), string(pw))
	}
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	if rc != nil {
		err = rc.Set(c, genIcalKey(string(un), string(pw)+exam), []byte(vCal.GetICS("")), cacheTTL).Err()
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("set cache with id %s failed", genIcalKey(string(un), string(pw)+exam))
		} else {
			log.Ctx(c).Info().Msgf("set cache with id %s success", genIcalKey(string(un), string(pw)+exam))
		}
	}
	ctx.Header("Content-Type", "text/calendar")
	ctx.Data(http.StatusOK, "text/calendar", []byte(vCal.GetICS("")))
	return
}
