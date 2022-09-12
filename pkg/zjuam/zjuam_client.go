package zjuam

import (
	"context"
	"net/http"
	"net/http/cookiejar"
)

type ZjuLogin interface {
	Login(ctx context.Context, payloadUrl, username, password string) error
	Client() *http.Client
}

type ZjuamClient struct {
	HttpClient *http.Client
}

func NewClient() *ZjuamClient {
	jar, _ := cookiejar.New(nil)
	return &ZjuamClient{
		HttpClient: &http.Client{Jar: jar},
	}
}
