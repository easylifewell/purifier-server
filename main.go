package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"

	"github.com/easylifewell/purifier-server/controller"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "purifier-server"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug output for logging",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "text",
			Usage: "set the format used by logs ('text' (default), or 'json')",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "6060",
			Usage: "set the port the server listen on(default: 6060)",
		},
	}
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		switch context.GlobalString("log-format") {
		case "text":
			// retain logrus's default.
		case "json":
			logrus.SetFormatter(new(logrus.JSONFormatter))
		default:
			return fmt.Errorf("unknown log-format %q", context.GlobalString("log-format"))
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		router := mux.NewRouter().StrictSlash(true)

		// View Files
		router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("./view"))))
		// Static Files
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

		// Get a DataController instance.
		dc := controller.NewDataController()
		router.HandleFunc("/", dc.Index)
		router.HandleFunc("/data", dc.DataCreate)
		router.HandleFunc("/data/{deviceID}", dc.DataShow)

		// Get a SMSController instance
		sc := controller.NewSMSController()
		//  登录的第一个步骤，发送登录验证码给手机号码
		router.HandleFunc("/api/login/{phone:[0-9]+}", sc.SendSMS)
		//  登录的第二个步骤，进行验证
		router.HandleFunc("/api/login/checkSms", sc.CheckSMS).Queries("phone", "{phone:[0-9]+}", "smscode", "{smscode:[0-9]+}")

		// Get a UserController instance.
		uc := controller.NewUserController()
		router.HandleFunc("/api/user", uc.GetUser)
		router.HandleFunc("/api/user/setting", uc.SetNickName).Queries("nickname", "{nickname}")
		router.HandleFunc("/api/user/setting", uc.SetRealName).Queries("realname", "{realname}")

		// Start the server
		addr := fmt.Sprintf("0.0.0.0:%s", c.GlobalString("port"))
		logrus.Infof("Server started and Listened on %s", addr)
		logrus.Fatal(http.ListenAndServe(addr, router))
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
