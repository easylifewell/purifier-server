package main

import (
	"encoding/json"
	"errors"
	"time"
)

func RepoCreateData(body []byte, token string) (Data, error) {
	var data Data
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

	_, err = AddData(data)
	if err != nil {
		return data, err
	}
	return data, nil
}
