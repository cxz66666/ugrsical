package ugrsical

import (
	"sort"
	"strconv"
	"strings"
	"zju-ical/pkg/zjuservice/zjuconst"
)

type ZjuWeeklyScheduleClass struct {
	WeekArrangement string   `json:"dsz"`
	Periods         []string `json:"jc"`
	TeacherName     string   `json:"jsxm"`
	ClassCode       string   `json:"kcdm"`
	ClassId         string   `json:"kcid"`
	ClassName       string   `json:"mc"`
	ClassLocation   string   `json:"skdd"`
	TermArrangement string   `json:"xq"`
	DayNumber       int      `json:"xqj"`
	IsConfirmed     int      `json:"sfqd"`
}

func (zwsc ZjuWeeklyScheduleClass) ToZJUClass() *zjuconst.ZJUClass {
	if zwsc.IsConfirmed == 0 {
		return nil
	}
	var res zjuconst.ZJUClass
	if strings.Contains(zwsc.TermArrangement, "秋") {
		res.TermArrangements = append(res.TermArrangements, zjuconst.Autumn)
	}
	if strings.Contains(zwsc.TermArrangement, "冬") {
		res.TermArrangements = append(res.TermArrangements, zjuconst.Winter)
	}
	if strings.Contains(zwsc.TermArrangement, "春") {
		res.TermArrangements = append(res.TermArrangements, zjuconst.Spring)
	}
	if strings.Contains(zwsc.TermArrangement, "夏") {
		res.TermArrangements = append(res.TermArrangements, zjuconst.Summer)
	}
	if len(res.TermArrangements) == 0 {
		return nil
	}
	year, err := strconv.ParseInt(zwsc.ClassId[1:5], 10, 64)
	if err != nil {
		return nil
	}
	res.ClassYear = int(year)

	switch zwsc.WeekArrangement {
	case "0":
		res.WeekArrangement = zjuconst.OddOnly
	case "1":
		res.WeekArrangement = zjuconst.EvenOnly
	default:
		res.WeekArrangement = zjuconst.Normal
	}

	periods := make([]int, 0)
	for _, v := range zwsc.Periods {
		period, _ := strconv.ParseInt(v, 10, 64)
		periods = append(periods, int(period))
	}
	sort.Ints(periods)
	res.StartPeriod = periods[0]
	res.EndPeriod = periods[len(periods)-1]
	res.TeacherName = zwsc.TeacherName
	res.ClassCode = zwsc.ClassCode
	res.ClassName = zwsc.ClassName
	res.ClassLocation = zwsc.ClassLocation
	res.DayNumber = zwsc.DayNumber

	return &res
}
