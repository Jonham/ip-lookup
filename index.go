package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"sort"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

type webServiceConf struct {
	Host string
	Port string
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mode, m",
			Value: "server",
			Usage: "`server`(default) run app as service, `client` will run interval sending local to server.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "client",
			Aliases: []string{"m"},
			Usage:   "run as client mode",
			Action: func(c *cli.Context) error {
				// 定时提交Get网络下的IP
				startInterval(sendGetNetIP, 1, 0) // 1 min
				return nil
			},
		},
		{
			Name:    "dev",
			Aliases: []string{"m"},
			Usage:   "run as client mode `DEV`",
			Action: func(c *cli.Context) error {
				// 定时提交Get网络下的IP
				startInterval(sendGetNetIP, 0, 5) // 5 second
				return nil
			},
		},
		{
			Name:  "release",
			Usage: "run Gin as `release` mode",
			Action: func(c *cli.Context) error {
				gin.SetMode(gin.ReleaseMode)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Name = "IP Lookup Service"
	app.Usage = "default run as service"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Command Not Found. Bye!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadFile("config.yml")
	config := webServiceConf{}
	yaml.Unmarshal(data, &config)

	r := setupRouter()

	r.Run(config.Host + ":" + config.Port)
}
