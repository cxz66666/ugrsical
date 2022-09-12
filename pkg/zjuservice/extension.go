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
