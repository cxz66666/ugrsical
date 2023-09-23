package zjuam

import (
	"context"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

var httpProxyUrl, _ = url.Parse("")
var proxyUrlMutex sync.RWMutex

type ZjuLogin interface {
	Login(ctx context.Context, payloadUrl, username, password string) error
	Client() *http.Client
}

type ZjuamClient struct {
	HttpClient *http.Client
}

func NewClient() *ZjuamClient {
	jar, _ := cookiejar.New(nil)
	proxyUrlMutex.RLock()
	defer proxyUrlMutex.RUnlock()
	if len(httpProxyUrl.String()) == 0 {
		return &ZjuamClient{
			HttpClient: &http.Client{Jar: jar, Transport: &http.Transport{Proxy: nil}},
		}
	}
	return &ZjuamClient{
		HttpClient: &http.Client{Jar: jar, Transport: &http.Transport{Proxy: http.ProxyURL(httpProxyUrl)}},
	}
}

func testProxyUrl(httpProxyStr string) error {
	proxyUrl, err := url.Parse(httpProxyStr)
	if err != nil {
		return err
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, err := client.Get("https://baidu.com")
	resp.Body.Close()
	return err
}

func UpdateProxyUrl(interval time.Duration, httpProxyStr string) {
	if err := testProxyUrl(httpProxyStr); err == nil {
		log.Info().Msgf("[server] use http proxy %s", httpProxyStr)
		proxyUrlMutex.Lock()
		httpProxyUrl, _ = url.Parse(httpProxyStr)
		proxyUrlMutex.Unlock()
	} else {
		log.Info().Msgf("[server] test http proxy fail %s, %s", time.Now().Format("2006.01.02 15:04:05"), err)
	}

	log.Info().Msgf("[server] update http proxy every %s", interval.String())
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			proxyUrlMutex.Lock()
			if err := testProxyUrl(httpProxyStr); err != nil {
				log.Info().Msgf("[server] test http proxy fail %s, %s", time.Now().Format("2006.01.02 15:04:05"), err)
				httpProxyUrl, _ = url.Parse("")
			} else {
				httpProxyUrl, _ = url.Parse(httpProxyStr)
			}
			proxyUrlMutex.Unlock()
		}
	}
}
