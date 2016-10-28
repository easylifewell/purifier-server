package store

import (
	"math/rand"
	"time"

	"github.com/easylifewell/purifier-server/model"
)

func GetNickName() string {
	length := len(model.Animal)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return model.Animal[r1.Intn(length)]
}
