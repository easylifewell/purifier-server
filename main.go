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
		router.HandleFunc("/data/{deviceID}", DataShow)
		logrus.Fatal(http.ListenAndServe("0.0.0.0:6060", router))
		return nil
	}

	app.Run(os.Args)
}

func Index(w http.ResponseWriter, r *http.Request) {
	datas, err := GetDatas()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Get datas from database failed")
		if err != nil {
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(datas); err != nil {
		panic(err)
	}
}

func DataCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Please use POST method")
		if err != nil {
			panic(err)
		}
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	token := r.Header.Get("TOKEN")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Please add TOKEN in http Header")
		if err != nil {
			panic(err)
		}
	}

	d, err := RepoCreateData(body, token)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, fmt.Sprintf("ERROR: %v\n", err))
		if err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}
}

func DataShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["deviceID"]
	datas, err := GetDataByDeviceID(deviceID)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Get data from database failed")
		if err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(datas); err != nil {
		panic(err)
	}
}
