// Code generated by "stringer -type=endState"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[esUnknown-0]
	_ = x[esSuccess-1]
	_ = x[esPreconditionFailure-2]
	_ = x[esCanaryFailure-3]
	_ = x[esMaxFailures-4]
}

const _endState_name = "esUnknownesSuccessesPreconditionFailureesCanaryFailureesMaxFailures"

var _endState_index = [...]uint8{0, 9, 18, 39, 54, 67}

func (i endState) String() string {
	if i < 0 || i >= endState(len(_endState_index)-1) {
		return "endState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _endState_name[_endState_index[i]:_endState_index[i+1]]
}
