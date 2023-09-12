package e

var MsgFlags = map[int]string{
	Success:                    "success",
	Error:                      "fail",
	InvalidParams:              "参数错误",
	ErrorExistUser:             "该用户名已存在",
	ErrorDatabase:              "数据库错误",
	ErrorFailEncryption:        "密码加密失败",
	ErrorExistUserNotFound:     "该用户不存在",
	ErrorNotCompare:            "密码错误",
	ErrorAuthToken:             "token 认证失败",
	ErrorAuthCheckTokenTimeout: "token 过期",
	ErrorExistFavorite:         "该商品已经存在于收藏夹",
	ErrorAddressNotExist:       "该收货地址不存在",
	ErrorProductNotExist:       "该商品不存在",

	WebsocketSuccess: "消息发送成功",
	WebsocketError:   "消息发送失败",
}

//GstMag 获取状态码对应信息

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
