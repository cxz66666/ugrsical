package zjuservice

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type WeekArrangement int

type ClassTerm = int

const (
	Normal   WeekArrangement = iota //每周都有
	OddOnly                         //单周
	EvenOnly                        //双周
)

const (
	Autumn         ClassTerm = iota //秋学期
	Winter                          //冬学期
	ShortA                          //短学期A
	SummerVacation                  //小学期
	Spring                          //春学期
	Summer                          //夏学期
	ShortB                          //短学期B
)

type ZjuClass struct {
	WeekArrangement  WeekArrangement
	StartPeriod      int
	EndPeriod        int
	TeacherName      string
	ClassCode        string
	ClassName        string
	ClassLocation    string
	TermArrangements []ClassTerm
	DayNumber        int
	ClassYear        int
}

func NewZjuClass() ZjuClass {
	return ZjuClass{
		TermArrangements: make([]ClassTerm, 0),
	}
}

func (zc *ZjuClass) GetStartDateTime(day time.Time) time.Time {
	period := NewClassPeriod(zc.StartPeriod)
	return period.ToStartDateTime(day)
}

func (zc *ZjuClass) GetEndDateTime(day time.Time) time.Time {
	period := NewClassPeriod(zc.EndPeriod)
	return period.ToEndDateTime(day)
}

func (zc *ZjuClass) ArrangementDescription() string {
	var b strings.Builder
	sort.Ints(zc.TermArrangements)
	for _, term := range zc.TermArrangements {
		b.WriteString(ClassTermToDescriptionString(term))
	}
	b.WriteString(" ")
	if zc.StartPeriod == zc.EndPeriod {
		b.WriteString(fmt.Sprintf("第%d节", zc.StartPeriod))
	} else {
		b.WriteString(fmt.Sprintf("第%d-%d节", zc.StartPeriod, zc.EndPeriod))
	}
	return b.String()
}

type ClassPeriod struct {
	hour   int
	minute int
}

func NewClassPeriod(periodNumber int) ClassPeriod {
	Hour := 0
	Minute := 0
	switch periodNumber {
	case 1:
		Hour = 8
		Minute = 0
	case 2:
		Hour = 8
		Minute = 50
	case 3:
		Hour = 9
		Minute = 50
	case 4:
		Hour = 10
		Minute = 40
	case 5:
		Hour = 11
		Minute = 30
	case 6:
		Hour = 13
		Minute = 15
	case 7:
		Hour = 14
		Minute = 5
	case 8:
		Hour = 14
		Minute = 55
	case 9:
		Hour = 15
		Minute = 55
	case 10:
		Hour = 16
		Minute = 45
	case 11:
		Hour = 18
		Minute = 30
	case 12:
		Hour = 19
		Minute = 20
	case 13:
		Hour = 20
		Minute = 10
	default:
		Hour = 21
		Minute = 0
	}
	return ClassPeriod{
		hour:   Hour,
		minute: Minute,
	}
}

func (cp *ClassPeriod) ToStartDateTime(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), cp.hour, cp.minute, 0, 0, time.Local)
}

func (cp *ClassPeriod) ToEndDateTime(day time.Time) time.Time {
	return cp.ToStartDateTime(day).Add(time.Minute * 45)
}
