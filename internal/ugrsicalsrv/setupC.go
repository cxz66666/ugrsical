package ugrsicalsrv

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"ugrs-ical/pkg/zjuservice"

	"github.com/gin-gonic/gin"
)

func encrypt(b []byte) ([]byte, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte(""), err
	}
	d := gcm.Seal(nonce, nonce, b, nil)
	return d, nil
}

func SetupPage(ctx *gin.Context) {
	u, _ := ctx.GetPostForm("username")
	p, _ := ctx.GetPostForm("password")
	if u == "" || p == "" {
		ctx.String(http.StatusOK, "用户名或密码未输入")
		return
	}
	uP := bytes.Repeat([]byte("#"), 12)
	l := 12
	if len(u) < 12 {
		l = len(u)
	}
	for i := 0; i < l; i++ {
		uP[i] = u[i]
	}

	b := append(uP, []byte(p)...)
	b, err := encrypt(b)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	en := base64.URLEncoding.EncodeToString(b)

	d := sd
	d.Classes = zjuservice.GetConfig().GetClassYearAndSemester()
	d.Exams = zjuservice.GetConfig().GetExamYearAndSemester()
	d.LastUpdated = zjuservice.GetConfig().GetLastUpdated()
	d.LastUpdatedTime = zjuservice.GetConfig().GetLastUpdatedTime()
	d.Link = fmt.Sprintf("%s/ical?p=%s", _serverConfig.Host, en)
	d.SubLink = fmt.Sprintf("%s/sub?p=%s", _serverConfig.Host, en)
	d.ScoreSubLink = fmt.Sprintf("%s/subScore?p=%s", _serverConfig.Host, en)

	ctx.Header("Content-Type", "text/html")
	buffer := bytes.NewBuffer([]byte(""))
	err = setupTpl.Execute(buffer, d)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.String(http.StatusOK, buffer.String())
	return
}
