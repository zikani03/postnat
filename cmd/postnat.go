package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zikani03/postnat"
)

type Globals struct {
	Config  string      `help:"Location of configuration file" default:"postnat.toml" type:"path"`
	Debug   bool        `help:"Enable debug mode"`
	Version VersionFlag `name:"version" help:"Show version and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type RunCmd struct {
	Daemonize bool `help:"Daemonize or run in foreground"`
}

func (r *RunCmd) Run(globals *Globals) error {
	config, err := postnat.ParseConfig(globals.Config)
	if err != nil {
		return err
	}

	log.Debug().
		Str("config", fmt.Sprint(config)).
		Msg("Starting")

	app, err := postnat.New(config)
	if err != nil {
		return err
	}
	// TODO: if r.Daemonize
	<-app.Run()

	app.Shutdown()
	return nil
}

type CLI struct {
	Globals

	Run RunCmd `cmd:"" help:"Start the postnat daemon"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cli := CLI{
		Globals: Globals{
			Version: VersionFlag("1.0.0"),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("postnat"),
		kong.Description("Publish messages to NATS from PostgreSQL LISTEN/NOTIFY messages"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "1.0.0",
		})

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
