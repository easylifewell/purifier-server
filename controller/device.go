package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/easylifewell/purifier-server/model"
	"github.com/easylifewell/purifier-server/store"
	"github.com/gorilla/mux"
)

type DeviceStatus struct {
	DeviceID  string `json:"id"`
	IsStarted int    `json:"on"`
}

type DeviceController struct{}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (uc DeviceController) GetDevices(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	user := store.GetUserByPhone(phone)
	if user.Phone == "" {
		Response400(w, "用户不存在，请您先注册系统")
		return
	}

	devices, err := store.GetDeviceByUserID(user.SID)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := io.WriteString(w, "Get devices from database failed")
		if err != nil {
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	logrus.WithFields(logrus.Fields{
		"user":    user.Phone,
		"devices": user.Devices,
	}).Info("查询用户的设备成功")
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		panic(err)
	}
}

func (uc DeviceController) BindDeviceWithUser(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	user := store.GetUserByPhone(phone)
	if user.Phone == "" {
		Response400(w, "用户不存在，请您先注册系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}

	var device model.Device
	device.DeviceID = deviceid
	device.CarName = ""
	device.SPhone = ""
	device.RPhone = ""
	device.Started = false
	device.UserID = user.SID

	_, err := store.AddDevice(device)
	if err != nil {
		Response500(w, "添加设备错误")
		return
	}

	if Contains(user.Devices, deviceid) {
		logrus.WithFields(logrus.Fields{
			"user":     user.Phone,
			"deviceid": deviceid,
		}).Info("设备已经添加,请不要重复添加")
		Response400(w, "设备已经添加,请不要重复添加")
		return
	}
	user.Devices = append(user.Devices, deviceid)
	if _, err := store.UpdateUser(*user); err != nil {
		Response500(w, "更新User设备列表失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"user":     user.Phone,
		"deviceid": deviceid,
	}).Info("添加设备成功")
	Response200(w, "添加设备成功")
	return
}

func (uc DeviceController) SetCarName(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	carname := vars["carname"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}
	if carname == "" {
		Response400(w, "请输入请求参数car")
		return
	}

	device := store.GetDeviceByID(deviceid)
	device.CarName = carname
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "更新CarName失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"deviceid": deviceid,
		"car name": carname,
	}).Info("更新CarName成功")
	Response200(w, "更新CarName成功")
	return
}

func (uc DeviceController) SetSPhone(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	sphone := vars["sphone"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}
	if sphone == "" {
		Response400(w, "请输入请求参数sphone")
		return
	}
	if !Phone.MatchString(sphone) {
		Response400(w, "无效的手机号码")
		return
	}

	device := store.GetDeviceByID(deviceid)
	device.SPhone = sphone
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "更新sphone失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"deviceid": deviceid,
		"sphone":   sphone,
	}).Info("更新sphone成功")
	Response200(w, "更新sphone成功")
	return
}

func (uc DeviceController) SetRPhone(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	rphone := vars["rphone"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}
	if rphone == "" {
		Response400(w, "请输入请求参数rphone")
		return
	}
	if !Phone.MatchString(rphone) {
		Response400(w, "无效的手机号码")
		return
	}

	device := store.GetDeviceByID(deviceid)
	device.RPhone = rphone
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "更新rphone失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"deviceid": deviceid,
		"rphone":   rphone,
	}).Info("更新rphone成功")
	Response200(w, "更新rphone成功")
	return
}

func (uc DeviceController) On(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}

	device := store.GetDeviceByID(deviceid)
	device.Started = true
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "打开净化器失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"deviceid": deviceid,
		"started":  true,
	}).Info("打开净化器成功")
	Response200(w, "打开净化器成功")
	return
}

func (uc DeviceController) Off(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}

	device := store.GetDeviceByID(deviceid)
	device.Started = false
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "关闭净化器失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"deviceid": deviceid,
		"started":  false,
	}).Info("关闭净化器成功")
	Response200(w, "关闭净化器成功")
	return
}

func (uc DeviceController) IsStarted(w http.ResponseWriter, r *http.Request) {
	key1 := "Edb^@u2T"
	key2 := "aZEqm5ph"
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

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}
	if len(deviceid) != 10 {
		Response400(w, "设备ID的长度必须是10")
	}
	t := "VIEW" + key1 + deviceid + key2
	if t != token {
		Response400(w, "无效的TOKEN")
	}

	device := store.GetDeviceByID(deviceid)
	d := new(DeviceStatus)
	d.DeviceID = device.DeviceID
	if device.Started {
		d.IsStarted = 1
	} else {
		d.IsStarted = 0
	}
	b, _ := json.Marshal(d)
	ResponseWithJSON(w, b, http.StatusOK)
	logrus.WithFields(logrus.Fields{
		"DeviceID":  d.DeviceID,
		"IsStarted": d.IsStarted,
	}).Info("查询净化器状态成功")
	return
}
