package common

// ErrorCode
type ErrorCode int32

const (
	_ int32 = iota + 9999
	// StatusOK : 10000 OK
	StatusOK
	// StatusParamInvalid :  10001
	StatusParamInvalid
	// StatusServerError : 10002
	StatusServerError
	// StatusRegisterFailed : 10003
	StatusRegisterFailed
	// StatusLoginFailed : 10004
	StatusLoginFailed
	// StatusTokenInvalid : 10005
	StatusTokenInvalid
	// StatusUserNotExists: 10006
	StatusUserNotExists
)
