package zjuconst

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ConfigDefaultPath = "configs/config.json"
const OnlineConfigPath = "https://mirror.ghproxy.com/https://raw.githubusercontent.com/cxz66666/zju-ical/master/configs/config.json"

var UseOnlineConfig = false

var _schedule ZjuScheduleConfig
var ScheduleRwMutex sync.RWMutex
var ScheduleCtxKey = "ctx_schedule"

type ClassYearAndTerm struct {
	Year string
	Term ClassTerm
}

type ExamYearAndTerm struct {
	Year string
	Term ExamTerm
}

type YearAndSemester struct {
	Year     string `json:"year"`
	Semester string `json:"semester"`
}

type ZjuScheduleConfig struct {
	LastUpdated     int              `json:"lastUpdated"`
	LastUpdatedTime string           `json:"lastUpdatedTime"`
	Tweaks          []TweakJson      `json:"tweaks"`
	TermConfigs     []TermConfigJson `json:"termConfigs"`
	ClassTerms      []string         `json:"classTerms"`
	ExamTerms       []string         `json:"examTerms"`
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

func (config *ZjuScheduleConfig) GetClassYearAndSemester() []YearAndSemester {
	var res []YearAndSemester
	for _, item := range config.ClassTerms {
		splits := strings.Split(item, ":")
		res = append(res, YearAndSemester{
			Year: splits[0],
			// convert like "1" to "冬学期"
			Semester: ClassTermStrToStr(splits[1]),
		})
	}
	return res
}

func (config *ZjuScheduleConfig) GetExamYearAndSemester() []YearAndSemester {
	var res []YearAndSemester
	for _, item := range config.ExamTerms {
		splits := strings.Split(item, ":")
		res = append(res, YearAndSemester{
			Year: splits[0],
			// convert like "1" to "春夏学期"
			Semester: ExamStrToStr(splits[1]),
		})
	}
	return res
}

func (config *ZjuScheduleConfig) GetLastUpdated() int {
	return config.LastUpdated
}

func (config *ZjuScheduleConfig) GetLastUpdatedTime() string {
	return config.LastUpdatedTime
}

// LoadConfig need be used in program init!!!
func LoadConfig(path string) error {
	var r io.Reader
	if len(path) == 0 {
		path = OnlineConfigPath
	}
	if strings.HasPrefix(path, "http") {
		res, err := http.DefaultClient.Get(path)
		if err != nil {
			return err
		}
		r = res.Body
		defer res.Body.Close()
		UseOnlineConfig = true
	} else {
		f, err := os.Open(path)
		defer f.Close()
		r = f
		if err != nil {
			return err
		}
	}
	log.Info().Msgf("[server] using config uri %s", path)

	cfd, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	ScheduleRwMutex.Lock()
	err = json.Unmarshal(cfd, &_schedule)
	_schedule.LastUpdatedTime = time.Now().Format("2006.01.02 15:04:05")
	ScheduleRwMutex.Unlock()
	return err
}

func UpdateConfig(interval time.Duration) {
	log.Info().Msgf("[server] update online config every %s", interval.String())
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := LoadConfig(OnlineConfigPath); err != nil {
				log.Warn().Msgf("[server] update online error %v", err)
			} else {
				log.Info().Msgf("[server] update online success %s", time.Now().Format("2006.01.02 15:04:05"))
			}
		}
	}
}

func GetConfig() *ZjuScheduleConfig {
	ScheduleRwMutex.RLock()
	defer ScheduleRwMutex.RUnlock()
	return &_schedule
}
