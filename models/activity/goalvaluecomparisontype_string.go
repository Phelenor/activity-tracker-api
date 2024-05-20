// Code generated by "enumer -type=GoalValueComparisonType -json --output goalvaluecomparisontype_string.go -trimprefix GoalValueComparisonType -transform snake-upper"; DO NOT EDIT.

package activity

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _GoalValueComparisonTypeName = "UNDEFINEDGREATERLESS"

var _GoalValueComparisonTypeIndex = [...]uint8{0, 9, 16, 20}

const _GoalValueComparisonTypeLowerName = "undefinedgreaterless"

func (i GoalValueComparisonType) String() string {
	if i < 0 || i >= GoalValueComparisonType(len(_GoalValueComparisonTypeIndex)-1) {
		return fmt.Sprintf("GoalValueComparisonType(%d)", i)
	}
	return _GoalValueComparisonTypeName[_GoalValueComparisonTypeIndex[i]:_GoalValueComparisonTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _GoalValueComparisonTypeNoOp() {
	var x [1]struct{}
	_ = x[GoalValueComparisonTypeUndefined-(0)]
	_ = x[GoalValueComparisonTypeGreater-(1)]
	_ = x[GoalValueComparisonTypeLess-(2)]
}

var _GoalValueComparisonTypeValues = []GoalValueComparisonType{GoalValueComparisonTypeUndefined, GoalValueComparisonTypeGreater, GoalValueComparisonTypeLess}

var _GoalValueComparisonTypeNameToValueMap = map[string]GoalValueComparisonType{
	_GoalValueComparisonTypeName[0:9]:        GoalValueComparisonTypeUndefined,
	_GoalValueComparisonTypeLowerName[0:9]:   GoalValueComparisonTypeUndefined,
	_GoalValueComparisonTypeName[9:16]:       GoalValueComparisonTypeGreater,
	_GoalValueComparisonTypeLowerName[9:16]:  GoalValueComparisonTypeGreater,
	_GoalValueComparisonTypeName[16:20]:      GoalValueComparisonTypeLess,
	_GoalValueComparisonTypeLowerName[16:20]: GoalValueComparisonTypeLess,
}

var _GoalValueComparisonTypeNames = []string{
	_GoalValueComparisonTypeName[0:9],
	_GoalValueComparisonTypeName[9:16],
	_GoalValueComparisonTypeName[16:20],
}

// GoalValueComparisonTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func GoalValueComparisonTypeString(s string) (GoalValueComparisonType, error) {
	if val, ok := _GoalValueComparisonTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _GoalValueComparisonTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to GoalValueComparisonType values", s)
}

// GoalValueComparisonTypeValues returns all values of the enum
func GoalValueComparisonTypeValues() []GoalValueComparisonType {
	return _GoalValueComparisonTypeValues
}

// GoalValueComparisonTypeStrings returns a slice of all String values of the enum
func GoalValueComparisonTypeStrings() []string {
	strs := make([]string, len(_GoalValueComparisonTypeNames))
	copy(strs, _GoalValueComparisonTypeNames)
	return strs
}

// IsAGoalValueComparisonType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i GoalValueComparisonType) IsAGoalValueComparisonType() bool {
	for _, v := range _GoalValueComparisonTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for GoalValueComparisonType
func (i GoalValueComparisonType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for GoalValueComparisonType
func (i *GoalValueComparisonType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("GoalValueComparisonType should be a string, got %s", data)
	}

	var err error
	*i, err = GoalValueComparisonTypeString(s)
	return err
}