package controller

import (
	"net/http"

	"github.com/easylifewell/purifier-server/store"
	"github.com/gorilla/mux"
)

type DeviceController struct{}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (uc DeviceController) BindDeviceWithUser(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
		return
	}

	user := store.GetUserByPhone(phone)
	if user.Phone == "" {
		Response200(w, "用户不存在，请您先注册系统")
		return
	}

	vars := mux.Vars(r)
	nickname := vars["deviceid"]
	if nickname == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}
	// 添加用户绑定的操作

	return
}

func (uc DeviceController) SetCarName(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
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
	Response200(w, "更新CarName成功")
	return
}

func (uc DeviceController) SetSPhone(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
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

	device := store.GetDeviceByID(deviceid)
	device.SPhone = sphone
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "更新sphone失败")
		return
	}
	Response200(w, "更新sphone成功")
	return
}

func (uc DeviceController) SetRPhone(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
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

	device := store.GetDeviceByID(deviceid)
	device.RPhone = rphone
	if _, err := store.UpdateDevice(*device); err != nil {
		Response500(w, "更新rphone失败")
		return
	}
	Response200(w, "更新rphone成功")
	return
}

func (uc DeviceController) On(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
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
		Response500(w, "更新statred失败")
		return
	}
	Response200(w, "更新started成功")
	return
}

func (uc DeviceController) Off(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
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
		Response500(w, "更新statred失败")
		return
	}
	Response200(w, "更新started成功")
	return
}

func (uc DeviceController) IsStarted(w http.ResponseWriter, r *http.Request) {
	_, ok := isLogin(r)
	if !ok {
		Response200(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	deviceid := vars["deviceid"]
	if deviceid == "" {
		Response400(w, "请输入请求参数deviceid")
		return
	}

	//device := store.GetDeviceByID(deviceid)
	// 有待进一步完善，返回什么字段给用户呢
	Response200(w, "更新started成功")
	return
}
