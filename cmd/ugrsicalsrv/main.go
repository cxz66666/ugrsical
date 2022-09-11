package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"grs-ical/internal/grsicalsrv"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Fatal().Msg(grsicalsrv.ListenAndServe(":3000").Error())
}
