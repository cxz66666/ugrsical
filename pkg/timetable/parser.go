package timetable

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strconv"
	"strings"
	"unicode"
)

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

func parseClass(a *html.Node) (Class, error) {
	class := Class{}
	children := getValidChildren(a)
	if len(children) != 5 {
		return class, fmt.Errorf("invalid children length")
	}
	// Name
	if children[0].FirstChild == nil {
		return class, fmt.Errorf("can not get class name")
	}
	class.Name = children[0].FirstChild.Data
	// Term & Repeat
	tr := strings.Split(children[1].Data, "||")
	if len(tr) != 2 {
		return class, fmt.Errorf("can not parse time and repeat: %s", children[1].Data)
	}
	switch strings.TrimSpace(tr[0]) {
	case "秋":
		class.Term = Autumn
	case "冬":
		class.Term = Winter
	case "秋冬":
		class.Term = AutumnWinter
	case "春":
		class.Term = Spring
	case "夏":
		class.Term = Summer
	case "春夏":
		class.Term = SpringSummer
	default:
		return class, fmt.Errorf("invalid term: %s", children[1].Data)
	}
	switch strings.TrimSpace(tr[1]) {
	case "每周":
		class.Repeat = EveryWeek
	default:
		// 不是严重错误
		class.Repeat = EveryWeek
		log.Warn().Msgf("unsupported repeat pattern: %s", strings.TrimSpace(tr[1]))
	}
	// RawDuration
	var b strings.Builder
	b.Grow(len(children[2].Data))
	for _, ch := range children[2].Data {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	class.RawDuration = b.String()
	class.Teacher = strings.TrimSpace(children[3].Data)
	class.Location = strings.TrimSpace(children[4].Data)
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

func ParseTable(node *html.Node) (*[]Class, error) {
	trs := goquery.NewDocumentFromNode(node).Children()
	if trs.Length() > 16 {
		log.Warn().Msgf("too many tr element: %d", trs.Length())
	} else if trs.Length() < 16 {
		log.Error().Msgf("insufficient tr element: %d", trs.Length())
		return nil, fmt.Errorf("insufficient tr element")
	}

	var classes []Class
	mask := [7][15]bool{}
	for i := 1; i < 16; i++ {
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
						log.Error().Msgf("failed to parse rowspan, val was %s", attr.Val)
					}
				}
				for k := 0; k < span; k++ {
					mask[j-1][i-1+k] = true
				}
			}
			for _, a := range aChildren {
				class, err := parseClass(a)
				if err != nil {
					log.Error().Msgf("failed to parse class: %s", err.Error())
					return nil, nil
				}
				class.Duration.Starts = i
				class.Duration.Ends = i + span - 1
				class.DayOfWeek = j
				classes = append(classes, class)
			}
		}
	}

	return &classes, nil
}
