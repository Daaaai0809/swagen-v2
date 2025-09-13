package constants

const (
	FORMAT_NONE = "None"
)

const (
	FORMAT_DATE                  = "date"
	FORMAT_DATE_TIME             = "date-time"
	FORMAT_PASSWORD              = "password"
	FORMAT_BYTE                  = "byte"
	FORMAT_BINARY                = "binary"
	FORMAT_EMAIL                 = "email"
	FORMAT_UUID                  = "uuid"
	FORMAT_HOSTNAME              = "hostname"
	FORMAT_IPV4                  = "ipv4"
	FORMAT_IPV6                  = "ipv6"
	FORMAT_URI                   = "uri"
	FORMAT_URI_REF               = "uri-reference"
	FORMAT_URI_TEMPLATE          = "uri-template"
	FORMAT_JSON_POINTER          = "json-pointer"
	FORMAT_RELATIVE_JSON_POINTER = "relative-json-pointer"
	FORMAT_REGEX                 = "regex"
)

var FormatStringList = []string{
	FORMAT_NONE,
	FORMAT_DATE,
	FORMAT_DATE_TIME,
	FORMAT_PASSWORD,
	FORMAT_BYTE,
	FORMAT_BINARY,
	FORMAT_EMAIL,
	FORMAT_UUID,
	FORMAT_HOSTNAME,
	FORMAT_IPV4,
	FORMAT_IPV6,
	FORMAT_URI,
	FORMAT_URI_REF,
	FORMAT_URI_TEMPLATE,
	FORMAT_JSON_POINTER,
	FORMAT_RELATIVE_JSON_POINTER,
	FORMAT_REGEX,
}

const (
	FORMAT_FLOAT  = "float"
	FORMAT_DOUBLE = "double"
)

var FormatNumberList = []string{
	FORMAT_NONE,
	FORMAT_FLOAT,
	FORMAT_DOUBLE,
}

const (
	FORMAT_INT32 = "int32"
	FORMAT_INT64 = "int64"
)

var FormatIntegerList = []string{
	FORMAT_NONE,
	FORMAT_INT32,
	FORMAT_INT64,
}

var FormatList = map[string][]string{
	STRING_TYPE:  FormatStringList,
	NUMBER_TYPE:  FormatNumberList,
	INTEGER_TYPE: FormatIntegerList,
}
