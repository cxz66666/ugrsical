package grsical

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	common2 "grs-ical/internal/common"

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
	tweaksFile   string
	outputFile   string
	forceWrite   bool
	rootCmd      = &cobra.Command{
		Use:           "ugrsical -u username -p password -c config [-t tweak] [-o output] [-f]",
		Short:         "ugrsical is a tool for generating class schedules iCalendar file",
		Long:          `A command-line utility for generating class schedule iCalender file from extracting data from ZJU Graduate School web pages.`,
		SilenceErrors: true,
		RunE:          CliMain,
	}
	version = "dirty"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "ZJUAM username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "ZJUAM password")
	rootCmd.PersistentFlags().StringVarP(&userPassFile, "upfile", "i", "", "username and password json")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "config file")
	rootCmd.PersistentFlags().StringVarP(&tweaksFile, "tweak", "t", "", "tweaks file")
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
	if userName == "" && password == "" {
		return errors.New("no username or password set, exiting")
	}

	c, err := common2.GetConfig(configFile)
	tc, err := common2.GetTweakConfig(tweaksFile)

	if _, err := os.Stat(outputFile); !errors.Is(err, os.ErrNotExist) && !forceWrite {
		return errors.New("output file exists, exiting")
	}
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := common2.FetchToMemory(ctx, userName, password, c, tc)
	if err != nil {
		return err
	}
	_, err = f.WriteString(r)
	if err != nil {
		return err
	}

	return nil
}

func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
