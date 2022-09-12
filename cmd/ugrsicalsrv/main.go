package main

import (
	"os"

	"ugrs-ical/internal/ugrsicalsrv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	//不想写日志文件了，好麻烦，stderr凑合看一下吧
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Fatal().Msg(ugrsicalsrv.ListenAndServe().Error())
}
