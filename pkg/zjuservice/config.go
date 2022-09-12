package zjuservice

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const defaultPath = "configs/config.json"

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

// LoadConfig need be used in program init!!!
func LoadConfig(path string) error {
	var r io.Reader
	if len(path) == 0 {
		path = defaultPath
	}
	if strings.HasPrefix(path, "http") {
		res, err := http.DefaultClient.Get(path)
		if err != nil {
			return err
		}
		r = res.Body
		defer res.Body.Close()
	} else {
		f, err := os.Open(path)
		defer f.Close()
		r = f
		if err != nil {
			return err
		}
	}
	cfd, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(cfd, &_schedule)
	return err
}
