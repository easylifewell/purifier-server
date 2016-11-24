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
		router.PathPrefix("/app").Handler(http.StripPrefix("/app", http.FileServer(http.Dir("./view"))))
		// Static Files
		router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

		// Get a DataController instance.
		dc := controller.NewDataController()
		router.HandleFunc("/", dc.Index)
		router.HandleFunc("/data", dc.DataCreate)
		router.HandleFunc("/data/{deviceID}", dc.DataShow)

		// Get a SMSController instance
		sc := controller.NewSMSController()
		// 正常使用密码登录
		router.HandleFunc("/api/login", sc.Login).
			Queries("phone", "{phone:[0-9]+}", "password", "{password:.{6,20}}")
		//  发送登录验证码给手机号码
		router.HandleFunc("/api/sms/{phone:[0-9]+}", sc.SendSMS)
		//  注册并登录进行，此时保存密码
		router.HandleFunc("/api/register", sc.Register).
			Queries("phone", "{phone:[0-9]+}", "smscode", "{smscode:[0-9]+}", "password", "{password:.{6,20}}")

		// Get a UserController instance.
		uc := controller.NewUserController()
		router.HandleFunc("/api/user", uc.GetUser)
		router.HandleFunc("/api/user/setting", uc.SetNickName).Queries("nickname", "{nickname}")
		router.HandleFunc("/api/user/setting", uc.SetRealName).Queries("realname", "{realname}")

		// Get a DeviceController instance
		deviceController := controller.NewDeviceController()
		router.HandleFunc("/api/device/add", deviceController.BindDeviceWithUser).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}")
		router.HandleFunc("/api/device/started", deviceController.IsStarted).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}")
		router.HandleFunc("/api/device/on", deviceController.On).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}")
		router.HandleFunc("/api/device/off", deviceController.Off).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}")
		router.HandleFunc("/api/device/setting", deviceController.SetCarName).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}", "carname", "{carname}")
		router.HandleFunc("/api/device/setting", deviceController.SetSPhone).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}", "sphone", "{sphone[0-9]+}")
		router.HandleFunc("/api/device/setting", deviceController.SetRPhone).
			Queries("deviceid", "{deviceid:[a-zA-Z0-9]{10}}", "rphone", "{rphone[0-9]}+")

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
