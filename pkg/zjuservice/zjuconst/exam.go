package zjuconst

import (
	"fmt"
	"strings"
	"ugrs-ical/pkg/date"
	"ugrs-ical/pkg/ical"
)

type ZJUExam struct {
	ClassName           string
	FinalExamAgenda     string
	FinalExamLocation   string
	FinalExamSeatNum    string
	MidtermExamAgenda   string
	MidtermExamLocation string
	MidtermExamSeatNum  string
}

func (zjuexam *ZJUExam) ToVEventList() []ical.VEvent {
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
