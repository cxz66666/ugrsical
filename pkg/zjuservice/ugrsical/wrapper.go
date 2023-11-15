package ugrsical

import (
	"ugrs-ical/pkg/zjuservice/zjuconst"
)

type ZjuResWrapperStr[T ZjuWeeklyScheduleRes | ZjuExamOutlineRes | ZjuClassScoreRes] struct {
	Data      T      `json:"data"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type ZjuWeeklyScheduleRes struct {
	ClassList []ZjuWeeklyScheduleClass `json:"kblist"`
}

type ZjuExamOutlineRes struct {
	ExamOutlineList []ZjuExamOutline `json:"list"`
	Zt              string           `json:"zt"`
}

type ZjuClassScoreRes struct {
	ClassScoreList []zjuconst.ZJUClassScore `json:"list"`
	Zt             string                   `json:"zt"`
}
