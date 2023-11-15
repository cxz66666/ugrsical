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
var proxyTransport http.Transport

type ZjuLogin interface {
	Login(ctx context.Context, payloadUrl, username, password string) error
	// UgrsExtraLogin used for ugrs students select grs course
	UgrsExtraLogin(ctx context.Context, payloadUrl string) error
	Client() *http.Client
}

type ZjuamClient struct {
	HttpClient *http.Client
}

func NewClient() *ZjuamClient {
	jar, _ := cookiejar.New(nil)
	proxyUrlMutex.RLock()
	defer proxyUrlMutex.RUnlock()
	return &ZjuamClient{
		HttpClient: &http.Client{Jar: jar, Transport: &proxyTransport},
	}
}

func testProxyUrl(httpProxyStr string) error {
	proxyUrl, err := url.Parse(httpProxyStr)
	if err != nil {
		return err
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), DisableKeepAlives: true}}
	if resp, err := client.Get("https://baidu.com"); err != nil {
		return err
	} else {
		resp.Body.Close()
		return nil
	}
}

func UpdateProxyUrl(interval time.Duration, httpProxyStr string) {
	if err := testProxyUrl(httpProxyStr); err == nil {
		log.Info().Msgf("[server] use http proxy %s", httpProxyStr)
		proxyUrlMutex.Lock()
		httpProxyUrl, _ = url.Parse(httpProxyStr)
		proxyTransport = http.Transport{Proxy: http.ProxyURL(httpProxyUrl)}
		proxyUrlMutex.Unlock()
	} else {
		log.Info().Msgf("[server] test http proxy fail %s, %s", time.Now().Format("2006.01.02 15:04:05"), err)
		proxyUrlMutex.Lock()
		proxyTransport = http.Transport{Proxy: nil}
		proxyUrlMutex.Unlock()

	}

	log.Info().Msgf("[server] update http proxy every %s", interval.String())
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := testProxyUrl(httpProxyStr); err != nil {
				proxyUrlMutex.Lock()
				log.Info().Msgf("[server] test http proxy fail %s, %s", time.Now().Format("2006.01.02 15:04:05"), err)
				httpProxyUrl, _ = url.Parse("")
				proxyTransport.CloseIdleConnections()
				proxyTransport = http.Transport{Proxy: nil}
			} else {
				proxyUrlMutex.Lock()
				httpProxyUrl, _ = url.Parse(httpProxyStr)
				if proxyTransport.Proxy == nil {
					proxyTransport = http.Transport{Proxy: http.ProxyURL(httpProxyUrl)}
				}
			}
			proxyUrlMutex.Unlock()
		}
	}
}
