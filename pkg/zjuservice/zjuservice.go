package zjuservice

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"ugrs-ical/pkg/zjuam"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	kAppServiceLoginUrl           = "https://zjuam.zju.edu.cn/cas/login?service=http%3A%2F%2Fappservice.zju.edu.cn%2Findex"
	kAppServiceGetWeekClassInfo   = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kbdy_cxXsZKbxx"
	kAppServiceGetExamOutlineInfo = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kkqk_cxXsksxx"
	kAppServiceGetClassScore      = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kkqk_cxXscjxx"
)

type IZjuService interface {
	Login(username, password string) error
	GetClassTimeTable(academicYear string, term ClassTerm, stuId string) ([]ZjuClass, error)
	GetExamInfo(academicYear string, term ExamTerm, stuId string) ([]ZjuExamOutline, error)
	GetScoreInfo(stuId string) ([]ZjuClassScore, error)
	GetTermConfigs() []TermConfig
	GetTweaks() []Tweak
	GetClassTerms() []ClassYearAndTerm
	GetExamTerms() []ExamYearAndTerm
	GetCtxConfig() *ZjuScheduleConfig
	UpdateConfig() bool
}

type ZjuService struct {
	ZjuClient zjuam.ZjuLogin
	ctx       context.Context
}

func (zs *ZjuService) Login(username, password string) error {
	if zs.ZjuClient == nil {
		zs.ZjuClient = zjuam.NewClient()
	}

	return zs.ZjuClient.Login(context.Background(), kAppServiceLoginUrl, username, password)
}

func (zs *ZjuService) GetClassTimeTable(academicYear string, term ClassTerm, stuId string) ([]ZjuClass, error) {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", ClassTermToQueryString(term))
	data.Set("xh", stuId)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", kAppServiceGetWeekClassInfo, strings.NewReader(encodedData))
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("new request failed")
		return nil, errors.New("new request failed")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO check
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("POST to Class API failed")
		return nil, errors.New("POST to Class API failed")
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	classTimeTable := ZjuResWrapperStr[ZjuWeeklyScheduleRes]{}
	if err = json.Unmarshal(content, &classTimeTable); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msgf("unmarshal failed %s", stuId)
		return nil, errors.New("无法获取课表, 大概率为浙大钉钉服务端问题，请打开浙大钉中的个人课表功能排查，若浙大钉课表无法打开/没有内容，则确认是浙大钉问题，请稍后重试，否则可能是学号密码错误")
	}

	res := make([]ZjuClass, 0)
	for _, item := range classTimeTable.Data.ClassList {
		tmp := item.ToZjuClass()
		if tmp != nil {
			res = append(res, *tmp)
		}
	}
	return res, nil
}

func (zs *ZjuService) GetExamInfo(academicYear string, term ExamTerm, stuId string) ([]ZjuExamOutline, error) {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", ExamTermToQueryString(term))
	data.Set("xh", stuId)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", kAppServiceGetExamOutlineInfo, strings.NewReader(encodedData))
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("new request failed")
		return nil, errors.New("new request failed")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO check
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("POST to Class API failed")
		return nil, errors.New("POST to Class API failed")
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	examOutlines := ZjuResWrapperStr[ZjuExamOutlineRes]{}
	if err = json.Unmarshal(content, &examOutlines); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msgf("unmarshal failed %s", stuId)
		return nil, errors.New("无法获取课表, 大概率为浙大钉钉服务端问题，请打开浙大钉中的个人课表功能排查，若浙大钉课表无法打开/没有内容，则确认是浙大钉问题，请稍后重试，否则可能是学号密码错误")
	}

	return examOutlines.Data.ExamOutlineList, nil
}

func (zs *ZjuService) GetScoreInfo(stuId string) ([]ZjuClassScore, error) {
	data := url.Values{}
	data.Set("lx", "0")
	data.Set("xh", stuId)
	data.Set("xn", "")
	data.Set("xq", "")
	data.Set("cjd", "")

	encodedData := data.Encode()
	req, err := http.NewRequest("POST", kAppServiceGetClassScore, strings.NewReader(encodedData))
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("new request failed")
		return nil, errors.New("new request failed")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO check
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("POST to Class API failed")
		return nil, errors.New("POST to Class API failed")
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	classScores := ZjuResWrapperStr[ZjuClassScoreRes]{}
	if err = json.Unmarshal(content, &classScores); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msgf("unmarshal failed %s", stuId)
		return nil, errors.New("无法获取课表, 大概率为浙大钉钉服务端问题，请打开浙大钉中的个人课表功能排查，若浙大钉课表无法打开/没有内容，则确认是浙大钉问题，请稍后重试，否则可能是学号密码错误")
	}

	return classScores.Data.ClassScoreList, nil
}

func (zs *ZjuService) GetTermConfigs() []TermConfig {
	var res []TermConfig
	for _, item := range zs.GetCtxConfig().TermConfigs {
		res = append(res, item.ToTermConfig())
	}
	return res
}

func (zs *ZjuService) GetTweaks() []Tweak {
	var res []Tweak
	for _, item := range zs.GetCtxConfig().Tweaks {
		res = append(res, item.ToTweak())
	}
	return res
}

func (zs *ZjuService) GetClassTerms() []ClassYearAndTerm {
	return zs.GetCtxConfig().GetClassYearAndTerms()
}

func (zs *ZjuService) GetExamTerms() []ExamYearAndTerm {
	return zs.GetCtxConfig().GetExamYearAndTerms()
}

func (zs *ZjuService) GetCtxConfig() *ZjuScheduleConfig {
	if config := zs.ctx.Value(ScheduleCtxKey); config != nil {
		return config.(*ZjuScheduleConfig)
	} else {
		ScheduleRwMutex.RLock()
		defer ScheduleRwMutex.RUnlock()
		return &_schedule
	}
}

func (zs *ZjuService) UpdateConfig() bool {
	//TODO
	//如此设计我们是否需要update？
	//或者后台开个协程每分钟和fs同步？
	return true
}

func NewZjuService(ctx context.Context) *ZjuService {
	return &ZjuService{
		ZjuClient: zjuam.NewClient(),
		ctx:       log.With().Str("reqid", uuid.NewString()).Logger().WithContext(ctx),
	}
}
