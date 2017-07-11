package commands

import (
	"strconv"

	"github.com/aisk/chrysanthemum"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/leancloud/lean-cli/api"
	"github.com/leancloud/lean-cli/apps"
	"github.com/leancloud/lean-cli/console"
)

func debugAction(c *cli.Context) error {
	remote := c.String("remote")
	port := strconv.Itoa(c.Int("port"))

	appID, err := apps.GetCurrentAppID(".")
	if err != nil {
		return newCliError(err)
	}

	bar := chrysanthemum.New("获取应用信息").Start()
	appInfo, err := api.GetAppInfo(appID)
	if err != nil {
		bar.Failed()
		return newCliError(err)
	}
	bar.Successed()
	chrysanthemum.Printf("当前应用：%s (%s)\r\n", color.RedString(appInfo.AppName), appID)

	cons := &console.Server{
		AppID:       appInfo.AppID,
		AppKey:      appInfo.AppKey,
		MasterKey:   appInfo.MasterKey,
		HookKey:     appInfo.HookKey,
		RemoteURL:   remote,
		ConsolePort: port,
		Errors:      make(chan error),
	}

	cons.Run()
	for {
		select {
		case err = <-cons.Errors:
			panic(err)
		}
	}
}
