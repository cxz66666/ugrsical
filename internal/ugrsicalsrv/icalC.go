package ugrsicalsrv

import (
	"context"
	"encoding/base64"
	"errors"

	common2 "ugrs-ical/internal/common"

	"github.com/gofiber/fiber/v2"
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

func FetchCal(ctx *fiber.Ctx) error {
	p := ctx.Query("p")
	if p == "" {
		return ctx.SendString("invalid p")
	}
	b, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		return ctx.SendString("invalid p2")
	}
	unpw, err := decrypt(b)
	if err != nil {
		return ctx.SendString("invalid p3")
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
		return ctx.SendString(err.Error())
	}

	ctx.Set("Content-Type", "text/calendar")
	return ctx.SendString(vCal.GetICS(""))
}
