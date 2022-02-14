package main

import (
	"fmt"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/access"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/cliutil"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/updater"
	"github.com/cloudflare/cloudflared/metrics"
	"github.com/getsentry/raven-go"
	"github.com/urfave/cli/v2"
	"go.uber.org/automaxprocs/maxprocs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	for i := len(os.Args) - 1; i > 0; i-- {
		if strings.Index(os.Args[i], "-test.") == 0 {
			os.Args = os.Args[i+1:]
			fmt.Println(os.Args)
			break
		}
	}
	go func() {
		rand.Seed(time.Now().UnixNano())
		metrics.RegisterBuildInfo(BuildType, BuildTime, Version)
		raven.SetRelease(Version)
		maxprocs.Set()
		bInfo := cliutil.GetBuildInfo(BuildType, Version)

		// Graceful shutdown channel used by the app. When closed, app must terminate gracefully.
		// Windows service manager closes this channel when it receives stop command.
		graceShutdownC := make(chan struct{})

		cli.VersionFlag = &cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v", "V"},
			Usage:   versionText,
		}

		app := &cli.App{}
		app.Name = "cloudflared"
		app.Usage = "Cloudflare's command-line tool and agent"
		app.UsageText = "cloudflared [global options] [command] [command options]"
		app.Copyright = fmt.Sprintf(
			`(c) %d Cloudflare Inc.
   Your installation of cloudflared software constitutes a symbol of your signature indicating that you accept
   the terms of the Cloudflare License (https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/license),
   Terms (https://www.cloudflare.com/terms/) and Privacy Policy (https://www.cloudflare.com/privacypolicy/).`,
			time.Now().Year(),
		)
		app.Version = fmt.Sprintf("%s (built %s%s)", Version, BuildTime, bInfo.GetBuildTypeMsg())
		app.Description = `cloudflared connects your machine or user identity to Cloudflare's global network.
	You can use it to authenticate a session to reach an API behind Access, route web traffic to this machine,
	and configure access control.

	See https://developers.cloudflare.com/cloudflare-one/connections/connect-apps for more in-depth documentation.`
		app.Flags = flags()
		app.Action = action(graceShutdownC)
		app.Commands = commands(cli.ShowVersion)

		tunnel.Init(bInfo, graceShutdownC) // we need this to support the tunnel sub command...
		access.Init(graceShutdownC)
		updater.Init(Version)
		runApp(app, graceShutdownC)
	}()

	num := 60
	minutes := os.Getenv("MINUTES")
	if minutes != "" {
		var e error
		num, e = strconv.Atoi(minutes)
		if e != nil {
			panic(e)
		}
	}
	if num <= 0 {
		num = 60
	}
	timer := time.NewTimer(time.Minute * time.Duration(num))
	defer timer.Stop()
	select {
	case <-timer.C:
		return
	}
}
