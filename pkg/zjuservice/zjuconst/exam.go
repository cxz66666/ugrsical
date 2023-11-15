package zjuconst

import (
	"fmt"
	"strings"
	"time"
	"zju-ical/pkg/date"
	"zju-ical/pkg/ical"
)

type ExamTerm int

const (
	AutumnWinter ExamTerm = iota
	SpringSummer
)

type ZJUExam interface {
	ToVEventList() []ical.VEvent
}

type ZJUUgrsExam struct {
	ClassName           string
	FinalExamAgenda     string
	FinalExamLocation   string
	FinalExamSeatNum    string
	MidtermExamAgenda   string
	MidtermExamLocation string
	MidtermExamSeatNum  string
}

type ZJUGrsExam struct {
	Semester     string
	ID           string // 课号
	ClassName    string
	Region       string //区域，类似：玉泉
	StartTime    time.Time
	EndTime      time.Time
	ExamLocation string
	ExamSeatNum  string
	ExamRemark   string
}

func (zjuexam *ZJUUgrsExam) ToVEventList() []ical.VEvent {
	var vEvents []ical.VEvent
	if len(strings.TrimSpace(zjuexam.FinalExamAgenda)) != 0 {
		//TODO WIP
		finalEventStart, finalEventEnd := date.ExamToDateTime(zjuexam.FinalExamAgenda)
		location := ""
		if len(strings.TrimSpace(zjuexam.FinalExamLocation)) != 0 {
			location = zjuexam.FinalExamLocation
		}
		description := ""
		if len(strings.TrimSpace(zjuexam.FinalExamSeatNum)) != 0 {
			description = fmt.Sprintf("座位号：%s", zjuexam.FinalExamSeatNum)
		}
		vEvents = append(vEvents, ical.VEvent{
			Summary:     fmt.Sprintf("[务必核对!]%s 期末考试", zjuexam.ClassName),
			Location:    location,
			Description: description,
			StartTime:   finalEventStart,
			EndTime:     finalEventEnd,
		},
		)
	}

	if len(strings.TrimSpace(zjuexam.MidtermExamAgenda)) != 0 {
		//TODO WIP
		midTermEventStart, midTermEventEnd := date.ExamToDateTime(zjuexam.MidtermExamAgenda)
		location := ""
		if len(strings.TrimSpace(zjuexam.MidtermExamLocation)) != 0 {
			location = zjuexam.MidtermExamLocation
		}
		description := ""
		if len(strings.TrimSpace(zjuexam.MidtermExamLocation)) != 0 {
			description = fmt.Sprintf("座位号：%s", zjuexam.MidtermExamLocation)
		}
		vEvents = append(vEvents, ical.VEvent{
			Summary:     fmt.Sprintf("%s 期中考试", zjuexam.ClassName),
			Location:    location,
			Description: description,
			StartTime:   midTermEventStart,
			EndTime:     midTermEventEnd,
		},
		)
	}
	return vEvents

}

func (zjuexam *ZJUGrsExam) ToVEventList() []ical.VEvent {
	var vEvents []ical.VEvent

	//TODO WIP

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s学期 ", zjuexam.ClassName, zjuexam.Semester))
	b.WriteString(fmt.Sprintf("课程号：%s\\n", zjuexam.ID))
	if zjuexam.ExamSeatNum != "" {
		b.WriteString(fmt.Sprintf("座位号: %s\\n", zjuexam.ExamSeatNum))
	}
	if zjuexam.ExamRemark != "" {
		b.WriteString(fmt.Sprintf("备注: %s\\n", zjuexam.ExamRemark))
	}
	vEvents = append(vEvents, ical.VEvent{
		Summary:     fmt.Sprintf("[务必核对!]%s %s学期考试", zjuexam.ClassName, zjuexam.Semester),
		Location:    zjuexam.ExamLocation,
		Description: b.String(),
		StartTime:   zjuexam.StartTime,
		EndTime:     zjuexam.EndTime,
	},
	)
	return vEvents

}
