package store

import (
	"errors"
	"time"

	"github.com/easylifewell/purifier-server/model"
)

func register(guest model.Guest) (model.User, error) {
	lastUser := GetLastUser()
	var id int64
	id = 1
	if lastUser.Phone != "" {
		id = lastUser.SID + 1
		lastUser.LastUser = false
		if _, err := UpdateUser(*lastUser); err != nil {
			return *lastUser, errors.New("注册用户失败")
		}
	}
	user := new(model.User)
	user.Phone = guest.Phone
	user.SMSCode = guest.SMSCode
	user.SMSChangeDate = guest.SMSChangeDate
	user.SMSSendDate = guest.SMSSendDate
	user.NickName = GetNickName()
	user.RealName = ""
	user.CreateDate = time.Now()
	user.LoginDate = time.Now()
	user.Email = ""
	user.SID = id
	user.LastUser = true

	if _, err := AddUser(*user); err != nil {
		return *user, errors.New("注册用户失败")
	}

	return *user, nil
}

func RegisterWithSMS(phone, smscode string) (model.User, error) {
	unused := new(model.User)
	guest, err := checkSMSForGuest(phone, smscode)
	if err != nil {
		return *unused, err
	}

	//  向数据库中添加User
	return register(guest)

}

func CheckSMSCode(sid, smscode string) (model.User, error) {
	user := GetUserBySID(sid)
	now := time.Now()

	if now.Sub(user.SMSChangeDate) >= 0 {
		return *user, errors.New("短信验证码已过期")
	}

	if user.SMSCode != smscode {
		return *user, errors.New("短信验证码不正确")
	}
	return *user, nil
}

func checkSMSForGuest(phone, smscode string) (model.Guest, error) {
	guest := GetGuestByPhone(phone)
	now := time.Now()

	if now.Sub(guest.SMSChangeDate) >= 0 {
		return *guest, errors.New("短信验证码已过期")
	}

	if guest.SMSCode != smscode {
		return *guest, errors.New("短信验证码不正确")
	}
	return *guest, nil
}
