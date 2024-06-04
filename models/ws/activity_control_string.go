// Code generated by "enumer -type=ActivityControl -json -text -output activity_control_string.go -trimprefix ActivityControl -transform snake-upper"; DO NOT EDIT.

package ws

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _ActivityControlName = "UNDEFINEDSTARTPAUSERESUMEFINISH"

var _ActivityControlIndex = [...]uint8{0, 9, 14, 19, 25, 31}

const _ActivityControlLowerName = "undefinedstartpauseresumefinish"

func (i ActivityControl) String() string {
	if i < 0 || i >= ActivityControl(len(_ActivityControlIndex)-1) {
		return fmt.Sprintf("ActivityControl(%d)", i)
	}
	return _ActivityControlName[_ActivityControlIndex[i]:_ActivityControlIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ActivityControlNoOp() {
	var x [1]struct{}
	_ = x[ActivityControlUndefined-(0)]
	_ = x[ActivityControlStart-(1)]
	_ = x[ActivityControlPause-(2)]
	_ = x[ActivityControlResume-(3)]
	_ = x[ActivityControlFinish-(4)]
}

var _ActivityControlValues = []ActivityControl{ActivityControlUndefined, ActivityControlStart, ActivityControlPause, ActivityControlResume, ActivityControlFinish}

var _ActivityControlNameToValueMap = map[string]ActivityControl{
	_ActivityControlName[0:9]:        ActivityControlUndefined,
	_ActivityControlLowerName[0:9]:   ActivityControlUndefined,
	_ActivityControlName[9:14]:       ActivityControlStart,
	_ActivityControlLowerName[9:14]:  ActivityControlStart,
	_ActivityControlName[14:19]:      ActivityControlPause,
	_ActivityControlLowerName[14:19]: ActivityControlPause,
	_ActivityControlName[19:25]:      ActivityControlResume,
	_ActivityControlLowerName[19:25]: ActivityControlResume,
	_ActivityControlName[25:31]:      ActivityControlFinish,
	_ActivityControlLowerName[25:31]: ActivityControlFinish,
}

var _ActivityControlNames = []string{
	_ActivityControlName[0:9],
	_ActivityControlName[9:14],
	_ActivityControlName[14:19],
	_ActivityControlName[19:25],
	_ActivityControlName[25:31],
}

// ActivityControlString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ActivityControlString(s string) (ActivityControl, error) {
	if val, ok := _ActivityControlNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ActivityControlNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ActivityControl values", s)
}

// ActivityControlValues returns all values of the enum
func ActivityControlValues() []ActivityControl {
	return _ActivityControlValues
}

// ActivityControlStrings returns a slice of all String values of the enum
func ActivityControlStrings() []string {
	strs := make([]string, len(_ActivityControlNames))
	copy(strs, _ActivityControlNames)
	return strs
}

// IsAActivityControl returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ActivityControl) IsAActivityControl() bool {
	for _, v := range _ActivityControlValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ActivityControl
func (i ActivityControl) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ActivityControl
func (i *ActivityControl) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ActivityControl should be a string, got %s", data)
	}

	var err error
	*i, err = ActivityControlString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for ActivityControl
func (i ActivityControl) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ActivityControl
func (i *ActivityControl) UnmarshalText(text []byte) error {
	var err error
	*i, err = ActivityControlString(string(text))
	return err
}
