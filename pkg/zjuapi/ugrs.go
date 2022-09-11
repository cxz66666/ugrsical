package zjuapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/rs/zerolog/log"
)

const GrsLoginUrl = "http://grs.zju.edu.cn/ssohome"
const changeLocaleUrl = "http://grs.zju.edu.cn/py/page/student/grkcb.htm?pageAction=changeLocale"

type GrsSemester int

const (
	Spring GrsSemester = 11
	Summer GrsSemester = 12
	Autumn GrsSemester = 13
	Winter GrsSemester = 14
)

func (c *ZJUAPIClient) FetchTimetable(ctx context.Context, year int, semester GrsSemester) (io.Reader, error) {
	// year - year+1 学年度
	_, err := c.HttpClient.PostForm(changeLocaleUrl, url.Values{
		"locale": {"zh_CN"},
	})
	r, err := c.HttpClient.Get(fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grkcb.htm?xj=%d&xn=%d", semester, year))
	if err != nil {
		e := fmt.Sprintf("failed to fetch timetable for %d-%d, error: %s", year, semester, err)
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read timetable for %d-%d, error: %s", year, semester, err)
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	return bytes.NewBuffer(rb), nil
}
