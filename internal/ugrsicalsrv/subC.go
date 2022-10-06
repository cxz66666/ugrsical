package ugrsicalsrv

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	common2 "ugrs-ical/internal/common"

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
			log.Ctx(c).Info().Msgf("don't find cache with id %s, will login and fetch", genIcalKey(string(un), string(pw)))
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache with id %s failed", genIcalKey(string(un), string(pw)))
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find cache with id %s, redirect sub url", genIcalKey(string(un), string(pw)))
			ctx.Redirect(http.StatusFound, subUrl)
			return
		}
	}
	vCal, err := common2.GetBothCalendar(c, string(un), string(pw))

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	if rc != nil {
		err = rc.Set(c, genIcalKey(string(un), string(pw)), []byte(vCal.GetICS("")), cacheTTL).Err()
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("set cache with id %s failed", genIcalKey(string(un), string(pw)))
		} else {
			log.Ctx(c).Info().Msgf("set cache with id %s success", genIcalKey(string(un), string(pw)))
		}
	}
	ctx.Redirect(http.StatusFound, subUrl)
	return
}
