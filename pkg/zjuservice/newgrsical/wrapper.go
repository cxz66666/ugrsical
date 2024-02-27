package newgrsical

import "encoding/json"

type ZjuResWrapperStr[T ZjuWeeklyScheduleRes | ZjuExamOutlineRes | ZjuLoginToken] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    T      `json:"result"`
}

type ZjuWeeklyScheduleRes struct {
	KCBMap map[string]json.RawMessage `json:"kcbMap"`
}

type ZjuWeeklyScheduleItem struct {
	Classes []ZjuWeeklyScheduleClass `json:"pyKcbjSjddVOList"`
}

type ZjuExamOutlineRes struct {
	ExamOutlineList []ZjuExamOutline `json:"records"`
}
type ZjuLoginToken struct {
	Token string `json:"token"`
}
