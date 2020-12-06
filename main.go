package main

import (
	"fmt"
	"github.com/inconshreveable/go-update"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/robfig/cron.v3"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	appConst "qrcode-server/src/appConst"
	"qrcode-server/src/webServer"
	"syscall"
	"time"
)

var (
	AppName      string
	AppVersion   string
	BuildVersion string
	BuildTime    string
	GoVersion    string
	GitBranch    string
)

func main() {
	app := &cli.App{
		Name:    AppName,
		Usage:   "Smart ID proxy",
		Version: AppVersion,
		Commands: []*cli.Command{
			{
				Name:    "info",
				Usage:   "info build",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					Version()
					return nil
				},
			},
			{
				Name:    "server",
				Usage:   "server --listen 0.0.0.0:80",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "listen, 0.0.0.0:80",
						Aliases: []string{"l"},
						Value:   "0.0.0.0:80",
					},
				},
				Action: func(c *cli.Context) error {

					listen := c.String("listen")

					app, err := appConst.GetApp(AppVersion)
					if err != nil {
						return err
					}
					app.SetAddr(listen)

					webServer, err := webServer.New()
					if err != nil {
						return err
					}

					log.WithFields(log.Fields{
						"package":  "main",
						"function": "main",
						"error":    err,
						//"notifier": viper,
					}).Info("Server start ", app.GetServerFullAddr())

					userInterrupt(webServer)
					if err := webServer.StartServer(); err != nil {
						fmt.Println(err)
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}
}

func Version() {
	fmt.Printf("App Name:\t%s\n", AppName)
	fmt.Printf("App Version:\t%s\n", AppVersion)
	fmt.Printf("Git branch:\t%s\n", GitBranch)
}

func userInterrupt(proxy *webServer.WebServer) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := proxy.Shutdown(); err != nil {
			fmt.Println(err.Error())
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "userInterrupt",
				"error":    err,
				//"notifier": viper,
			}).Error("Server stop on demand")
			panic(err)
		}
		close(c)
	}()
}

func getUpdate(url string, proxy *webServer.WebServer, c *cron.Cron) func() {
	//fmt.Println(runtime.GOOS)
	return func() {
		//fmt.Println("cron run")
		resp, err := http.Get(url)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "getUpdate",
				"error":    err,
				//"notifier": viper,
			}).Error("Update resource unavailable ", url)
			return
		}
		defer resp.Body.Close()
		err_ := update.Apply(resp.Body, update.Options{})
		if err_ != nil {
			// error handling
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "getUpdate",
				"error":    err,
				//"notifier": viper,
			}).Error("Failed to update file")
		}

		//os.Args

		go func() {
			time.Sleep(5 * time.Second)
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "getUpdate",
				//"notifier": viper,
			}).Info("Process restart after upgrade")
			cmd := exec.Command(os.Args[0], os.Args[1:]...)
			err = cmd.Start()
			proxy.Shutdown()
			c.Stop()
		}()

		//os.Exit(0)
		return
	}
}
