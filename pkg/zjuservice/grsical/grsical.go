package grsical

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cxz66666/zju-ical/pkg/zjuam"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/zjuconst"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
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
		err := zs.ZJUClient.Login(zs.ctx, ugrsLoginUrl, username, password)
		if err != nil {
			return err
		}
		return zs.ZJUClient.UgrsExtraLogin(zs.ctx, fmt.Sprintf("%s?gnmkdm=N253530&su=%s", ugrsLoginUrl2, username))
	}
	return zs.ZJUClient.Login(zs.ctx, grsLoginUrl, username, password)
}

func (zs *GrsService) fetchTimetable(academicYear string, term zjuconst.ClassTerm) (io.Reader, error) {
	var changeLocaleUrl string
	var fetchUrl string

	semester := zjuconst.GrsClassTermToClassQueryInt(term)
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
func (zs *GrsService) fetchExamtable(academicYear string, term zjuconst.ExamTerm) (io.Reader, error) {
	var fetchUrl string
	semester := zjuconst.GrsExamTermToQueryInt(term)
	if zs.isUGRS {
		fetchUrl = fmt.Sprintf("http://grs.zju.edu.cn/py/undergraduate/grksap.htm?xj=%d&xn=%d", semester, academicYear)
	} else {
		fetchUrl = fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grksap.htm?xj=%d&xn=%d", semester, academicYear)
	}
	r, err := zs.ZJUClient.Client().Get(fetchUrl)
	if err != nil {
		e := fmt.Sprintf("failed to fetch exam for %d-%d, error: %s", academicYear, semester, err.Error())
		log.Ctx(zs.ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read exam for %d-%d, error: %s", academicYear, semester, err.Error())
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
	table, err := GetTable(r)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("get class failed")
		return nil, nil
	}

	return ParseTable(zs.ctx, table, zs.isUGRS)
}

func (zs *GrsService) GetExamInfo(academicYear string, term zjuconst.ExamTerm, stuId string) ([]zjuconst.ZJUExam, error) {
	r, err := zs.fetchExamtable(academicYear, term)
	if err != nil {
		return nil, err
	}
	table, err := GetExamTable(r)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("get exam table failed")
		return nil, nil
	}
	return ParseExamTable(zs.ctx, table)
}

func (zs *GrsService) GetScoreInfo(stuId string) ([]zjuconst.ZJUClassScore, error) {
	// TODO
	return nil, errors.New("研究生版成绩解析模块正在开发中，敬请期待QAQ")
}

func (zs *GrsService) GetTermConfigs() []zjuconst.TermConfig {
	var res []zjuconst.TermConfig
	for _, item := range zs.GetCtxConfig().TermConfigs {
		res = append(res, item.ToTermConfig())
	}
	return res
}

func (zs *GrsService) GetTweaks() []zjuconst.Tweak {
	var res []zjuconst.Tweak
	for _, item := range zs.GetCtxConfig().Tweaks {
		res = append(res, item.ToTweak())
	}
	return res
}

func (zs *GrsService) GetClassTerms() []zjuconst.ClassYearAndTerm {
	return zs.GetCtxConfig().GetClassYearAndTerms()
}

func (zs *GrsService) GetExamTerms() []zjuconst.ExamYearAndTerm {
	return zs.GetCtxConfig().GetExamYearAndTerms()
}

func (zs *GrsService) GetCtxConfig() *zjuconst.ZjuScheduleConfig {
	if config := zs.ctx.Value(zjuconst.ScheduleCtxKey); config != nil {
		return config.(*zjuconst.ZjuScheduleConfig)
	} else {
		return zjuconst.GetConfig()
	}
}

func (zs *GrsService) UpdateConfig() bool {
	//TODO
	//如此设计我们是否需要update？
	//或者后台开个协程每分钟和fs同步？
	return true
}

func NewGrsService(ctx context.Context, isUgrs bool) *GrsService {
	return &GrsService{
		ZJUClient: zjuam.NewClient(),
		ctx:       log.With().Str("reqid", uuid.NewString()).Logger().WithContext(ctx),
		isUGRS:    isUgrs,
	}
}
