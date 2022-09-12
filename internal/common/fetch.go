package common

import (
	"context"
	"errors"

	"ugrs-ical/pkg/ical"
	"ugrs-ical/pkg/zjuservice"

	"github.com/rs/zerolog/log"
)

func firstMatchTerm(configs []zjuservice.TermConfig, target zjuservice.ClassYearAndTerm) int {
	for index, item := range configs {
		if item.Term == target.Term && item.Year == target.Year {
			return index
		}
	}
	return -1
}

func GetClassCalendar(ctx context.Context, username, password string) (ical.VCalendar, error) {
	var zs zjuservice.IZjuService
	zs = zjuservice.NewZjuService(ctx)

	if err := zs.Login(username, password); err != nil {
		return ical.VCalendar{}, err
	}

	termConfigs := zs.GetTermConfigs()
	tweaks := zs.GetTweaks()

	vCal := ical.VCalendar{}

	for _, item := range zs.GetClassTerms() {
		index := firstMatchTerm(termConfigs, item)
		if index == -1 {
			return ical.VCalendar{}, errors.New("配置文件错误，未找到指定学期的具体配置")
		}
		classOutline := zs.GetClassTimeTable(item.Year, item.Term, username)
		log.Ctx(ctx).Info().Msgf("generating class vevents %d-%d", item.Year, item.Term)
		// classes to events
		list := zjuservice.ClassToVEvents(classOutline, termConfigs[index], tweaks)
		vCal.VEvents = append(vCal.VEvents, list...)
		log.Ctx(ctx).Info().Msgf("generated class vevents %d-%d", item.Year, item.Term)
	}
	log.Ctx(ctx).Info().Msgf("get class vCal success ")

	// TODO cache
	return vCal, nil
}

func GetExamCalendar(ctx context.Context, username, password string) (ical.VCalendar, error) {
	var zs zjuservice.IZjuService
	zs = zjuservice.NewZjuService(ctx)

	if err := zs.Login(username, password); err != nil {
		return ical.VCalendar{}, err
	}

	vCal := ical.VCalendar{}

	for _, item := range zs.GetExamTerms() {
		examOutline := zs.GetExamInfo(item.Year, item.Term, username)
		log.Ctx(ctx).Info().Msgf("generating exam vevents %d-%d", item.Year, item.Term)
		// exam to events
		for _, exam := range examOutline {
			vCal.VEvents = append(vCal.VEvents, exam.ToVEventList()...)
		}
		log.Ctx(ctx).Info().Msgf("generated exam vevents %d-%d", item.Year, item.Term)
	}
	log.Ctx(ctx).Info().Msgf("get exam vCal success")

	// TODO cache
	return vCal, nil
}

func GetBothCalendar(ctx context.Context, username, password string) (ical.VCalendar, error) {
	var zs zjuservice.IZjuService
	zs = zjuservice.NewZjuService(ctx)

	if err := zs.Login(username, password); err != nil {
		return ical.VCalendar{}, err
	}

	termConfigs := zs.GetTermConfigs()
	tweaks := zs.GetTweaks()

	vCal := ical.VCalendar{}

	for _, item := range zs.GetClassTerms() {
		index := firstMatchTerm(termConfigs, item)
		if index == -1 {
			return ical.VCalendar{}, errors.New("配置文件错误，未找到指定学期的具体配置")
		}
		classOutline := zs.GetClassTimeTable(item.Year, item.Term, username)
		log.Ctx(ctx).Info().Msgf("generating class vevents %d-%d", item.Year, item.Term)
		// classes to events
		list := zjuservice.ClassToVEvents(classOutline, termConfigs[index], tweaks)
		vCal.VEvents = append(vCal.VEvents, list...)
		log.Ctx(ctx).Info().Msgf("generated class vevents %d-%d", item.Year, item.Term)
	}
	log.Ctx(ctx).Info().Msgf("get class vCal success ")

	for _, item := range zs.GetExamTerms() {
		examOutline := zs.GetExamInfo(item.Year, item.Term, username)
		log.Ctx(ctx).Info().Msgf("generating exam vevents %d-%d", item.Year, item.Term)
		// exam to events
		for _, exam := range examOutline {
			vCal.VEvents = append(vCal.VEvents, exam.ToVEventList()...)
		}
		log.Ctx(ctx).Info().Msgf("generated exam vevents %d-%d", item.Year, item.Term)
	}
	log.Ctx(ctx).Info().Msgf("get exam vCal success")

	// TODO cache
	return vCal, nil

}
