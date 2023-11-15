package zjuconst

import (
	"time"

	"zju-ical/pkg/date"
)

type TweakType int

const (
	Clear TweakType = iota
	Copy
	Exchange
)

type Tweak struct {
	Id          int
	TweakType   TweakType
	Description string
	From        time.Time
	To          time.Time
}

type TweakJson struct {
	TweakType   int    `json:"TweakType"`
	Description string `json:"Description"`
	From        int    `json:"From"`
	To          int    `json:"To"`
}

func (tj TweakJson) ToTweak() Tweak {
	return Tweak{
		TweakType:   TweakType(tj.TweakType),
		Description: tj.Description,
		From:        date.DayToDateTime(tj.From),
		To:          date.DayToDateTime(tj.To),
	}
}
