package zjuicalsrv

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	common2 "github.com/cxz66666/zju-ical/internal/common"
	"github.com/cxz66666/zju-ical/pkg/ical"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const DurationIcalCache = time.Hour * 48

func decryptWithGCM(b []byte, usedGCM cipher.AEAD) ([]byte, error) {
	ns := usedGCM.NonceSize()
	if len(b) < ns {
		return []byte(""), errors.New("invalid data")
	}
	nonce, ct := b[:ns], b[ns:]
	p, err := usedGCM.Open(nil, nonce, ct, nil)
	if err != nil {
		return []byte(""), err
	}
	return p, nil
}

func decrypt(b []byte) ([]byte, error) {
	p, err := decryptWithGCM(b, gcm)
	if err == nil {
		return p, nil
	}
	if gcm2 != nil {
		p, err = decryptWithGCM(b, gcm2)
		if err == nil {
			return p, nil
		}
	}
	return []byte(""), err
}

// genIcalKey return key for redis ical data, format "319010xxxx" + "###" + sha256(passwd)[0:16] + "change"
func genIcalKey(username, passwd, change string) string {
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
	if change != "" {
		b = append(b, []byte("change"+change)...)
	}
	return string(b)
}

func FetchCal(ctx *gin.Context) {
	p := ctx.Query("p")
	if p == "" {
		ctx.String(http.StatusOK, "invalid p")
		return
	}

	exam := ctx.Query("exam")
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

	c := log.With().Str("u", string(un)).Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	if rc != nil {
		data, err := rc.Get(c, genIcalKey(string(un), string(pw)+exam, change)).Bytes()
		if err == redis.Nil {
			log.Ctx(c).Info().Msgf("don't find ical cache")
		} else if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("fetch cache failed")
			ctx.String(http.StatusOK, "redis 内部错误，请查看日志")
			return
		} else {
			//get cache
			log.Ctx(c).Info().Msgf("find ical cache")
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
	var vCal ical.VCalendar
	if exam == "0" {
		vCal, err = common2.GetClassCalendar(c, string(un), string(pw), isGRS)
	} else {
		vCal, err = common2.GetBothCalendar(c, string(un), string(pw), isGRS)
	}
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	sdMutex.Lock()
	sd.LastSuccessIcal = time.Now().Format("2006.01.02 15:04:05")
	sdMutex.Unlock()

	if rc != nil {
		err = rc.Set(c, genIcalKey(string(un), string(pw)+exam, change), []byte(vCal.GetICS("")), cacheTTL).Err()
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("set ical cache failed, url = %s", "/ical?p="+p)
		} else {
			log.Ctx(c).Info().Msgf("set ical cache success, url = %s", "/ical?p="+p)
		}
	}
	ctx.Header("Content-Type", "text/calendar")
	ctx.Data(http.StatusOK, "text/calendar", []byte(vCal.GetICS("")))
	return
}
