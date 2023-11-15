package zjuconst

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"ugrs-ical/pkg/ical"

	"github.com/mattn/go-runewidth"
)

const BlankSpace = 25

var scoreMap = map[string]float64{
	"A+":  95,
	"A":   90,
	"A-":  87,
	"B+":  83,
	"B":   80,
	"B-":  77,
	"C+":  73,
	"C":   70,
	"C-":  67,
	"D":   60,
	"F":   0,
	"优":   90,
	"优秀":  90,
	"良":   80,
	"良好":  80,
	"中":   70,
	"中等":  70,
	"及格":  60,
	"不及格": 0,
	"合格":  75,
	"不合格": 0,
}

type ZJUClassScore struct {
	Score     string  `json:"cj"`
	Gpa       float64 `json:"jd"`
	ClassName string  `json:"kcmc"`
	//TODO what is PX?
	PX          int    `json:"px"`
	ClassCredit string `json:"xf"`
	ClassYear   string `json:"xn"`
	ClassTerm   string `json:"xq"`
}

func (zcs *ZJUClassScore) GetScore() float64 {
	if score, ok := scoreMap[zcs.Score]; ok {
		return score
	} else {
		s, err := strconv.ParseFloat(zcs.Score, 64)
		if err != nil {
			return 0
		} else {
			return s
		}
	}
}

func (zcs *ZJUClassScore) GetDescription() string {
	runewidth.EastAsianWidth = true
	className := runewidth.Truncate(zcs.ClassName, BlankSpace, "..")
	return fmt.Sprintf("%s(%s):  %s/%.1f\\n", className, zcs.ClassCredit, zcs.Score, zcs.Gpa)
}

func (zcs *ZJUClassScore) GetScoreTermDescription() string {
	term := ""
	switch zcs.PX {
	case 1:
		term = "秋冬学期"
	case 2:
		term = "春夏学期"
	default:
		term = "未知学期"
	}
	return fmt.Sprintf("%s %s", zcs.ClassYear, term)
}

func GetAverageGPA(scores []ZJUClassScore) (float64, error) {
	var sumGpa float64
	var sumCredits float64
	for _, score := range scores {
		credit, err := strconv.ParseFloat(score.ClassCredit, 64)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("无法解析%s课程的学分，学分为%s", score.ClassName, score.ClassCredit))
		}
		sumGpa += score.Gpa * credit
		sumCredits += credit
	}
	return sumGpa / sumCredits, nil
}

func GetAverageScore(scores []ZJUClassScore) (float64, error) {
	var sumScore float64
	var sumCredits float64
	for _, score := range scores {
		credit, err := strconv.ParseFloat(score.ClassCredit, 64)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("无法解析%s课程的学分，学分为%s", score.ClassName, score.ClassCredit))
		}
		sumScore += credit * score.GetScore()
		sumCredits += credit
	}
	return sumScore / sumCredits, nil
}

func GetTotalCredit(scores []ZJUClassScore) float64 {
	var sumCredits float64
	for _, score := range scores {
		credit, _ := strconv.ParseFloat(score.ClassCredit, 64)
		sumCredits += credit
	}
	return sumCredits
}

func ScoresCleanUp(scores []ZJUClassScore) []ZJUClassScore {
	bestScores := make(map[string]float64)
	for _, score := range scores {
		if score.Score == "弃修" || score.Score == "缓考" || score.Score == "缺考" {
			continue
		}
		if val, exist := bestScores[score.ClassName]; !exist {
			bestScores[score.ClassName] = score.GetScore()
		} else {
			bestScores[score.ClassName] = math.Max(val, score.GetScore())
		}
	}
	var newScores []ZJUClassScore
	for _, score := range scores {
		if val, exist := bestScores[score.ClassName]; exist && val == score.GetScore() {
			newScores = append(newScores, score)
		}
	}
	return newScores
}

func ScoresToVEventList(scores []ZJUClassScore) ([]ical.VEvent, error) {
	scoresMap := make(map[int][]ZJUClassScore)

	termsCount := 0
	termsMap := make(map[string]int)

	nowTermDesc := ""

	for _, score := range scores {
		termDesc := score.GetScoreTermDescription()
		if termDesc != nowTermDesc {
			nowTermDesc = termDesc
			termsCount++
			termsMap[termDesc] = termsCount
		}
		scoresMap[termsCount] = append(scoresMap[termsCount], score)
	}

	var vEvents []ical.VEvent

	event, err := ScoresToVEvent(scores, "GPA 总览", 0)
	if err != nil {
		return nil, err
	}
	vEvents = append(vEvents, event)

	for i := 1; i <= termsCount; i++ {
		classScores := scoresMap[i]
		if len(classScores) == 0 {
			continue
		}
		summary := classScores[0].GetScoreTermDescription()
		event, err = ScoresToVEvent(classScores, summary, i)
		if err != nil {
			continue
		}
		vEvents = append(vEvents, event)
	}

	return vEvents, nil
}

func ScoresToVEvent(scores []ZJUClassScore, summary string, index int) (ical.VEvent, error) {
	if summary == "" {
		summary = "GPA Helper"
	}
	averageGpa, err := GetAverageGPA(scores)
	if err != nil {
		return ical.VEvent{}, err
	}
	averageScore, err := GetAverageScore(scores)
	if err != nil {
		return ical.VEvent{}, err
	}
	totalCredit := GetTotalCredit(scores)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("平均绩点: %.2f\\n平均分: %.2f\\n总学分: %.2f\\n\\n", averageGpa, averageScore, totalCredit))

	if index != 0 {
		for _, score := range scores {
			b.WriteString(score.GetDescription())
		}
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	beginTime := time.Minute * 15 * time.Duration(index)
	return ical.VEvent{
		Summary:     summary,
		Description: b.String(),
		StartTime:   time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 21, 30, 0, 0, time.Local).Add(beginTime),
		EndTime:     time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 21, 30, 0, 0, time.Local).Add(beginTime).Add(time.Minute * 15)}, nil
}
