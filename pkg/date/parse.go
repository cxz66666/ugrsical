package date

import "time"

// TODO check
const layout = "2006-01-02  15:04:05"

// DayToDateTime convert 20210614 format to time.Time
func DayToDateTime(from int) time.Time {
	return time.Date(from/10000, time.Month(from/100%100), from%100, 0, 0, 0, 0, time.Local)
}

func ExamToDateTime(str string) (time.Time, time.Time) {
	return time.Time{}, time.Time{}
}
