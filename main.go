package main

import (
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
	app.Action = func(c *cli.Context) error {
		router := mux.NewRouter().StrictSlash(true)
		// Get a DataController instance.
		dc := controller.NewDataController()
		// Get a UserController instance.
		// uc := controller.NewUserController()
		// Get a SMSController instance
		sc := controller.NewSMSController()

		router.HandleFunc("/", dc.Index)
		router.HandleFunc("/data", dc.DataCreate)
		router.HandleFunc("/data/{deviceID}", dc.DataShow)

		//  登录的第一个步骤，发送登录验证码给手机号码
		router.HandleFunc("/api/login/{phone}", sc.SendSMS)
		router.HandleFunc("/api/login/checkSms/{phone}/{smscode}", sc.CheckSMS)

		logrus.Fatal(http.ListenAndServe("0.0.0.0:6060", router))
		return nil
	}

	app.Run(os.Args)
}
