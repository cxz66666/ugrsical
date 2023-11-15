package ugrsical

import (
	"ugrs-ical/pkg/zjuservice/zjuconst"
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

func (zjuexo *ZjuExamOutline) ToZJUExam() zjuconst.ZJUExam {
	return &zjuconst.ZJUUgrsExam{
		ClassName:           zjuexo.ClassName,
		FinalExamAgenda:     zjuexo.FinalExamAgenda,
		FinalExamLocation:   zjuexo.FinalExamLocation,
		FinalExamSeatNum:    zjuexo.FinalExamSeatNum,
		MidtermExamAgenda:   zjuexo.MidtermExamAgenda,
		MidtermExamLocation: zjuexo.MidtermExamLocation,
	}
}
