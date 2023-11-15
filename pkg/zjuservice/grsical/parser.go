package grsical

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"ugrs-ical/pkg/zjuservice/zjuconst"

	"io"
	"strconv"
	"strings"
	"time"
)

const GrsTableRowNum = 16
const UgrsTableRowNum = 15

func GetTable(r io.Reader) (*html.Node, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	table := doc.Find(".table-course").Find("tbody")
	if len(table.Nodes) == 0 {
		return nil, fmt.Errorf("can not find table")
	}
	return table.Nodes[0], nil
}

func GetExamTable(r io.Reader) (*html.Node, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	table := doc.Find("table").Find("tbody")
	if len(table.Nodes) == 0 {
		return nil, fmt.Errorf("can not find table")
	}
	return table.Nodes[0], nil
}

func parseClass(ctx context.Context, a *html.Node) (zjuconst.ZJUClass, error) {
	class := zjuconst.ZJUClass{}
	children := getValidChildren(a)
	if len(children) == 6 {
		// 单双周的字体被特殊标记，导致第二个元素分裂
		if children[2].FirstChild != nil {
			children[1].Data += children[2].FirstChild.Data
			children = append(children[:2], children[3:]...)
		}
	}
	if len(children) != 5 {
		return class, fmt.Errorf("invalid children length")
	}
	// Name
	if children[0].FirstChild == nil {
		return class, fmt.Errorf("can not get class name")
	}
	class.ClassName = children[0].FirstChild.Data
	// Semester & Repeat
	tr := strings.Split(children[1].Data, "||")
	if len(tr) != 2 {
		return class, fmt.Errorf("can not parse time and repeat: %s", children[1].Data)
	}
	switch strings.TrimSpace(tr[0]) {
	case "秋":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Autumn)
	case "冬":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Winter)
	case "秋冬":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Autumn, zjuconst.Winter)
	case "春":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Spring)
	case "夏":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Summer)
	case "春夏":
		class.TermArrangements = append(class.TermArrangements, zjuconst.Spring, zjuconst.Summer)
	default:
		return class, fmt.Errorf("invalid semester: %s", children[1].Data)
	}
	switch strings.TrimSpace(tr[1]) {
	case "每周":
		class.WeekArrangement = zjuconst.Normal
	case "单周":
		class.WeekArrangement = zjuconst.OddOnly
	case "双周":
		class.WeekArrangement = zjuconst.EvenOnly
	default:
		// 不是严重错误
		class.WeekArrangement = zjuconst.Normal
		log.Ctx(ctx).Warn().Msgf("unsupported repeat pattern: %s", strings.TrimSpace(tr[1]))
	}
	// RawDuration
	//var b strings.Builder
	//b.Grow(len(children[2].Data))
	//for _, ch := range children[2].Data {
	//	if !unicode.IsSpace(ch) {
	//		b.WriteRune(ch)
	//	}
	//}
	//class.RawDuration = b.String()
	class.TeacherName = strings.TrimSpace(children[3].Data)
	class.ClassLocation = strings.TrimSpace(children[4].Data)
	//TODO
	class.ClassCode = "未知"
	return class, nil
}

func findChildren(node *html.Node, atom atom.Atom) []*html.Node {
	var res []*html.Node
	n := node.FirstChild
	if n == nil {
		return nil
	}
	nn := &html.Node{
		NextSibling: n,
	}
	for nn.NextSibling != nil {
		if nn.NextSibling.DataAtom == atom {
			res = append(res, nn.NextSibling)
		}
		nn = nn.NextSibling
	}
	return res
}

func getValidChildren(node *html.Node) []*html.Node {
	n := &html.Node{
		NextSibling: node.FirstChild,
	}
	if n.NextSibling == nil {
		return nil
	}

	var c []*html.Node
	for n.NextSibling != nil {
		// 如果是空的文本节点或br节点
		if (n.NextSibling.Type == html.TextNode &&
			n.NextSibling.FirstChild == nil &&
			strings.TrimSpace(n.NextSibling.Data) == "") ||
			n.NextSibling.DataAtom == atom.Br {
			n = n.NextSibling
			continue
		}
		c = append(c, n.NextSibling)
		n = n.NextSibling
	}
	return c
}

func ParseTable(ctx context.Context, node *html.Node, isUGRS bool) ([]zjuconst.ZJUClass, error) {
	var rowNum int
	if isUGRS {
		rowNum = UgrsTableRowNum
	} else {
		rowNum = GrsTableRowNum
	}
	trs := goquery.NewDocumentFromNode(node).Children()
	if trs.Length() > rowNum {
		log.Ctx(ctx).Warn().Msgf("too many tr element: %d", trs.Length())
	} else if trs.Length() < rowNum {
		log.Ctx(ctx).Error().Msgf("insufficient tr element: %d", trs.Length())
		return nil, fmt.Errorf("insufficient tr element")
	}

	var classes []zjuconst.ZJUClass
	mask := [7][15]bool{}
	for i := 1; i < rowNum; i++ {
		tr := trs.Nodes[i]
		tds := goquery.NewDocumentFromNode(tr).Children().Nodes
		if i == 1 || i == 6 || i == 11 {
			// 删除时间的第一列（上午、下午、晚上）
			tds = tds[1:]
		}
		// 删除节数信息
		tds = tds[1:]

		j := 0
		for _, td := range tds {
			j++
			for mask[j-1][i-1] {
				j++
			}
			aChildren := findChildren(td, atom.A)
			if len(aChildren) == 0 {
				mask[j-1][i-1] = true
				continue
			}

			var err error
			span := 1
			for _, attr := range td.Attr {
				if attr.Key == "rowspan" {
					span, err = strconv.Atoi(attr.Val)
					if err != nil || span < 1 || span > 15 {
						log.Ctx(ctx).Error().Msgf("failed to parse rowspan, val was %s", attr.Val)
					}
				}
				for k := 0; k < span; k++ {
					mask[j-1][i-1+k] = true
				}
			}
			for _, a := range aChildren {
				class, err := parseClass(ctx, a)
				if err != nil {
					log.Ctx(ctx).Warn().Msgf("failed to parse class: %s", err.Error())
					continue
				}
				class.StartPeriod = i
				class.EndPeriod = i + span - 1
				class.DayNumber = j
				classes = append(classes, class)
			}
		}
	}

	return classes, nil
}

func getFirstChildData(node *html.Node, def string) string {
	if node.FirstChild != nil && node.FirstChild.Data != "" {
		return node.FirstChild.Data
	} else {
		return def
	}
}

func ParseExamTable(ctx context.Context, node *html.Node) ([]zjuconst.ZJUExam, error) {
	var err error
	trs := goquery.NewDocumentFromNode(node).Find("tr").Nodes
	var exams []zjuconst.ZJUExam
	for _, tr := range trs {
		tds := goquery.NewDocumentFromNode(tr).Find("td").Nodes
		l := len(tds)
		if l > 9 {
			err = errors.New("long node")
			l = 9
		} else if l < 6 {
			err = errors.New("short node")
			continue
		}

		exam := zjuconst.ZJUGrsExam{}
		valid := true
	FOR:
		for i := 0; i < l; i++ {
			switch i {
			case 0:
				exam.Semester = getFirstChildData(tds[i], "未知")
			case 1:
				exam.ID = getFirstChildData(tds[i], "未知课号")
			case 2:
				exam.ClassName = getFirstChildData(tds[i], "未知课程")
			case 3:
				exam.Region = getFirstChildData(tds[i], "")
			case 4:
				t2 := strings.Split(getFirstChildData(tds[5], ""), "->")
				if len(t2) < 2 {
					err = errors.New("malformed time 1")
					valid = false
					break FOR
				}
				d := getFirstChildData(tds[4], "")
				t0 := fmt.Sprintf("%s %s", d, t2[0])
				t1 := fmt.Sprintf("%s %s", d, t2[1])
				exam.StartTime, err = time.Parse("2006-01-02 15:04", t0)
				if err != nil {
					valid = false
					break FOR
				}
				exam.EndTime, err = time.Parse("2006-01-02 15:04", t1)
				if err != nil {
					valid = false
					break FOR
				}
			case 6:
				exam.ExamLocation = getFirstChildData(tds[i], "未知地点")
			case 7:
				exam.ExamSeatNum = getFirstChildData(tds[i], "")
			case 8:
				exam.ExamRemark = getFirstChildData(tds[i], "")
			}
		}
		if valid {
			exams = append(exams, &exam)
		}
	}
	return exams, err
}
