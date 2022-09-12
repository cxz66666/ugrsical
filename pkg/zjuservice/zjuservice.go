package zjuservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"ugrs-ical/pkg/zjuam"
)

const (
	kAppServiceLoginUrl         = "https://zjuam.zju.edu.cn/cas/login?service=http%3A%2F%2Fappservice.zju.edu.cn%2Findex"
	kAppServiceGetWeekClassInfo = "http://appservice.zju.edu.cn/zju-smartcampus/zdydjw/api/kbdy_cxXsZKbxx"
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
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if err != nil {
		//TODO log
		return nil
	}
	resp, err := zs.ZjuClient.Client().Do(req)
	if err != nil {
		//TODO log
		return nil
	}
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(content)
	return nil
}

func (zs *ZjuService) GetExamInfo(academicYear string, term ExamTerm, stuId string) []ZjuExamOutline {
	//TODO implement me
	panic("implement me")
}

func (zs *ZjuService) GetTermConfigs() []TermConfig {
	//TODO implement me
	panic("implement me")
}

func (zs *ZjuService) GetTweaks() []Tweak {
	//TODO implement me
	panic("implement me")
}

func (zs *ZjuService) GetClassTerms() []ClassYearAndTerm {
	//TODO implement me
	panic("implement me")
}

func (zs *ZjuService) GetExamTerms() []ExamYearAndTerm {
	//TODO implement me
	panic("implement me")
}

func (zs *ZjuService) UpdateConfig() bool {
	//TODO implement me
	panic("implement me")
}

func NewZjuService() ZjuService {
	return ZjuService{
		ZjuClient: zjuam.NewClient(),
	}
}
