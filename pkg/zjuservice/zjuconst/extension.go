package zjuconst

func UgrsClassTermToQueryString(t ClassTerm) string {
	switch t {
	case Autumn:
		return "1|秋"
	case Winter:
		return "1|冬"
	case ShortA:
		return "1|短"
	case SummerVacation:
		return "1|暑"
	case Spring:
		return "2|春"
	case Summer:
		return "2|夏"
	case ShortB:
		return "2|短"
	default:
		return ""
	}
}

func GrsClassTermToClassQueryInt(t ClassTerm) int {
	switch t {
	case Autumn:
		return 13
	case Winter:
		return 14
	case ShortA:
		return 0
	case SummerVacation:
		return 0
	case Spring:
		return 11
	case Summer:
		return 12
	case ShortB:
		return 0
	default:
		return 0
	}
}

// GrsClassQueryStringToClassTerm used for newgrs system
func GrsClassQueryStringToClassTerm(str string) []ClassTerm {
	switch str {
	case "11":
		return []ClassTerm{Spring}
	case "12":
		return []ClassTerm{Summer}
	case "13":
		return []ClassTerm{Autumn}
	case "14":
		return []ClassTerm{Winter}
	case "15":
		return []ClassTerm{Spring, Summer}
	case "16":
		return []ClassTerm{Autumn, Winter}
	default:
		return []ClassTerm{}
	}
}

func ClassTermToDescriptionString(t ClassTerm) string {
	switch t {
	case Autumn:
		return "秋"
	case Winter:
		return "冬"
	case ShortA:
		return "短"
	case SummerVacation:
		return "暑"
	case Spring:
		return "春"
	case Summer:
		return "夏"
	case ShortB:
		return "短"
	default:
		return ""
	}
}

func ExamTermToDescriptionString(t ExamTerm) string {
	switch t {
	case AutumnWinter:
		return "秋冬"
	case SpringSummer:
		return "春夏"
	default:
		return ""
	}
}

func UgrsExamTermToQueryString(t ExamTerm) string {
	switch t {
	case AutumnWinter:
		return "1"
	case SpringSummer:
		return "2"
	default:
		return ""
	}
}

func GrsExamTermToQueryInt(t ExamTerm) int {
	switch t {
	case AutumnWinter:
		return 16
	case SpringSummer:
		return 15
	default:
		return -1
	}
}

func NewGrsExamTermToQueryInt(t ExamTerm) int {
	switch t {
	case AutumnWinter:
		return 12
	case SpringSummer:
		return 11
	default:
		return -1
	}
}

// ClassTermStrToStr converts a string like "1" to "冬学期", used for config
func ClassTermStrToStr(str string) string {
	switch str {
	case "0":
		return "秋学期"
	case "1":
		return "冬学期"
	case "2":
		return "短学期A"
	case "3":
		return "小学期"
	case "4":
		return "春学期"
	case "5":
		return "夏学期"
	case "6":
		return "短学期B"
	default:
		return "未知学期"
	}
}

// ExamStrToStr converts a string like "1" to "春夏学期", used for config
func ExamStrToStr(str string) string {
	switch str {
	case "0":
		return "秋冬学期"
	case "1":
		return "春夏学期"
	default:
		return "未知学期"
	}
}
