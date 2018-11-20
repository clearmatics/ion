// Copyright (c) 2018 Clearmatics Technologies Ltd
package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func ConvertToType(str string, typ *abi.Type) (interface{}, error) {
    switch typ.Kind {
    case reflect.String:
        return str, nil
    case reflect.Bool:
        b, err := ConvertToBool(str)
        return b, err
    case reflect.Int8:
        i, err := ConvertToInt(true, 8, str)
        return i, err
    case reflect.Int16:
        i, err := ConvertToInt(true, 16, str)
        return i, err
    case reflect.Int32:
        i, err := ConvertToInt(true, 32, str)
        return i, err
    case reflect.Int64:
        i, err := ConvertToInt(true, 64, str)
        return i, err
    case reflect.Uint8:
        u, err := ConvertToInt(false, 8, str)
        return u, err
    case reflect.Uint16:
        u, err := ConvertToInt(false, 16, str)
        return u, err
    case reflect.Uint32:
        u, err := ConvertToInt(false, 32, str)
        return u, err
    case reflect.Uint64:
        u, err := ConvertToInt(false, 64, str)
        return u, err
    case reflect.Ptr:
        i, err := ConvertToInt(false, typ.Size, str)
        return i, err
    case reflect.Array:
        if typ.Type == reflect.TypeOf(common.Address{}) {
            return common.HexToAddress(str), nil
        } else {
            return nil, errors.New("Conversion failed. Item is array type, cannot parse")
        }
    default:
        errStr := fmt.Sprintf("Error, type not found: %s", typ.Kind)
        return nil, errors.New(errStr)
    }
}

func ConvertToInt(signed bool, size int, value string) (interface{}, error) {
    if size % 8 > 0 {
        return nil, errors.New("Integer is not a multiple of 8")
    } else if !isGoIntSize(size) {
        newInt := new(big.Int)
        newInt, ok := newInt.SetString(value, 10)
        if !ok {
            return nil, errors.New("Could not convert string to big.int")
        }

        return newInt, nil
    } else {
        if signed {
            i, err := strconv.ParseInt(value, 10, size)
            if err != nil {
                return nil, err
            }
            return i, nil
        } else {
            u, err := strconv.ParseUint(value, 10, size)
            if err != nil {
                return nil, err
            }
            return u, nil
        }
    }
}

// MUST CHECK RETURNED ERROR ELSE WILL RETURN FALSE FOR ANY ERRONEOUS INPUT
func ConvertToBool(value string) (bool, error) {
    b, err := strconv.ParseBool(value)
    if err != nil {
        return false, err
    }
    return b, nil
}

func isGoIntSize(size int) (isGoPrimitive bool) {
    switch size {
    case 8, 16, 32, 64:
        return true
    default:
        return false
    }
}