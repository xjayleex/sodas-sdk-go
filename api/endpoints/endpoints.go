package endpoints

var Endpoints = map[string]string{}

// DEVOPS
const (
	API_V1_DEVOPS_DEVELOPMENT_ENVIRONMENT_GET  = "/devops/development/environment/get"
	API_V1_DEVOPS_DEVELOPMENT_ENVIRONMENT_LIST = "/devops/development/environment/list"
)

// GATEWAY
const (
	API_V1_GATEWAY_AUTHENTICATION_USER_LOGIN       = "/gateway/authentication/user/login"
	API_V1_GATEWAY_AUTHENTICATION_USER_REFRESHUSER = "/gateway/authentication/user/refreshUser"
)
