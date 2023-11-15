package zjuconst

import (
	"time"
	"zju-ical/pkg/date"
)

type TermConfig struct {
	Id          int
	Year        string
	Term        ClassTerm
	Begin       time.Time
	End         time.Time
	FirstWeekNo int
}

type TermConfigJson struct {
	Year        string `json:"Year"`
	Term        int    `json:"Term"`
	Begin       int    `json:"Begin"`
	End         int    `json:"End"`
	FirstWeekNo int    `json:"FirstWeekNo"`
}

func (tcj TermConfigJson) ToTermConfig() TermConfig {
	return TermConfig{
		Year:        tcj.Year,
		Term:        ClassTerm(tcj.Term),
		Begin:       date.DayToDateTime(tcj.Begin),
		End:         date.DayToDateTime(tcj.End),
		FirstWeekNo: tcj.FirstWeekNo,
	}
}
