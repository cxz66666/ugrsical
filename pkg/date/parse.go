package date

import (
	"strconv"
	"time"
)

// DayToDateTime convert 20210614 format to time.Time
func DayToDateTime(from int) time.Time {
	return time.Date(from/10000, time.Month(from/100%100), from%100, 0, 0, 0, 0, time.Local)
}

// ExamToDateTime convert  2021年07月02日(10:30-12:30) format to time.Time
func ExamToDateTime(str string) (time.Time, time.Time) {
	strRune := []rune(str)
	if len(strRune) < 23 {
		return time.Time{}, time.Time{}
	}
	year, _ := strconv.ParseInt(string(strRune[0:4]), 10, 64)
	month, _ := strconv.ParseInt(string(strRune[5:7]), 10, 64)
	day, _ := strconv.ParseInt(string(strRune[8:10]), 10, 64)
	hour, _ := strconv.ParseInt(string(strRune[12:14]), 10, 64)
	minute, _ := strconv.ParseInt(string(strRune[15:17]), 10, 64)
	endHour, _ := strconv.ParseInt(string(strRune[18:20]), 10, 64)
	endMinute, _ := strconv.ParseInt(string(strRune[21:23]), 10, 64)
	return time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), 0, 0, time.Local), time.Date(int(year), time.Month(month), int(day), int(endHour), int(endMinute), 0, 0, time.Local)
}
