package grsical

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"ugrs-ical/pkg/zjuam"
	"ugrs-ical/pkg/zjuservice/zjuconst"
)

const (
	grsLoginUrl         = "http://grs.zju.edu.cn/ssohome"
	ugrsLoginUrl        = "http://zdbk.zju.edu.cn/jwglxt/xtgl/login_ssologin.html"
	ugrsLoginUrl2       = "http://zdbk.zju.edu.cn/jwglxt/xsxk/zzxkghb_cxZzxkGhbJumpYjsCourses.html"
	grcChangeLocaleUrl  = "http://grs.zju.edu.cn/py/page/student/grkcb.htm?pageAction=changeLocale"
	ugrsChangeLocaleUrl = "http://grs.zju.edu.cn/py/undergraduate/grkcb.htm?pageAction=changeLocale"
)

type GrsService struct {
	ZJUClient zjuam.ZjuLogin
	ctx       context.Context
	isUGRS    bool
}

func (zs *GrsService) Login(username, password string) error {
	if zs.ZJUClient == nil {
		zs.ZJUClient = zjuam.NewClient()
	}
	if zs.isUGRS {
		err := zs.ZJUClient.Login(context.Background(), ugrsLoginUrl, username, password)
		if err != nil {
			return err
		}
		return zs.ZJUClient.UgrsExtraLogin(context.Background(), fmt.Sprintf("%s?gnmkdm=N253530&su=%s", ugrsLoginUrl2, username))
	}
	return zs.ZJUClient.Login(context.Background(), grsLoginUrl, username, password)
}

func (zs *GrsService) fetchTimetable(academicYear string, term zjuconst.ClassTerm) (io.Reader, error) {
	var changeLocaleUrl string
	var fetchUrl string

	semester := zjuconst.GrsClassTermToQueryInt(term)
	if zs.isUGRS {
		changeLocaleUrl = ugrsChangeLocaleUrl
		fetchUrl = fmt.Sprintf("http://grs.zju.edu.cn/py/undergraduate/grkcb.htm?xj=%d&xn=%d", semester, academicYear)
	} else {
		changeLocaleUrl = grcChangeLocaleUrl
		fetchUrl = fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grkcb.htm?xj=%d&xn=%d", semester, academicYear)
	}

	_, err := zs.ZJUClient.Client().PostForm(changeLocaleUrl, url.Values{
		"locale": {"zh_CN"},
	})
	r, err := zs.ZJUClient.Client().Get(fetchUrl)
	if err != nil {
		e := fmt.Sprintf("failed to fetch timetable for %d-%d, error: %s", academicYear, semester, err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read timetable for %d-%d, error: %s", academicYear, semester, err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	return bytes.NewBuffer(rb), nil

}
func (zs *GrsService) GetClassTimeTable(academicYear string, term zjuconst.ClassTerm, stuId string) ([]zjuconst.ZJUClass, error) {
	r, err := zs.fetchTimetable(academicYear, term)
	if err != nil {
		return nil, err
	}

}
