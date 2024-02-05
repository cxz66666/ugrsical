package newgrsical

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxz66666/zju-ical/pkg/zjuam"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/zjuconst"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

const (
	newGrsLoginUrl = "https://zjuam.zju.edu.cn/cas/login?service=https%3A%2F%2Fyjsy.zju.edu.cn%2F"
)

type NewGrsService struct {
	ZJUClient zjuam.ZjuLogin
	ctx       context.Context
	isUGRS    bool
	token     string
}

func (zs *NewGrsService) Login(username, password string) error {
	if zs.ZJUClient == nil {
		zs.ZJUClient = zjuam.NewClient()
	}
	if zs.isUGRS {
		e := fmt.Sprintf("本科生请在订阅页面选择使用本科生版本，新版研究生系统已不支持查看本科生选研究生课功能")
		log.Ctx(zs.ctx).Error().Msg(e)
		return errors.New(e)
	}
	if err := zs.ZJUClient.Login(zs.ctx, newGrsLoginUrl, username, password); err != nil {
		return err
	}
	lRes, err := zs.ZJUClient.Client().Get(newGrsLoginUrl)
	if err != nil {
		e := fmt.Sprintf("can not access login page: %s", err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return errors.New(e)
	}
	lRes.Body.Close()
	ticketID := lRes.Request.URL.RawQuery
	if ticketID == "" {
		log.Ctx(zs.ctx).Error().Msg("登录失败[获取ticketID失败]，请检查用户名密码是否正确")
		return errors.New("登录失败[获取ticketID失败]，请检查用户名密码是否正确")
	}
	validateURL := fmt.Sprintf("https://yjsy.zju.edu.cn/dataapi/sys/cas/client/validateLogin?service=https:%%2F%%2Fyjsy.zju.edu.cn%%2F&%s", ticketID)
	vRes, err := zs.ZJUClient.Client().Get(validateURL)
	if err != nil {
		e := fmt.Sprintf("登录失败[发送验证ticketID请求失败]: %s", err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return errors.New(e)
	}
	tokenWrapperStr, err := io.ReadAll(vRes.Body)
	if err != nil {
		e := fmt.Sprintf("登录失败[读取验证ticketID请求失败]: %s", err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return errors.New(e)
	}
	vRes.Body.Close()
	tokenWrapper := ZjuResWrapperStr[ZjuLoginToken]{}
	if err = json.Unmarshal(tokenWrapperStr, &tokenWrapper); err != nil || len(tokenWrapper.Data.Token) == 0 {
		e := fmt.Sprintf("登录失败[解析验证ticketID请求失败，怀疑是教务网502或限制内网访问]: %s", err)
		log.Ctx(zs.ctx).Error().Msg(e)
		return errors.New(e)
	}
	zs.token = tokenWrapper.Data.Token
	return nil
}

func (zs *NewGrsService) GetClassTimeTable(academicYear string, term zjuconst.ClassTerm, stuId string) ([]zjuconst.ZJUClass, error) {
	var fetchUrl string

	semester := zjuconst.GrsClassTermToClassQueryInt(term)
	academicYearNum, _ := strconv.Atoi(academicYear[:4])
	fetchUrl = fmt.Sprintf("https://yjsy.zju.edu.cn/dataapi/py/pyKcbj/queryXskbByLoginUser?xn=%d&pkxq=%d", academicYearNum, semester)
	fetchReq, _ := http.NewRequest("GET", fetchUrl, nil)
	fetchReq.Header.Set("X-Access-Token", zs.token)

	r, err := zs.ZJUClient.Client().Do(fetchReq)
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
	classTimeTable := ZjuResWrapperStr[ZjuWeeklyScheduleRes]{}
	if err = json.Unmarshal(rb, &classTimeTable); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msgf("failed to marshal class struct %s", err)
		return nil, errors.New("无法获取课表, 大概率为研究生教务系统服务端问题，请打开研究生教务系统中的个人课表功能排查，若新版研究生教务系统无法打开/没有内容，则确认其问题，请稍后重试，否则可能是学号密码错误")
	}

	res := make([]zjuconst.ZJUClass, 0)
	class_set := make(map[string]bool)
	for _, weekdayItemRaw := range classTimeTable.Data {
		var weekdayItem = make(map[string]json.RawMessage)
		err = json.Unmarshal(weekdayItemRaw, &weekdayItem)
		if err != nil {
			continue
		}
		for _, classItemsRaw := range weekdayItem {
			classItems := ZjuWeeklyScheduleItem{}
			err = json.Unmarshal(classItemsRaw, &classItems)
			if err != nil {
				continue
			}
			for _, item := range classItems.Classes {
				tmp := item.ToZJUClass()
				if tmp != nil && class_set[tmp.ClassCode] == false {
					res = append(res, *tmp)
					class_set[tmp.ClassCode] = true
				}
			}
		}
	}
	return res, nil
}

func (zs *NewGrsService) GetExamInfo(academicYear string, term zjuconst.ExamTerm, stuId string) ([]zjuconst.ZJUExam, error) {
	var fetchUrl string

	semester := zjuconst.NewGrsExamTermToQueryInt(term)
	academicYearNum, _ := strconv.Atoi(academicYear[:4])
	fetchUrl = fmt.Sprintf("https://yjsy.zju.edu.cn/dataapi/py/pyKsxsxx/queryPageByXs?dm=py_grks&mode=2&role=1&column=createTime&order=desc&queryMode=1&field=id,,kcbh,kcmc,rq,ksTime,xn,xq_dictText,ksdd,zwh&pageNo=1&pageSize=100&xn=%d&xq=%d", academicYearNum, semester)
	fetchReq, _ := http.NewRequest("GET", fetchUrl, nil)
	fetchReq.Header.Set("X-Access-Token", zs.token)

	r, err := zs.ZJUClient.Client().Do(fetchReq)
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

	examOutlines := ZjuResWrapperStr[ZjuExamOutlineRes]{}
	if err = json.Unmarshal(rb, &examOutlines); err != nil {
		log.Ctx(zs.ctx).Error().Err(err).Msgf("failed to marshal exam struct %s", err)
		return nil, errors.New("无法获取考试, 大概率为研究生教务系统服务端问题，请打开研究生教务系统中的考试功能排查，若新版研究生教务系统无法打开/没有内容，则确认其问题，请稍后重试，否则可能是学号密码错误")
	}

	var res []zjuconst.ZJUExam
	for _, item := range examOutlines.Data.ExamOutlineList {
		tmp := item.ToZJUExam()
		if tmp != nil {
			res = append(res, tmp)
		} else {
			log.Ctx(zs.ctx).Warn().Msgf("failed to parse exam: %v", item)
			return nil, errors.New(fmt.Sprintf("解析%s考试时间失败，具体解析失败考试时间为：%d", item.ClassName, item.ExamDay))
		}
	}
	return res, nil
}

func (zs *NewGrsService) GetScoreInfo(stuId string) ([]zjuconst.ZJUClassScore, error) {
	// TODO
	return nil, errors.New("研究生版成绩解析模块正在开发中，敬请期待QAQ")
}

func (zs *NewGrsService) GetTermConfigs() []zjuconst.TermConfig {
	var res []zjuconst.TermConfig
	for _, item := range zs.GetCtxConfig().TermConfigs {
		res = append(res, item.ToTermConfig())
	}
	return res
}

func (zs *NewGrsService) GetTweaks() []zjuconst.Tweak {
	var res []zjuconst.Tweak
	for _, item := range zs.GetCtxConfig().Tweaks {
		res = append(res, item.ToTweak())
	}
	return res
}

func (zs *NewGrsService) GetClassTerms() []zjuconst.ClassYearAndTerm {
	return zs.GetCtxConfig().GetClassYearAndTerms()
}

func (zs *NewGrsService) GetExamTerms() []zjuconst.ExamYearAndTerm {
	return zs.GetCtxConfig().GetExamYearAndTerms()
}

func (zs *NewGrsService) GetCtxConfig() *zjuconst.ZjuScheduleConfig {
	if config := zs.ctx.Value(zjuconst.ScheduleCtxKey); config != nil {
		return config.(*zjuconst.ZjuScheduleConfig)
	} else {
		return zjuconst.GetConfig()
	}
}

func (zs *NewGrsService) UpdateConfig() bool {
	//TODO
	//如此设计我们是否需要update？
	//或者后台开个协程每分钟和fs同步？
	return true
}

func NewNewGrsService(ctx context.Context, isUgrs bool) *NewGrsService {
	return &NewGrsService{
		ZJUClient: zjuam.NewClient(),
		ctx:       log.With().Str("reqid", uuid.NewString()).Logger().WithContext(ctx),
		isUGRS:    isUgrs,
	}
}
