package constants

import "fmt"

const (
	HTTP_GET     = "GET"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"
	HTTP_DELETE  = "DELETE"
	HTTP_PATCH   = "PATCH"
	HTTP_HEAD    = "HEAD"
	HTTP_OPTIONS = "OPTIONS"

	_HTTP_GET     = "get"
	_HTTP_POST    = "post"
	_HTTP_PUT     = "put"
	_HTTP_DELETE  = "delete"
	_HTTP_PATCH   = "patch"
	_HTTP_HEAD    = "head"
	_HTTP_OPTIONS = "options"
)

var HTTPMethods = []string{
	HTTP_GET,
	HTTP_POST,
	HTTP_PUT,
	HTTP_DELETE,
	HTTP_PATCH,
	HTTP_HEAD,
	HTTP_OPTIONS,
}

var HTTPMethodsMap = map[string]string{
	HTTP_GET:     _HTTP_GET,
	HTTP_POST:    _HTTP_POST,
	HTTP_PUT:     _HTTP_PUT,
	HTTP_DELETE:  _HTTP_DELETE,
	HTTP_PATCH:   _HTTP_PATCH,
	HTTP_HEAD:    _HTTP_HEAD,
	HTTP_OPTIONS: _HTTP_OPTIONS,
}

var httpMethodsMap = map[string]string{
	_HTTP_GET:     HTTP_GET,
	_HTTP_POST:    HTTP_POST,
	_HTTP_PUT:     HTTP_PUT,
	_HTTP_DELETE:  HTTP_DELETE,
	_HTTP_PATCH:   HTTP_PATCH,
	_HTTP_HEAD:    HTTP_HEAD,
	_HTTP_OPTIONS: HTTP_OPTIONS,
}

func GetNotExistingMethods(existedMethods []string) []string {
	notExisted := make([]string, 0)
	existedMap := make(map[string]bool)
	for _, method := range existedMethods {
		existedMap[method] = true
	}

	for method := range httpMethodsMap {
		if !existedMap[method] {
			notExisted = append(notExisted, httpMethodsMap[method])
		}
	}

	fmt.Printf("Not existed methods: %v\n", notExisted)

	return notExisted
}
