package main

import "fmt"

var currentID int
var datas Datas

// Give us some seed data
func init() {
	RepoCreateData(Data{Value: "default vaule1"})
	RepoCreateData(Data{Value: "default vaule2"})
}

func RepoFindData(id int) Data {
	for _, d := range datas {
		if d.DeviceID == id {
			return d
		}
	}
	// return empty Data if not found
	return Data{}
}

func RepoCreateData(d Data) Data {
	currentID++
	d.DeviceID = currentID
	datas = append(datas, d)
	return d
}

func RepoDestroyData(id int) error {
	for i, d := range datas {
		if d.DeviceID == id {
			datas = append(datas[:i], datas[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("could not find Data with id of %d to delete", id)
}
