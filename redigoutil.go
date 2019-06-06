package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

/*====================辅助方法====================*/

func Bool(val interface{}, err error) (bool, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return false, false, nil
		}
		return false, false, err
	}
	switch val := val.(type) {
	case nil:
		return false, false, nil
	case redigo.Error:
		return false, false, val
	case int64:
		ret := val != 0
		return ret, true, nil
	case []byte:
		ret, err := strconv.ParseBool(string(val))
		return ret, true, err
	case string:
		ret, err := strconv.ParseBool(val)
		return ret, true, err
	}
	return false, true, fmt.Errorf("unexpected type for bool: %T", val)
}

func Int(val interface{}, err error) (int, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return 0, false, nil
		}
		return 0, false, err
	}
	switch val := val.(type) {
	case nil:
		return 0, false, nil
	case redigo.Error:
		return 0, false, val
	case int64:
		ret := int(val)
		return ret, true, nil
	case []byte:
		ret, err := strconv.Atoi(string(val))
		return ret, true, err
	case string:
		ret, err := strconv.Atoi(val)
		return ret, true, err
	}
	return 0, true, fmt.Errorf("unexpected type for int: %T", val)
}

func Int64(val interface{}, err error) (int64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return 0, false, nil
		}
		return 0, false, err
	}
	switch val := val.(type) {
	case nil:
		return 0, false, nil
	case redigo.Error:
		return 0, false, val
	case int64:
		return val, true, nil
	case []byte:
		ret, err := strconv.ParseInt(string(val), 10, 64)
		return ret, true, err
	case string:
		ret, err := strconv.ParseInt(val, 10, 64)
		return ret, true, err
	}
	return 0, true, fmt.Errorf("unexpected type for int64: %T", val)
}

func Float64(val interface{}, err error) (float64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return 0, false, nil
		}
		return 0, false, err
	}
	switch val := val.(type) {
	case nil:
		return 0, false, nil
	case redigo.Error:
		return 0, false, val
	case []byte:
		ret, err := strconv.ParseFloat(string(val), 64)
		return ret, true, err
	case int64:
		ret := float64(val)
		return ret, true, nil
	case string:
		ret, err := strconv.ParseFloat(val, 64)
		return ret, true, err
	}
	return 0, false, fmt.Errorf("unexpected type for float64: %T", val)
}

func Bytes(val interface{}, err error) ([]byte, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch val := val.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, val
	case []byte:
		return val, true, nil
	case int64:
		return []byte(strconv.FormatInt(val, 10)), true, nil
	case string:
		return []byte(val), true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for []byte: %T", val)
}

func String(val interface{}, err error) (string, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return "", false, nil
		}
		return "", false, err
	}
	switch val := val.(type) {
	case nil:
		return "", false, nil
	case redigo.Error:
		return "", false, val
	case []byte:
		ret := string(val)
		return ret, true, nil
	case int64:
		ret := strconv.FormatInt(val, 10)
		return ret, true, nil
	case string:
		return val, true, nil
	}
	return "", true, fmt.Errorf("unexpected type for string: %T", val)
}

func IntSlice(reply interface{}, err error) ([]int, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		ret := make([]int, len(reply))
		for i, vi := range reply {
			v, _, err := Int(vi, nil)
			if err != nil {
				return nil, true, err
			} else {
				ret[i] = v
			}
		}
		return ret, true, nil
	default:
		v, _, err := Int(reply, nil)
		if err != nil {
			return nil, true, err
		} else {
			return []int{v}, true, nil
		}
	}

}

func Int64Slice(reply interface{}, err error) ([]int64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		ret := make([]int64, len(reply))
		for i, vi := range reply {
			v, _, err := Int64(vi, nil)
			if err != nil {
				return nil, true, err
			} else {
				ret[i] = v
			}
		}
		return ret, true, nil
	default:
		v, _, err := Int64(reply, nil)
		if err != nil {
			return nil, true, err
		} else {
			return []int64{v}, true, nil
		}
	}
}

func Float64Slice(reply interface{}, err error) ([]float64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		ret := make([]float64, len(reply))
		for i, vi := range reply {
			v, _, err := Float64(vi, nil)
			if err != nil {
				return nil, true, err
			} else {
				ret[i] = v
			}
		}
		return ret, true, nil
	default:
		v, _, err := Float64(reply, nil)
		if err != nil {
			return nil, true, err
		} else {
			return []float64{v}, true, nil
		}
	}
}

func StringSlice(reply interface{}, err error) ([]string, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		ret := make([]string, len(reply))
		for i, vi := range reply {
			v, _, err := String(vi, nil)
			if err != nil {
				return nil, true, err
			} else {
				ret[i] = v
			}
		}
		return ret, true, nil
	default:
		v, _, err := String(reply, nil)
		if err != nil {
			return nil, true, err
		} else {
			return []string{v}, true, nil
		}
	}
}

func BytesSlice(reply interface{}, err error) ([][]byte, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		ret := make([][]byte, len(reply))
		for i, vi := range reply {
			v, _, err := Bytes(vi, nil)
			if err != nil {
				return nil, true, err
			}
			ret[i] = v
		}
		return ret, true, nil
	default:
		v, _, err := Bytes(reply, nil)
		if err != nil {
			return nil, true, err
		}
		return [][]byte{v}, true, nil
	}
}

func ValueSlice(reply interface{}, err error) ([]interface{}, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}
	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		return reply, true, nil
	default:
		return []interface{}{reply}, true, nil
	}
}

func IntMap(reply interface{}, err error) (map[string]int, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make(map[string]int, len/2)
		for i := 1; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			v, _, err := Int(reply[i], nil)
			if err != nil {
				return nil, true, err
			}
			ret[k] = v
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for map[string]int: %T", reply)
}

func Int64Map(reply interface{}, err error) (map[string]int64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make(map[string]int64, len/2)
		for i := 1; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			v, _, err := Int64(reply[i], nil)
			if err != nil {
				return nil, true, err
			}
			ret[k] = v
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for map[string]int64: %T", reply)
}
func Float64Map(reply interface{}, err error) (map[string]float64, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make(map[string]float64, len/2)
		for i := 1; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			v, _, err := Float64(reply[i], nil)
			if err != nil {
				return nil, true, err
			}
			ret[k] = v
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for map[string]float64: %T", reply)
}
func StringMap(reply interface{}, err error) (map[string]string, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make(map[string]string, len/2)
		for i := 1; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			v, _, err := String(reply[i], nil)
			if err != nil {
				return nil, true, err
			}
			ret[k] = v
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for map[string]string: %T", reply)
}

func ValueMap(reply interface{}, err error) (map[string]interface{}, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make(map[string]interface{}, len/2)
		for i := 1; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			ret[k] = reply[i]
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for map[string]interface{}: %T", reply)
}

func ValueScoreSlice(reply interface{}, err error) ([]*ValueScore, bool, error) {
	if err != nil {
		if err == redigo.ErrNil {
			return nil, false, nil
		}
		return nil, false, err
	}

	switch reply := reply.(type) {
	case nil:
		return nil, false, nil
	case redigo.Error:
		return nil, false, reply
	case []interface{}:
		len := len(reply)
		ret := make([]*ValueScore, len/2)
		for i, j := 1, 0; i < len; i += 2 {
			k, _, err := String(reply[i-1], nil)
			if err != nil {
				return nil, true, err
			}
			v, _, err := Float64(reply[i], nil)
			if err != nil {
				return nil, true, err
			}
			ret[j] = &ValueScore{
				Value: k,
				Score: v,
			}
			j++
		}
		return ret, true, nil
	}
	return nil, true, fmt.Errorf("unexpected type for []*ValueScore: %T", reply)
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
