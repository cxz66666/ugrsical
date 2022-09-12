package zjuservice

import (
	"context"
	"encoding/json"
	"fmt"
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
)

type IZjuService interface {
	Login(username, password string) error
	GetClassTimeTable(academicYear string, term ClassTerm, stuId string) []ZjuClass
	GetExamInfo(academicYear string, term ExamTerm, stuId string) []ZjuExamOutline
	GetTermConfigs() []TermConfig
	GetTweaks() []Tweak
	GetClassTerms() []ClassYearAndTerm
	GetExamTerms() []ExamYearAndTerm
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

func (zs *ZjuService) GetClassTimeTable(academicYear string, term ClassTerm, stuId string) []ZjuClass {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", ClassTermToQueryString(term))
	data.Set("xh", stuId)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", kAppServiceGetWeekClassInfo, strings.NewReader(encodedData))
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("new request failed")
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO check
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("POST to Class API failed")
		return nil
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(content)

	classTimeTable := ZjuResWrapperStr[ZjuWeeklyScheduleRes]{}
	if err = json.Unmarshal(content, &classTimeTable); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("unmarshal failed")
		return nil
	}

	res := make([]ZjuClass, 0)
	for _, item := range classTimeTable.Data.ClassList {
		tmp := item.ToZjuClass()
		if tmp != nil {
			res = append(res, *tmp)
		}
	}
	return res
}

func (zs *ZjuService) GetExamInfo(academicYear string, term ExamTerm, stuId string) []ZjuExamOutline {
	data := url.Values{}
	data.Set("xn", academicYear)
	data.Set("xq", ExamTermToQueryString(term))
	data.Set("xh", stuId)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", kAppServiceGetExamOutlineInfo, strings.NewReader(encodedData))
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("new request failed")
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO check
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("POST to Class API failed")
		return nil
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(content)
	examOutlines := ZjuResWrapperStr[ZjuExamOutlineRes]{}
	if err = json.Unmarshal(content, &examOutlines); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msg("unmarshal failed")
		return nil
	}

	return examOutlines.Data.ExamOutlineList
}

func (zs *ZjuService) GetTermConfigs() []TermConfig {
	var res []TermConfig
	for _, item := range _schedule.TermConfigs {
		res = append(res, item.ToTermConfig())
	}
	return res
}

func (zs *ZjuService) GetTweaks() []Tweak {
	var res []Tweak
	for _, item := range _schedule.Tweaks {
		res = append(res, item.ToTweak())
	}
	return res
}

func (zs *ZjuService) GetClassTerms() []ClassYearAndTerm {
	return _schedule.GetClassYearAndTerms()
}

func (zs *ZjuService) GetExamTerms() []ExamYearAndTerm {
	return _schedule.GetExamYearAndTerms()
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
