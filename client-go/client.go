package clientgo

import (
	"net/http"
)

const (
	HttpScheme = "http://"
)

var client *http.Client

func init() {
	// TODO : make singleton
	client = &http.Client{}
}
