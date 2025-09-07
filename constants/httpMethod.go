package constants

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
