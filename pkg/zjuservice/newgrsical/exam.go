package newgrsical

import (
	"github.com/cxz66666/zju-ical/pkg/zjuservice/zjuconst"
	"strconv"
	"strings"
	"time"
)

type ZjuExamOutline struct {
	// "秋冬学期" or "春夏学期"
	Semester string `json:"xq_dictText"`
	//课号
	ClassId   string `json:"kcbh"`
	ClassName string `json:"kcmc"`
	Region    string `json:"xqmc"`
	// 1030 or 1230
	StartTime int `json:"kssj"`
	// 1030 or 1230
	EndTime      int    `json:"jssj"`
	ExamDay      int    `json:"rq"`
	ExamLocation string `json:"mc"`
	ExamSeatNum  int    `json:"zwh"`
}

func (zjuexo *ZjuExamOutline) ToZJUExam() zjuconst.ZJUExam {
	var res zjuconst.ZJUGrsExam

	semester := zjuexo.Semester
	if strings.HasSuffix(semester, "学期") {
		semester = strings.TrimSuffix(semester, "学期")
	}
	res.Semester = semester
	res.ID = zjuexo.ClassId
	res.ClassName = zjuexo.ClassName
	res.Region = zjuexo.Region
	res.ExamLocation = zjuexo.ExamLocation
	res.ExamSeatNum = strconv.Itoa(zjuexo.ExamSeatNum)
	// TODO?
	res.ExamRemark = ""

	dayTime, err := time.Parse("20060102", strconv.Itoa(zjuexo.ExamDay))
	if err != nil {
		return nil
	}
	startTime := dayTime.Add(time.Hour*time.Duration(zjuexo.StartTime/100) + time.Minute*time.Duration(zjuexo.StartTime%100))
	endTime := dayTime.Add(time.Hour*time.Duration(zjuexo.EndTime/100) + time.Minute*time.Duration(zjuexo.EndTime%100))

	res.StartTime = startTime
	res.EndTime = endTime
	return &res
}
