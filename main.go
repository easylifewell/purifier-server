package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "purifier-server"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Action = func(c *cli.Context) error {
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", Index)
		router.HandleFunc("/data", DataCreate)
		router.HandleFunc("/data/{dataId}", DataShow)
		logrus.Fatal(http.ListenAndServe("0.0.0.0:6060", router))
		return nil
	}

	app.Run(os.Args)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(datas); err != nil {
		panic(err)
	}
}

func DataCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	d := RepoCreateData(Data{Value: string(body[:])})
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}

}

func DataShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataID := vars["dataId"]
	fmt.Fprintln(w, "Data Show:", dataID)
}
