package zjuservice

import (
	"strconv"
	"strings"
)

var _schedule ZjuScheduleConfig

type ClassYearAndTerm struct {
	Year string
	Term ClassTerm
}

type ExamYearAndTerm struct {
	Year string
	Term ExamTerm
}

type ZjuScheduleConfig struct {
	LastUpdated int              `json:"lastUpdated"`
	Tweaks      []TweakJson      `json:"tweaks"`
	TermConfigs []TermConfigJson `json:"termConfigs"`
	ClassTerms  []string         `json:"classTerms"`
	ExamTerms   []string         `json:"examTerms"`
}

func (config *ZjuScheduleConfig) GetClassYearAndTerms() []ClassYearAndTerm {
	var res []ClassYearAndTerm
	for _, item := range config.ClassTerms {
		p := strings.Split(item, ":")
		term, _ := strconv.ParseInt(p[1], 10, 64)
		res = append(res, ClassYearAndTerm{
			Year: p[0],
			Term: ClassTerm(term),
		})
	}
	return res
}

func (config *ZjuScheduleConfig) GetExamYearAndTerms() []ExamYearAndTerm {
	var res []ExamYearAndTerm
	for _, item := range config.ExamTerms {
		p := strings.Split(item, ":")
		term, _ := strconv.ParseInt(p[1], 10, 64)
		res = append(res, ExamYearAndTerm{
			Year: p[0],
			Term: ExamTerm(term),
		})
	}
	return res
}
