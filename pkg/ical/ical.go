package ical

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

func toISOString(d time.Time) string {
	return fmt.Sprintf("%04d%02d%02dT%02d%02d%02d", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
}

const dateLayoutUTC = "20060102T150405Z"

type VEvent struct {
	Description string
	Summary     string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
}

func (v *VEvent) GetHash() string {
	m := sha1.New()
	m.Write([]byte(v.Description))
	m.Write([]byte(v.Summary))
	m.Write([]byte(v.Location))
	m.Write([]byte(toISOString(v.StartTime)))

	return fmt.Sprintf("%x", m.Sum([]byte("")))
}

func (v *VEvent) String() string {
	var b strings.Builder
	utcStr := time.Now().UTC().Format(dateLayoutUTC)
	stStr := toISOString(v.StartTime)
	etStr := toISOString(v.EndTime)

	b.WriteString(fmt.Sprintf("BEGIN:VEVENT\r\nCLASS:PUBLIC\r\nCREATED:%s\r\n", utcStr))
	if v.Description != "" {
		if strings.Contains(v.Description, "=0D=0A") {
			b.WriteString(fmt.Sprintf("DESCRIPTION;ENCODING=QUOTED-PRINTABLE:%s\r\n", v.Description))
		} else {
			b.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", v.Description))
		}
	}
	b.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", utcStr))
	b.WriteString(fmt.Sprintf("DTSTART;TZID=Asia/Shanghai:%s\r\n", stStr))
	b.WriteString(fmt.Sprintf("DTEND;TZID=Asia/Shanghai:%s\r\n", etStr))
	b.WriteString(fmt.Sprintf("LAST-MODIFIED:%s\r\n", utcStr))
	if v.Location != "" {
		b.WriteString(fmt.Sprintf("LOCATION: %s\r\n", v.Location))
	}
	b.WriteString(fmt.Sprintf("SEQUENCE:0\r\nSUMMARY;LANGUAGE=zh-cn:%s\r\nTRANSP:OPAQUE\r\nUID:%s\r\n", v.Summary, v.GetHash()))
	// Don't display !!!!
	b.WriteString("END:VEVENT\r\n")
	//b.WriteString("BEGIN:VALARM\r\nTRIGGER:-PT15M\r\nACTION:DISPLAY\r\nDESCRIPTION:提醒\r\nEND:VALARM\r\nEND:VEVENT\r\n")

	return b.String()
}

type VCalendar struct {
	VEvents []VEvent
}

func (v *VCalendar) GetICS(icalName string) string {
	var b strings.Builder
	if icalName == "" {
		icalName = "ZJU-ICAL 课程表"
	}
	b.WriteString(fmt.Sprintf("BEGIN:VCALENDAR\r\nX-WR-CALNAME:%s\r\nX-APPLE-CALENDAR-COLOR:#2BBFF0\r\nPRODID:-//ZJUEVA Workshop//Ejector 0.2//EN\r\nVERSION:2.0\r\nMETHOD:PUBLISH\r\nBEGIN:VTIMEZONE\r\nTZID:Asia/Shanghai\r\nBEGIN:STANDARD\r\nDTSTART:16010101T000000\r\nTZOFFSETFROM:+0800\r\nTZOFFSETTO:+0800\r\nEND:STANDARD\r\nEND:VTIMEZONE\r\n", icalName))
	for _, v := range v.VEvents {
		b.WriteString(v.String())
	}
	b.WriteString("END:VCALENDAR\r\n")
	return b.String()
}
