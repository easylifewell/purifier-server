package store

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/easylifewell/purifier-server/model"
)

func register(guest *model.Guest) (*model.User, error) {
	lastUser := GetLastUser()
	var id int64
	id = 1
	if lastUser.Phone != "" {
		id = lastUser.SID + 1
		lastUser.LastUser = false
		if _, err := UpdateUser(*lastUser); err != nil {
			return lastUser, errors.New("注册用户失败")
		}
	}
	user := new(model.User)
	user.Phone = guest.Phone
	user.SMSCode = guest.SMSCode
	user.SMSChangeDate = guest.SMSChangeDate
	user.SMSSendDate = guest.SMSSendDate
	user.NickName = GetNickName()
	user.RealName = ""
	user.Password = guest.Password
	user.CreateDate = time.Now()
	user.LoginDate = time.Now()
	user.Email = ""
	user.SID = id
	user.LastUser = true

	if _, err := AddUser(*user); err != nil {
		return user, errors.New("注册用户失败")
	}

	return user, nil
}

func Register(phone, smscode, password string) (*model.User, error) {
	unused := new(model.User)
	guest, err := checkSMSForGuest(phone, smscode)
	if err != nil {
		return unused, err
	}
	guest.Password = password

	//  向数据库中添加User
	return register(guest)

}

func Login(sid, password string) (*model.User, error) {
	s, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, errors.New("解析用户id错误")
	}
	user := GetUserBySID(s)
	now := time.Now()

	fmt.Printf("user = %v\n", user)
	fmt.Printf("now = %v\n", now)

	if user.Password != password {
		return user, errors.New("密码错误")
	}

	return user, nil
}

func ForgetPassword(sid, smscode, password string) (*model.User, error) {
	s, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, errors.New("解析用户id错误")
	}
	user := GetUserBySID(s)
	now := time.Now()

	fmt.Printf("user = %v\n", user)
	fmt.Printf("now = %v\n", now)
	fmt.Printf("changedate = %v\n", user.SMSChangeDate)
	if now.Sub(user.SMSChangeDate) >= 0 {
		return user, errors.New("短信验证码已过期")
	}

	if user.SMSCode != smscode {
		return user, errors.New("短信验证码不正确")
	}

	user.Password = password
	if _, err := UpdateUser(*user); err != nil {
		return user, errors.New("重置密码失败")
	}

	return user, nil
}

func checkSMSForGuest(phone, smscode string) (*model.Guest, error) {
	guest := GetGuestByPhone(phone)
	now := time.Now()

	if now.Sub(guest.SMSChangeDate) >= 0 {
		return guest, errors.New("短信验证码已过期")
	}

	if guest.SMSCode != smscode {
		return guest, errors.New("短信验证码不正确")
	}
	return guest, nil
}
