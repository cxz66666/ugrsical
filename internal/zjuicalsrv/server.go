package zjuicalsrv

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"text/template"
	"time"
	"zju-ical/pkg/zjuam"
	"zju-ical/pkg/zjuservice/zjuconst"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

const defaultServerConfigPath = "configs/server.json"

var _serverConfig ServerConfig

type SetupData struct {
	Classes          []zjuconst.YearAndSemester
	Exams            []zjuconst.YearAndSemester
	LastUpdated      int
	LastUpdatedTime  string
	LastSuccessIcal  string
	LastSuccessScore string
	Link             string
	SubLink          string
	ScoreSubLink     string
}

type ServerConfig struct {
	Enckey    string `json:"enckey"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	CfgPath   string `json:"config"`
	HttpProxy string `json:"http_proxy"`
	IpHeader  string `json:"ip_header"`
	RedisAddr string `json:"redis_addr"`
	RedisPass string `json:"redis_pass"`
	CacheTTL  int    `json:"cache_ttl"`
}

var setupTpl *template.Template

var sd = SetupData{
	Classes:          []zjuconst.YearAndSemester{},
	Exams:            []zjuconst.YearAndSemester{},
	Link:             "",
	SubLink:          "",
	ScoreSubLink:     "",
	LastSuccessIcal:  "暂无",
	LastSuccessScore: "暂无",
}
var sdMutex sync.RWMutex

var gcm cipher.AEAD
var rc *redis.Client
var cacheTTL time.Duration

func loadServerConfig() error {
	var r io.Reader
	f, err := os.Open(defaultServerConfigPath)
	defer f.Close()
	r = f
	if err != nil {
		return err
	}
	cfd, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(cfd, &_serverConfig)
	return err
}

func validConfig() error {
	if len(_serverConfig.Enckey) == 0 || len(_serverConfig.Host) == 0 || _serverConfig.Port == 0 {
		return errors.New("invalid server config file, please check enckey,host and port fields")
	}
	return nil
}

func ListenAndServe() error {
	if err := loadServerConfig(); err != nil {
		return err
	}
	if err := validConfig(); err != nil {
		return err
	}
	c, err := aes.NewCipher([]byte(_serverConfig.Enckey))
	if err != nil {
		return err
	}
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return err
	}

	if _serverConfig.IpHeader != "" {
		log.Info().Msgf("zjuicalsrv will get header from %s", _serverConfig.IpHeader)
	}

	if _serverConfig.RedisAddr == "" {
		log.Warn().Msg("redis not set, rate limit won't work")
	} else {
		rc = redis.NewClient(&redis.Options{
			Addr:     _serverConfig.RedisAddr,
			Password: _serverConfig.RedisPass,
			DB:       0,
		})
		pong, err := rc.Ping(context.Background()).Result()
		if err != nil {
			log.Error().Msgf("redis ping error: %s", err)
			return err
		}
		log.Info().Msgf("redis ping: %s", pong)
	}

	if _serverConfig.CacheTTL == 0 {
		cacheTTL = DurationIcalCache
	} else if _serverConfig.CacheTTL < 0 {
		return errors.New("cache ttl must be positive")
	} else {
		cacheTTL = time.Duration(_serverConfig.CacheTTL) * time.Hour
	}
	log.Info().Msgf("cache ttl: %f", cacheTTL.Hours())

	if len(_serverConfig.HttpProxy) > 0 {
		go zjuam.UpdateProxyUrl(time.Minute*10, _serverConfig.HttpProxy)
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
	if err = zjuconst.LoadConfig(_serverConfig.CfgPath); err != nil {
		return err
	}

	if zjuconst.UseOnlineConfig {
		go zjuconst.UpdateConfig(time.Hour * 1)
	}

	//TODO check config
	//当前生成的学期
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	setRoutes(app)

	log.Info().Msgf("[server] running on %d", _serverConfig.Port)
	return app.Run(fmt.Sprintf(":%d", _serverConfig.Port))

}

func setRoutes(app *gin.Engine) {
	app.Use(RateLimiterM())
	app.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static")
	})
	app.GET("/ical", FetchCal)
	app.GET("/sub", SubCal)
	app.GET("/subScore", SubScore)
	app.GET("/score", FetchScore)
	app.GET("/ping", PingEp)
	app.POST("/static", SetupPage)
	app.Static("/static", "./web/app")

}
