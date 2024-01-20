package user

import "strconv"

// FieldName - user type for storing field name
type FieldName string

const (
	name       FieldName = "name"
	surname    FieldName = "surname"
	patronymic FieldName = "patronymic"
	age        FieldName = "age"
	gender     FieldName = "gender"
	nation     FieldName = "nation"
)

// NewFieldName - constructor with checking correction of field name value
func NewFieldName(val string) (FieldName, error) {
	fieldName := FieldName(val)

	switch fieldName {
	case name, surname, patronymic, age, gender, nation:
		return fieldName, nil
	default:
		return "", newUnknownFieldNameError(fieldName)
	}
}

// ConvertQueryParam converts fieldVal to type by fieldName
// Example: fieldName age - int type; fieldName gender - Gender type; fieldName name - string type
func (fieldName FieldName) ConvertQueryParam(fieldVal string) (any, error) {
	switch fieldName {
	case name, surname, patronymic, nation:
		return fieldVal, nil
	case age:
		val, err := strconv.Atoi(fieldVal)
		if err != nil {
			return nil, err
		}

		return val, nil
	case gender:
		val := Gender(fieldVal)

		err := val.Check()
		if err != nil {
			return nil, err
		}

		return val, nil
	default:
		return nil, newUnknownFieldNameError(fieldName)
	}
}
