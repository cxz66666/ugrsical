package grsicalsrv

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"text/template"

	cmn "grs-ical/internal/common"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type SetupData struct {
	Year     int
	YearN    int
	Semester string
	Link     string
}

var EncKey, Host string
var setupTpl *template.Template
var cfg cmn.Config
var tweak cmn.TweakConfig
var sd = SetupData{
	Year:     0,
	YearN:    0,
	Semester: "",
	Link:     "",
}
var gcm cipher.AEAD
var rc *redis.Client
var ipHeader string

func ListenAndServe(address string) error {
	EncKey = os.Getenv("GRSICALSRV_ENCKEY")
	if EncKey == "" {
		return errors.New("encryption key not set")
	}
	c, err := aes.NewCipher([]byte(EncKey))
	if err != nil {
		return err
	}
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return err
	}

	Host = os.Getenv("GRSICALSRV_HOST")
	if Host == "" {
		return errors.New("host not set")
	}
	configFile := os.Getenv("GRSICALSRV_CFG")
	if configFile == "" {
		return errors.New("config file not set")
	}
	tweaksFile := os.Getenv("GRSICALSRV_TWEAKS")
	if tweaksFile == "" {
		return errors.New("config file not set")
	}
	ipHeader = os.Getenv("GRSICALSRV_IP_HEADER")
	if ipHeader != "" {
		log.Info().Msgf("ugrsicalsrv will get header from %s", ipHeader)
	}
	redisAddr := os.Getenv("GRSICALSRV_REDIS_ADDR")
	redisPass := os.Getenv("GRSICALSRV_REDIS_PASS")
	if redisAddr == "" || redisPass == "" {
		log.Warn().Msg("redis not set, rate limit won't work")
	} else {
		rc = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPass,
			DB:       0,
		})
	}
	// read template
	f, err := os.Open("./web/template/setup.html")
	if err != nil {
		return err
	}
	fc, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	setupTpl, err = template.New("setup").Parse(string(fc))
	if err != nil {
		return err
	}
	// read config
	cfg, err = cmn.GetConfig(configFile)
	if err != nil {
		return err
	}
	sd.Year = cfg.FetchConfig[0].Year
	sd.YearN = sd.Year + 1
	if cfg.FetchConfig[0].Semester > 12 {
		sd.Semester = "秋冬"
	} else {
		sd.Semester = "春夏"
	}
	tweak, err = cmn.GetTweakConfig(tweaksFile)
	if err != nil {
		return err
	}

	app := fiber.New()
	setRoutes(app)
	return app.Listen(address)
}

func setRoutes(app *fiber.App) {
	app.Use(RateLimiterM)
	app.Static("/", "./web/app")
	app.Post("/", SetupPage)
	app.Get("/ical", FetchCal)
	app.Get("/ping", PingEp)
}
