package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/easylifewell/purifier-server/model"
	"github.com/easylifewell/purifier-server/store"
	"github.com/gorilla/mux"
)

type DataController struct {
}

func NewDataController() *DataController {
	return &DataController{}
}

func (dc DataController) Index(w http.ResponseWriter, r *http.Request) {
	datas, err := store.GetDatas()
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

func (dc DataController) DataCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Please use POST method")
		if err != nil {
			panic(err)
		}
		return
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
		return
	}

	d, err := createData(body, token)
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

func (dc DataController) DataShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["deviceID"]
	datas, err := store.GetDataByDeviceID(deviceID)
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

func createData(body []byte, token string) (model.Data, error) {
	var data model.Data
	key1 := "Edb^@u2T"
	key2 := "aZEqm5ph"
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	if data.Version != "0.1" {
		return data, errors.New("bad version data")
	}

	if len(data.DeviceID) != 8 {
		return data, errors.New("DeviceID must be 8 length of digital")
	}
	t := "VIEW" + key1 + data.DeviceID + key2
	if t != token {
		return data, errors.New("TOKEN is invalid")
	}
	data.Time = time.Now()

	_, err = store.AddData(data)
	if err != nil {
		return data, err
	}
	return data, nil
}
