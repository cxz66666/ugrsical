package newgrsical

import (
	"github.com/cxz66666/zju-ical/pkg/zjuservice/zjuconst"
	"strconv"
	"strings"
)

type ZjuWeeklyScheduleClass struct {
	// TODO what means 13?
	// we can use WeekArrangementExtra to check dsz
	WeekArrangement string `json:"dsz"`
	BeginPeriod     int    `json:"ksjc"`
	EndPeriod       int    `json:"jsjc"`
	// TODO xm always be null
	TeacherName    string `json:"xm"`
	ClassId        string `json:"bjbh"`
	ClassYear      string `json:"xn"`
	ClassName      string `json:"kcmc"`
	ClassLocation  string `json:"cdmc"`
	TermArrangment string `json:"pkxq"`
	DayNumber      int    `json:"xqj"`

	// format like "9,10,11,12,13,14,15,16" or "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16"
	// because we don't know what means of dsz, we can use this field as extra check
	WeekArrangementExtra string `json:"zc"`
}

func (zwsc ZjuWeeklyScheduleClass) ToZJUClass() *zjuconst.ZJUClass {
	var res zjuconst.ZJUClass

	res.TermArrangements = zjuconst.GrsClassQueryStringToClassTerm(zwsc.TermArrangment)
	if len(res.TermArrangements) == 0 {
		return nil
	}
	year, err := strconv.Atoi(zwsc.ClassYear)
	if err != nil {
		year = -1
	}
	res.ClassYear = year

	if strings.Contains(zwsc.WeekArrangementExtra, "1,2") || strings.Contains(zwsc.WeekArrangementExtra, "9,10") {
		res.WeekArrangement = zjuconst.Normal
	} else if strings.Contains(zwsc.WeekArrangementExtra, "1") || strings.Contains(zwsc.WeekArrangementExtra, "9") {
		res.WeekArrangement = zjuconst.OddOnly
	} else {
		res.WeekArrangement = zjuconst.EvenOnly
	}

	res.StartPeriod = zwsc.BeginPeriod
	res.EndPeriod = zwsc.EndPeriod
	res.TeacherName = zwsc.TeacherName
	if len(res.TeacherName) == 0 {
		res.TeacherName = "未知"
	}
	res.ClassCode = zwsc.ClassId
	res.ClassName = zwsc.ClassName
	res.ClassLocation = zwsc.ClassLocation
	res.DayNumber = zwsc.DayNumber
	return &res
}
