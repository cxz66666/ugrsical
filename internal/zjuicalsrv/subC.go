package zjuicalsrv

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	common2 "zju-ical/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func SubCal(ctx *gin.Context) {
	p := ctx.Query("p")
	if p == "" {
		ctx.String(http.StatusOK, "invalid p")
		return
	}
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
	subUrl := fmt.Sprintf("webcal://%s/ical?p=%s", _serverConfig.Host, p)

	c := log.With().Str("u", string(un)).Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	if rc != nil {
		val, err := rc.Exists(c, genIcalKey(string(un), string(pw))).Result()
		if val == 0 {
			log.Ctx(c).Info().Msgf("don't find ical cache")
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache failed")
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find ical cache, redirect sub url")
			ctx.Redirect(http.StatusFound, subUrl)
			return
		}
	}
	vCal, err := common2.GetBothCalendar(c, string(un), string(pw))

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	sdMutex.Lock()
	sd.LastSuccessIcal = time.Now().Format("2006.01.02 15:04:05")
	sdMutex.Unlock()

	if rc != nil {
		err = rc.Set(c, genIcalKey(string(un), string(pw)), []byte(vCal.GetICS("")), cacheTTL).Err()
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("set ical cache failed, url = %s", "/ical?p="+p)
		} else {
			log.Ctx(c).Info().Msgf("set ical cache success, url = %s", "/ical?p="+p)
		}
	}
	ctx.Redirect(http.StatusFound, subUrl)
	return
}

func SubScore(ctx *gin.Context) {
	p := ctx.Query("p")
	if p == "" {
		ctx.String(http.StatusOK, "invalid p")
		return
	}
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
	subUrl := fmt.Sprintf("webcal://%s/score?p=%s", _serverConfig.Host, p)

	c := log.With().Str("u", string(un)).Str("type", "score").Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	if rc != nil {
		val, err := rc.Exists(c, genScoreKey(string(un), string(pw))).Result()
		if val == 0 {
			log.Ctx(c).Info().Msgf("don't find score cache")
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache failed")
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find score cache, redirect sub url")
			ctx.Redirect(http.StatusFound, subUrl)
			return
		}
	}
	vCal, err := common2.GetScoreCalendar(c, string(un), string(pw))

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
	ctx.Redirect(http.StatusFound, subUrl)
	return
}
