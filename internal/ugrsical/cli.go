package ugrsical

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	common2 "ugrs-ical/internal/common"
	"ugrs-ical/pkg/ical"
	"ugrs-ical/pkg/zjuservice"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type pwFile struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	userName     string
	password     string
	userPassFile string
	configFile   string
	icalType     int
	outputFile   string
	forceWrite   bool
	rootCmd      = &cobra.Command{
		Use:           "ugrsical -u username -p password [-c config] [-t 0] [-o output] [-f]",
		Short:         "ugrsical is a tool for generating class schedules iCalendar file",
		Long:          `A command-line utility for generating class schedule iCalender file from extracting data from ZJU DingDing API.`,
		SilenceErrors: true,
		RunE:          CliMain,
	}
	version = "dirty"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "ZJUAM username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "ZJUAM password")
	rootCmd.PersistentFlags().StringVarP(&userPassFile, "upfile", "i", "", "username and password json")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", zjuservice.ConfigDefaultPath, "config file")
	rootCmd.PersistentFlags().IntVarP(&icalType, "type", "t", 0, "0(default) for both class and exam, 1 for only class, 2 for only exam, 3 for only scores")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "ugrsical.ics", "output file")
	rootCmd.PersistentFlags().BoolVarP(&forceWrite, "force", "f", false, "force write to target file")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("ugrsical build {{.Version}}\n")
}

func CliMain(cmd *cobra.Command, args []string) error {
	ctx := log.With().Str("reqid", uuid.NewString()).Logger().WithContext(context.Background())

	if userPassFile != "" {
		upf, err := os.Open(userPassFile)
		defer upf.Close()
		if err != nil {
			return err
		}
		upfc, err := io.ReadAll(upf)
		if err != nil {
			return err
		}
		var up pwFile
		err = json.Unmarshal(upfc, &up)
		userName = up.Username
		password = up.Password
	}

	if userName == "" || password == "" {
		return errors.New("no username or password set, exiting")
	}
	if icalType < 0 || icalType > 2 {
		return errors.New("invalid ical type")
	}
	// Load config before fetch !
	if err := zjuservice.LoadConfig(configFile); err != nil {
		return err
	}

	if _, err := os.Stat(outputFile); !errors.Is(err, os.ErrNotExist) && !forceWrite {
		return errors.New("output file exists, exiting")
	}
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}
	var vCal ical.VCalendar
	var icalName string
	switch icalType {
	case 0:
		vCal, err = common2.GetBothCalendar(ctx, userName, password)
		if err != nil {
			return err
		}
		icalName = ""
	case 1:
		vCal, err = common2.GetClassCalendar(ctx, userName, password)
		if err != nil {
			return err
		}
		icalName = "UGRSICAL 课程表"
	case 2:
		vCal, err = common2.GetExamCalendar(ctx, userName, password)
		if err != nil {
			return err
		}
		icalName = "UGRSICAL 考试表"
	case 3:
		vCal, err = common2.GetScoreCalendar(ctx, userName, password)
		if err != nil {
			return err
		}
		icalName = "UGRSICAL GPA表"
	}

	_, err = f.WriteString(vCal.GetICS(icalName))
	if err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msg("generate success!")
	return nil
}

func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
