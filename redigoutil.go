package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

/*====================辅助方法====================*/

func Bool(val interface{}, err error) (*bool, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case int64:
		ret := val != 0
		return &ret, nil
	case []byte:
		ret, err := strconv.ParseBool(string(val))
		return &ret, err
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		ret, err := strconv.ParseBool(val)
		return &ret, err
	}
	return nil, fmt.Errorf("unexpected type for bool: %T", val)
}

func Int(val interface{}, err error) (*int, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case int64:
		ret := int(val)
		return &ret, nil
	case []byte:
		ret, err := strconv.Atoi(string(val))
		return &ret, err
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		ret, err := strconv.Atoi(val)
		return &ret, err
	}
	return nil, fmt.Errorf("unexpected type for int: %T", val)
}

func Int64(val interface{}, err error) (*int64, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case int64:
		return &val, nil
	case []byte:
		ret, err := strconv.ParseInt(string(val), 10, 64)
		return &ret, err
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		ret, err := strconv.ParseInt(val, 10, 64)
		return &ret, err
	}
	return nil, fmt.Errorf("unexpected type for int64: %T", val)
}

func Float64(val interface{}, err error) (*float64, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case []byte:
		ret, err := strconv.ParseFloat(string(val), 64)
		return &ret, err
	case int64:
		ret := float64(val)
		return &ret, nil
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		ret, err := strconv.ParseFloat(val, 64)
		return &ret, err
	}
	return nil, fmt.Errorf("unexpected type for float64: %T", val)
}

func Bytes(val interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case []byte:
		return val, nil
	case int64:
		return []byte(strconv.FormatInt(val, 10)), nil
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		return []byte(val), nil
	}
	return nil, fmt.Errorf("unexpected type for []byte: %T", val)
}

func String(val interface{}, err error) (*string, error) {
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case []byte:
		ret := string(val)
		return &ret, nil
	case int64:
		ret := strconv.FormatInt(val, 10)
		return &ret, nil
	case nil:
		return nil, nil
	case redigo.Error:
		return nil, val
	case string:
		return &val, nil
	}
	return nil, fmt.Errorf("unexpected type for string: %T", val)
}

func IntSlice(reply interface{}, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		ret := make([]int, len(reply))
		for i, vi := range reply {
			v, err := Int(vi, nil)
			if err != nil {
				return nil, err
			} else if v != nil {
				ret[i] = *v
			} else {
				ret[i] = 0
			}
		}
		return ret, nil
	case nil:
		return nil, nil
	default:
		v, err := Int(reply, nil)
		if err != nil {
			return nil, err
		} else if v != nil {
			return []int{*v}, nil
		} else {
			return []int{0}, nil
		}
	}

}

func Int64Slice(reply interface{}, err error) ([]int64, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		ret := make([]int64, len(reply))
		for i, vi := range reply {
			v, err := Int64(vi, nil)
			if err != nil {
				return nil, err
			} else if v != nil {
				ret[i] = *v
			} else {
				ret[i] = 0
			}
		}
		return ret, nil
	case nil:
		return nil, nil
	default:
		v, err := Int64(reply, nil)
		if err != nil {
			return nil, err
		} else if v != nil {
			return []int64{*v}, nil
		} else {
			return []int64{0}, nil
		}
	}
}

func Float64Slice(reply interface{}, err error) ([]float64, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		ret := make([]float64, len(reply))
		for i, vi := range reply {
			v, err := Float64(vi, nil)
			if err != nil {
				return nil, err
			} else if v != nil {
				ret[i] = *v
			} else {
				ret[i] = 0
			}
		}
		return ret, nil
	case nil:
		return nil, nil
	default:
		v, err := Float64(reply, nil)
		if err != nil {
			return nil, err
		} else if v != nil {
			return []float64{*v}, nil
		} else {
			return []float64{0}, nil
		}
	}
}

func StringSlice(reply interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		ret := make([]string, len(reply))
		for i, vi := range reply {
			v, err := String(vi, nil)
			if err != nil {
				return nil, err
			} else if v != nil {
				ret[i] = *v
			} else {
				ret[i] = ""
			}
		}
		return ret, nil
	case nil:
		return nil, nil
	default:
		v, err := String(reply, nil)
		if err != nil {
			return nil, err
		} else if v != nil {
			return []string{*v}, nil
		} else {
			return []string{""}, nil
		}
	}
}

func BytesSlice(reply interface{}, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		ret := make([][]byte, len(reply))
		for i, vi := range reply {
			v, err := Bytes(vi, nil)
			if err != nil {
				return nil, err
			}
			ret[i] = v
		}
		return ret, nil
	case nil:
		return nil, nil
	default:
		v, err := Bytes(reply, nil)
		if err != nil {
			return nil, err
		}
		return [][]byte{v}, nil
	}
}

func ValueSlice(reply interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		return reply, nil
	case nil:
		return nil, nil
	default:
		return []interface{}{reply}, nil
	}
}

func IntMap(reply interface{}, err error) (map[string]int, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make(map[string]int, len/2)
		for i := 1; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			v, err := Int(reply[i], nil)
			if err != nil {
				return nil, err
			}
			ret[*k] = *v
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for map[string]int: %T", reply)
}

func Int64Map(reply interface{}, err error) (map[string]int64, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make(map[string]int64, len/2)
		for i := 1; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			v, err := Int64(reply[i], nil)
			if err != nil {
				return nil, err
			}
			ret[*k] = *v
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for map[string]int64: %T", reply)
}
func Float64Map(reply interface{}, err error) (map[string]float64, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make(map[string]float64, len/2)
		for i := 1; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			v, err := Float64(reply[i], nil)
			if err != nil {
				return nil, err
			}
			ret[*k] = *v
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for map[string]float64: %T", reply)
}
func StringMap(reply interface{}, err error) (map[string]string, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make(map[string]string, len/2)
		for i := 1; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			v, err := String(reply[i], nil)
			if err != nil {
				return nil, err
			}
			ret[*k] = *v
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for map[string]string: %T", reply)
}

func ValueMap(reply interface{}, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make(map[string]interface{}, len/2)
		for i := 1; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			ret[*k] = reply[i]
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for map[string]interface{}: %T", reply)
}

func ValueScoreSlice(reply interface{}, err error) ([]*ValueScore, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []interface{}:
		len := len(reply)
		ret := make([]*ValueScore, len/2)
		for i, j := 1, 0; i < len; i += 2 {
			k, err := String(reply[i-1], nil)
			if err != nil {
				return nil, err
			}
			v, err := Float64(reply[i], nil)
			if err != nil {
				return nil, err
			}
			ret[j] = &ValueScore{
				Value: *k,
				Score: *v,
			}
			j++
		}
		return ret, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected type for []*ValueScore: %T", reply)
}

func Scan(sc Scanner, src interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	return sc.Scan(src)
}

var ScanSlice = redigo.Scan
var ScanStruct = redigo.ScanStruct
var ScanStructSlice = redigo.ScanSlice

// func ScanStructSlice(src []interface{}, dest interface{}, fieldNames ...string) error
