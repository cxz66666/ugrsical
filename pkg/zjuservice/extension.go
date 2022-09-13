package zjuservice

func ClassTermToQueryString(t ClassTerm) string {
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

func ExamTermToQueryString(t ExamTerm) string {
	switch t {
	case AutumnWinter:
		return "1"
	case SpringSummer:
		return "2"
	default:
		return ""
	}
}

// ClassTermStrToStr converts a string like "1" to "冬学期"
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

// ExamStrToStr converts a string like "1" to "春夏学期"
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
