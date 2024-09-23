package ugrsical

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cxz66666/zju-ical/pkg/zjuam"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/zjuconst"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	kAppServiceLoginUrl           = "https://zjuam.zju.edu.cn/cas/login?service=http%3A%2F%2Fappservice.zju.edu.cn%2F"
	kAppServiceGetWeekClassInfo   = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kbdy_cxXsZKbxx"
	kAppServiceGetExamOutlineInfo = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kkqk_cxXsksxx"
	kAppServiceGetClassScore      = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kkqk_cxXscjxx"
)

type UgrsService struct {
	ZJUClient zjuam.ZjuLogin
	ctx       context.Context
}

func (zs *UgrsService) Login(username, password string) error {
	if zs.ZJUClient == nil {
		zs.ZJUClient = zjuam.NewClient()
	}

	return zs.ZJUClient.Login(zs.ctx, kAppServiceLoginUrl, username, password)
}

func (zs *UgrsService) GetClassTimeTable(academicYear string, term zjuconst.ClassTerm, stuId string) ([]zjuconst.ZJUClass, error) {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", zjuconst.UgrsClassTermToQueryString(term))
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
	resp, err := zs.ZJUClient.Client().Do(req)
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

	res := make([]zjuconst.ZJUClass, 0)
	for _, item := range classTimeTable.Data.ClassList {
		tmp := item.ToZJUClass()
		if tmp != nil {
			res = append(res, *tmp)
		}
	}
	return res, nil
}

func (zs *UgrsService) GetExamInfo(academicYear string, term zjuconst.ExamTerm, stuId string) ([]zjuconst.ZJUExam, error) {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", zjuconst.UgrsExamTermToQueryString(term))
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
	resp, err := zs.ZJUClient.Client().Do(req)
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

	var res []zjuconst.ZJUExam
	for index, _ := range examOutlines.Data.ExamOutlineList {
		res = append(res, examOutlines.Data.ExamOutlineList[index].ToZJUExam())
	}
	return res, nil
}

func (zs *UgrsService) GetScoreInfo(stuId string) ([]zjuconst.ZJUClassScore, error) {
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
	resp, err := zs.ZJUClient.Client().Do(req)
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

func (zs *UgrsService) GetTermConfigs() []zjuconst.TermConfig {
	var res []zjuconst.TermConfig
	for _, item := range zs.GetCtxConfig().TermConfigs {
		res = append(res, item.ToTermConfig())
	}
	return res
}

func (zs *UgrsService) GetTweaks() []zjuconst.Tweak {
	var res []zjuconst.Tweak
	for _, item := range zs.GetCtxConfig().Tweaks {
		res = append(res, item.ToTweak())
	}
	return res
}

func (zs *UgrsService) GetClassTerms() []zjuconst.ClassYearAndTerm {
	return zs.GetCtxConfig().GetClassYearAndTerms()
}

func (zs *UgrsService) GetExamTerms() []zjuconst.ExamYearAndTerm {
	return zs.GetCtxConfig().GetExamYearAndTerms()
}

func (zs *UgrsService) GetCtxConfig() *zjuconst.ZjuScheduleConfig {
	if config := zs.ctx.Value(zjuconst.ScheduleCtxKey); config != nil {
		return config.(*zjuconst.ZjuScheduleConfig)
	} else {
		return zjuconst.GetConfig()
	}
}

func (zs *UgrsService) UpdateConfig() bool {
	//TODO
	//如此设计我们是否需要update？
	//或者后台开个协程每分钟和fs同步？
	return true
}

func NewUgrsService(ctx context.Context) *UgrsService {
	return &UgrsService{
		ZJUClient: zjuam.NewClient(),
		ctx:       log.With().Str("reqid", uuid.NewString()).Logger().WithContext(ctx),
	}
}
