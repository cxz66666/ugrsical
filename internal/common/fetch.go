package common

import (
	"context"
	"errors"
	"time"

	"ugrs-ical/pkg/ical"
	"ugrs-ical/pkg/zjuservice"

	"github.com/rs/zerolog/log"
)

func FetchToMemory(ctx context.Context, username, password string, config Config, tweaks TweakConfig) (string, error) {
	c := zjuapi.NewClient()
	log.Ctx(ctx).Info().Msgf("logging in for %s", username)
	err := c.Login(ctx, zjuapi.GrsLoginUrl, username, password)
	if err != nil {
		return "", err
	}

	var ve []ical.VEvent
	for _, fc := range config.FetchConfig {
		log.Ctx(ctx).Info().Msgf("fetching %d-%d", fc.Year, fc.Semester)
		r, err := c.FetchTimetable(ctx, fc.Year, zjuapi.GrsSemester(fc.Semester))
		if err != nil {
			return "", err
		}

		table, err := timetable.GetTable(r)
		if err != nil {
			return "", err
		}

		log.Ctx(ctx).Info().Msgf("parsing %d-%d", fc.Year, fc.Semester)
		cl, err := timetable.ParseTable(ctx, table)
		if err != nil {
			return "", err
		}

		fm, err := time.ParseInLocation("20060102", fc.FirstDay, time.Local)
		if err != nil {
			return "", err
		}

		log.Ctx(ctx).Info().Msgf("generating vevents %d-%d", fc.Year, fc.Semester)
		vEvents, err := timetable.ClassToVEvents(ctx, fm, cl, &tweaks.Tweaks)
		if err != nil {
			return "", err
		}

		ve = append(ve, *vEvents...)
	}

	log.Ctx(ctx).Info().Msgf("generating iCalendar file")
	iCal := ical.VCalendar{VEvents: ve}
	return iCal.GetICS(""), nil
}
func firstMatchTerm(configs []zjuservice.TermConfig, target zjuservice.ClassYearAndTerm) int {
	for index, item := range configs {
		if item.Term == target.Term && item.Year == target.Year {
			return index
		}
	}
	return -1
}

func GetClassCalendar(ctx context.Context, username, password string) error {
	var zs zjuservice.IZjuService
	zs = zjuservice.NewZjuService(ctx)

	if err := zs.Login(username, password); err != nil {
		return err
	}

	termConfigs := zs.GetTermConfigs()
	tweaks := zs.GetTweaks()

	vCal := ical.VCalendar{}

	for _, item := range zs.GetClassTerms() {
		index := firstMatchTerm(termConfigs, item)
		if index == -1 {
			return errors.New("配置文件错误，未找到指定学期的具体配置")
		}
		classOutline := zs.GetClassTimeTable(item.Year, item.Term, username)

		vCal.VEvents = append(vCal.VEvents)
	}
}
