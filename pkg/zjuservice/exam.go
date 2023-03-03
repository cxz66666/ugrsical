package zjuservice

import (
	"fmt"
	"strings"

	"ugrs-ical/pkg/date"
	"ugrs-ical/pkg/ical"
)

type ExamTerm int

const (
	AutumnWinter ExamTerm = iota
	SpringSummer
)

type ZjuExamOutline struct {
	ClassIdWithStuId    string `json:"kcid"`
	ClassName           string `json:"kcmc"`
	FinalExamAgenda     string `json:"qmkssj"`
	FinalExamLocation   string `json:"qmksdd"`
	FinalExamSeatNum    string `json:"zwxh"`
	MidtermExamAgenda   string `json:"qzkssj"`
	MidtermExamLocation string `json:"qzksdd"`
	MidtermExamSeatNum  string `json:"qzzwxh"`
	ClassId             string `json:"xkkh"`
	ClassCredit         string `json:"xkxf"`
	ClassTerm           string `json:"xq"`
}

func (zjuexo *ZjuExamOutline) ToVEventList() []ical.VEvent {
	var vEvents []ical.VEvent
	if len(strings.TrimSpace(zjuexo.FinalExamAgenda)) != 0 {
		//TODO WIP
		finalEventStart, finalEventEnd := date.ExamToDateTime(zjuexo.FinalExamAgenda)
		location := ""
		if len(strings.TrimSpace(zjuexo.FinalExamLocation)) != 0 {
			location = zjuexo.FinalExamLocation
		}
		description := ""
		if len(strings.TrimSpace(zjuexo.FinalExamSeatNum)) != 0 {
			description = fmt.Sprintf("座位号：%s", zjuexo.FinalExamSeatNum)
		}
		vEvents = append(vEvents, ical.VEvent{
			Summary:     fmt.Sprintf("[务必核对!]%s 期末考试", zjuexo.ClassName),
			Location:    location,
			Description: description,
			StartTime:   finalEventStart,
			EndTime:     finalEventEnd,
		},
		)
	}

	if len(strings.TrimSpace(zjuexo.MidtermExamAgenda)) != 0 {
		//TODO WIP
		midTermEventStart, midTermEventEnd := date.ExamToDateTime(zjuexo.MidtermExamAgenda)
		location := ""
		if len(strings.TrimSpace(zjuexo.MidtermExamLocation)) != 0 {
			location = zjuexo.MidtermExamLocation
		}
		description := ""
		if len(strings.TrimSpace(zjuexo.MidtermExamLocation)) != 0 {
			description = fmt.Sprintf("座位号：%s", zjuexo.MidtermExamLocation)
		}
		vEvents = append(vEvents, ical.VEvent{
			Summary:     fmt.Sprintf("%s 期中考试", zjuexo.ClassName),
			Location:    location,
			Description: description,
			StartTime:   midTermEventStart,
			EndTime:     midTermEventEnd,
		},
		)
	}
	return vEvents
}
