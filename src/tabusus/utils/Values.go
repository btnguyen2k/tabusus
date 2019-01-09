package utils

// ToInt32 casts/converts a value to int32
func ToInt32(v interface{}) (int32, bool) {
	switch v := v.(type) {
	case int:
		return int32(v), true
	case uint:
		return int32(v), true
	case int8:
		return int32(v), true
	case uint8:
		return int32(v), true
	case int16:
		return int32(v), true
	case uint16:
		return int32(v), true
	case int32:
		return int32(v), true
	case uint32:
		return int32(v), true
	case int64:
		return int32(v), true
	case uint64:
		return int32(v), true
	}
	return 0, false
}

// ToString casts/converts a value to string
func ToString(v interface{}) (string, bool) {
	switch v := v.(type) {
	case string:
		return string(v), true
	}
	return "", false
}
