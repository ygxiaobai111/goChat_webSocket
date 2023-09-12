package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400
	// db
	ErrorDatabase              = 1000
	ErrorExistUser             = 30001
	ErrorFailEncryption        = 30002
	ErrorExistUserNotFound     = 30003
	ErrorNotCompare            = 30004
	ErrorAuthToken             = 30005
	ErrorAuthCheckTokenTimeout = 30006
	ErrorUploadFail            = 30007
	ErrorEmail                 = 30008

	//product Model
	ErrorProductImgUpload = 40001
	ErrorProductNotExist  = 40002

	//favorite
	ErrorExistFavorite = 50001

	//address
	ErrorAddressNotExist = 60001

	//socket
	WebsocketSuccess = 70000
	WebsocketError   = 700001
)
