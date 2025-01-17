package main

import (
	"context"
	"fmt"
	log "github.com/kataras/golog"
	"github.com/urfave/cli/v2"
	"github.com/zema1/suo5/ctrl"
	"os"
	"os/signal"
	"strings"
)

func main() {
	log.Default.SetTimeFormat("01-02 15:04")
	app := cli.NewApp()
	app.Name = "suo5"
	app.Usage = "A super http proxy tunnel"
	app.Version = "v0.2.0"

	defaultConfig := ctrl.DefaultSuo5Config()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "target",
			Aliases:  []string{"t"},
			Usage:    "set the memshell url, ex: http://localhost:8080/tomcat_debug_war_exploded/",
			Value:    defaultConfig.Target,
			Required: true,
		},
		&cli.StringFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "set the socks server port",
			Value:   defaultConfig.Listen,
		},
		&cli.BoolFlag{
			Name:  "no-auth",
			Usage: "disable socks5 authentication",
			Value: defaultConfig.NoAuth,
		},
		&cli.StringFlag{
			Name:  "auth",
			Usage: "socks5 creds, username:password, leave empty to auto generate",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "mode",
			Usage: "connection mode, choices are auto, full, half",
			Value: string(defaultConfig.Mode),
		},
		&cli.StringFlag{
			Name:  "ua",
			Usage: "the user-agent used to send request",
			Value: defaultConfig.UserAgent,
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "http request timeout in seconds",
			Value: defaultConfig.Timeout,
		},
		&cli.IntFlag{
			Name:  "buf-size",
			Usage: "set the request max body size",
			Value: defaultConfig.BufferSize,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "debug the traffic, print more details",
			Value:   defaultConfig.Debug,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.Default.SetLevel("debug")
		}
		return nil
	}
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {
	listen := c.String("listen")
	target := c.String("target")
	noAuth := c.Bool("no-auth")
	auth := c.String("auth")
	mode := ctrl.ConnectionType(c.String("mode"))
	ua := c.String("ua")
	bufSize := c.Int("buf-size")
	timeout := c.Int("timeout")
	debug := c.Bool("debug")

	var username, password string
	if auth == "" {
		username = "suo5"
		password = ctrl.RandString(8)
	} else {
		parts := strings.Split(auth, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid socks credentials, expected username:password")
		}
		username = parts[0]
		password = parts[1]
	}
	if !(mode == ctrl.AutoDuplex || mode == ctrl.FullDuplex || mode == ctrl.HalfDuplex) {
		return fmt.Errorf("invalid mode, expected auto or full or half")
	}

	if bufSize < 512 || bufSize > 1024000 {
		return fmt.Errorf("inproper buffer size, 512~1024000")
	}
	config := &ctrl.Suo5Config{
		Listen:     listen,
		Target:     target,
		NoAuth:     noAuth,
		Username:   username,
		Password:   password,
		Mode:       mode,
		UserAgent:  ua,
		BufferSize: bufSize,
		Timeout:    timeout,
		Debug:      debug,
	}
	ctx, cancel := signalCtx()
	defer cancel()
	return ctrl.Run(ctx, config)
}

func signalCtx() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		cancel()
	}()
	return ctx, cancel
}
