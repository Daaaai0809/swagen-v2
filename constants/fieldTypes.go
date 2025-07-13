package constants

const (
	STRING_TYPE  = "string"
	NUMBER_TYPE  = "number"
	INTEGER_TYPE = "integer"
	BOOLEAN_TYPE = "boolean"
	ARRAY_TYPE   = "array"
	OBJECT_TYPE  = "object"
)

var FieldTypeList = []string{
	STRING_TYPE,
	NUMBER_TYPE,
	INTEGER_TYPE,
	BOOLEAN_TYPE,
	ARRAY_TYPE,
	OBJECT_TYPE,
}

func IsExamplableType(fieldType string) bool {
	switch fieldType {
	case STRING_TYPE, NUMBER_TYPE, INTEGER_TYPE, BOOLEAN_TYPE:
		return true
	default:
		return false
	}
}
