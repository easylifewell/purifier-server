package store

import (
	"github.com/easylifewell/purifier-server/model"

	"golang.org/x/net/context"
)

type Store interface {

	// GetUser  gets a user by unique username
	GetUser(string) (*model.User, error)

	// GetDataList gets a list of all data in the database
	GetDataList() (*[]modle.Data, error)

	// GetData gets a list of data by unique deviceID
	GetData(string) (*[]modle.Data, error)
}
