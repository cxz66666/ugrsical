package common

import (
	"context"
	"github.com/cxz66666/zju-ical/pkg/zjuservice"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/grsical"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/newgrsical"
	"github.com/cxz66666/zju-ical/pkg/zjuservice/ugrsical"
)

const ContextReqId = "zjuical-reqid"

var useNewGrsService = false

var useNewUgrsService = false

func newGrsService(ctx context.Context, isUgrs bool) zjuservice.IZJUService {
	if useNewGrsService {
		return newgrsical.NewNewGrsService(ctx, isUgrs)
	} else {
		return grsical.NewGrsService(ctx, isUgrs)
	}
}

func newUgrsService(ctx context.Context) zjuservice.IZJUService {
	return ugrsical.NewUgrsService(ctx)
}

func SetServiceAPI(newGrs, newUgrs bool) {
	useNewGrsService = newGrs
	useNewUgrsService = newUgrs
}
