package zjuicalsrv

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	common2 "github.com/cxz66666/zju-ical/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const DurationScoreCache = time.Hour * 6

// genScoreKey return key for redis score data, format "319010xxxx" + "***" + sha256(passwd)[0:16]
func genScoreKey(username, passwd string) string {
	uP := bytes.Repeat([]byte("*"), 12)
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

func FetchScore(ctx *gin.Context) {
	p := ctx.Query("p")
	if p == "" {
		ctx.String(http.StatusOK, "invalid p")
		return
	}
	change := ctx.Query("change")
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

	c := log.With().Str("u", string(un)).Str("type", "score").Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	if rc != nil {
		data, err := rc.Get(c, genScoreKey(string(un), string(pw))).Bytes()
		if err == redis.Nil {
			log.Ctx(c).Info().Msgf("don't find score cache")
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache failed")
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find score cache")
			ctx.Header("Content-Type", "text/calendar")
			ctx.Data(http.StatusOK, "text/calendar", data)
			return
		}
	}
	var isGRS bool
	if strings.HasPrefix(string(un), "3") {
		isGRS = false
	} else {
		isGRS = true
	}
	if change == "1" {
		isGRS = !isGRS
	}
	vCal, err := common2.GetScoreCalendar(c, string(un), string(pw), isGRS)

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	sdMutex.Lock()
	sd.LastSuccessScore = time.Now().Format("2006.01.02 15:04:05")
	sdMutex.Unlock()

	if rc != nil {
		err = rc.Set(c, genScoreKey(string(un), string(pw)), []byte(vCal.GetICS("ZJU-ICAL GPA表")), DurationScoreCache).Err()
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("set score cache failed, url = %s", "/score?p="+p)
		} else {
			log.Ctx(c).Info().Msgf("set score cache success, url = %s", "/score?p="+p)
		}
	}
	ctx.Header("Content-Type", "text/calendar")
	ctx.Data(http.StatusOK, "text/calendar", []byte(vCal.GetICS("ZJU-ICAL GPA表")))
	return
}
