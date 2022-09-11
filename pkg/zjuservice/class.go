package zjuservice

type WeekArrangement int

type ClassTerm int

const (
	Normal   WeekArrangement = iota //每周都有
	OddOnly                         //单周
	EvenOnly                        //双周
)

const (
	Autumn         ClassTerm = iota //秋学期
	Winter                          //冬学期
	ShortA                          //短学期A
	SummerVacation                  //小学期
	Spring                          //春学期
	Summer                          //夏学期
	ShortB                          //短学期B
)

type ZjuClass struct {
	WeekArrangement  WeekArrangement
	StartPeriod      int
	EndPeriod        int
	TeacherName      string
	ClassCode        string
	ClassName        string
	ClassLocation    string
	TermArrangements []ClassTerm
	DayNumber        int
	ClassYear        int
}
