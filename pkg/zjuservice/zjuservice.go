package zjuservice

import (
	"ugrs-ical/pkg/zjuservice/zjuconst"
)

type IZJUService interface {
	Login(username, password string) error
	GetClassTimeTable(academicYear string, term zjuconst.ClassTerm, stuId string) ([]zjuconst.ZJUClass, error)
	GetExamInfo(academicYear string, term zjuconst.ExamTerm, stuId string) ([]zjuconst.ZJUExam, error)
	GetScoreInfo(stuId string) ([]zjuconst.ZJUClassScore, error)
	GetTermConfigs() []zjuconst.TermConfig
	GetTweaks() []zjuconst.Tweak
	GetClassTerms() []zjuconst.ClassYearAndTerm
	GetExamTerms() []zjuconst.ExamYearAndTerm
	GetCtxConfig() *zjuconst.ZjuScheduleConfig
	UpdateConfig() bool
}
