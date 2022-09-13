package ugrsicalsrv

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	common2 "ugrs-ical/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

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

func FetchCal(ctx *gin.Context) {
	p := ctx.Query("p")
	fmt.Println(p)
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

	c := log.With().Str("u", string(un)).Str("r", uuid.NewString()).Logger().WithContext(context.Background())
	vCal, err := common2.GetBothCalendar(c, string(un), string(pw))
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.Header("Content-Type", "text/calendar")
	ctx.Data(http.StatusOK, "text/calendar", []byte(vCal.GetICS("")))
	return
}
