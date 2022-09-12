package zjuservice

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"ugrs-ical/pkg/ical"
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

func GetClassOfDay(classes []ZjuClass, day int) []ZjuClass {
	var res []ZjuClass
	for _, item := range classes {
		if item.DayNumber == day {
			res = append(res, item)
		}
	}
	return res
}

func isEvenWeek(mondayOfTermBegin, target time.Time) bool {
	return ((target.Day()-mondayOfTermBegin.Day())/7)%2 == 1
}

func ClassToVEvents(classes []ZjuClass, termConfig TermConfig, tweaks []Tweak) []ical.VEvent {
	dataLength := termConfig.End.Day() - termConfig.Begin.Day() + 2
	shadowDates := make(map[time.Time]time.Time, dataLength)
	for currentDate := termConfig.Begin; currentDate.Before(termConfig.End) || currentDate.Equal(termConfig.End); currentDate = currentDate.Add(time.Hour * 24) {
		shadowDates[currentDate] = currentDate
	}
	modDescriptions := make(map[time.Time]string, dataLength)
	for _, tweak := range tweaks {
		if tweak.To.Before(termConfig.Begin) || tweak.To.After(termConfig.End) {
			continue
		}
		switch tweak.TweakType {
		case Clear:
			for d := tweak.From; d.Before(tweak.To) || d.Equal(tweak.To); d = d.Add(time.Hour * 24) {
				delete(shadowDates, d)
			}
		case Copy:
			shadowDates[tweak.To] = tweak.From
			modDescriptions[tweak.To] = tweak.Description
		case Exchange:
			shadowDates[tweak.To] = tweak.From
			shadowDates[tweak.From] = tweak.To
			modDescriptions[tweak.To] = tweak.Description
			modDescriptions[tweak.From] = tweak.Description
		}
	}
	classOfDay := make(map[time.Weekday][]ZjuClass, 7)
	classOfDay[time.Monday] = GetClassOfDay(classes, 1)
	classOfDay[time.Tuesday] = GetClassOfDay(classes, 2)
	classOfDay[time.Wednesday] = GetClassOfDay(classes, 3)
	classOfDay[time.Thursday] = GetClassOfDay(classes, 4)
	classOfDay[time.Friday] = GetClassOfDay(classes, 5)
	classOfDay[time.Saturday] = GetClassOfDay(classes, 6)
	classOfDay[time.Sunday] = GetClassOfDay(classes, 7)

	termBeginDayOfWeek := int(termConfig.Begin.Weekday())
	if termBeginDayOfWeek == 0 {
		termBeginDayOfWeek = 7
	}
	mondayOfFirstWeek := termConfig.Begin.Add(time.Hour * time.Duration(-24*(termBeginDayOfWeek-1))).
		Add(time.Hour * time.Duration(-24*7*(termConfig.FirstWeekNo-1)))

	events := make([]ical.VEvent, 3*dataLength)
	for actualDate, dateOfClass := range shadowDates {
		classesOfCurrentDate := classOfDay[dateOfClass.Weekday()]
		isCurrentDateEvenWeek := isEvenWeek(mondayOfFirstWeek, dateOfClass)
		for _, item := range classesOfCurrentDate {
			if (isCurrentDateEvenWeek && item.WeekArrangement == OddOnly) ||
				(!isCurrentDateEvenWeek && item.WeekArrangement == EvenOnly) {
				continue
			}
			description := ""

			if mod, exist := modDescriptions[actualDate]; !exist {
				description = fmt.Sprintf("教师: %s\\n课程代码: %s\\n教学时间安排: %s", item.TeacherName, item.ClassCode, item.ArrangementDescription())
			} else {
				description = fmt.Sprintf("%s\\n教师: %s\\n课程代码: %s\\n教学时间安排: %s", mod, item.TeacherName, item.ClassCode, item.ArrangementDescription())
			}
			events = append(events, ical.VEvent{
				Summary:     item.ClassName,
				StartTime:   item.GetStartDateTime(actualDate),
				EndTime:     item.GetEndDateTime(actualDate),
				Location:    item.ClassLocation,
				Description: description,
			})
		}
	}
	return events
}
