// Copyright (c) 2018 Clearmatics Technologies Ltd
package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/abiosoft/ishell"
)

func ConvertAndAppend(c *ishell.Context, input string, argument *abi.Argument, output []interface{}) ([]interface{}, error) {
	if argument.Type.Kind == reflect.Array || argument.Type.Kind == reflect.Slice {
		c.Println("Argument is array\n")

		// One dimensional byte array
		// Accepts all byte arrays as hex string with pre-pended '0x' only
		if argument.Type.Elem == nil {
			if argument.Type.Type == reflect.TypeOf(common.Address{}) {
				// address solidity type
				item, err := ConvertToType(input, &argument.Type)
				if err != nil {
					return nil, err
				}
				output = append(output, item)
				return output, nil
			} else if argument.Type.Type == reflect.TypeOf([]byte{}) {
				// bytes solidity type
				bytes, err := hex.DecodeString(input[2:])
				if err != nil {
					return nil, err
				}
				output = append(output, bytes)
				return output, nil
			} else {
				// Fixed byte array of size n; bytesn solidity type
				// Any submitted bytes longer than the expected size will be truncated

				bytes, err := hex.DecodeString(input[2:])
				if err != nil {
					return nil, err
				}

				// Fixed sized arrays can't be created with variables as size
				switch argument.Type.Size {
				case 1:
					var byteArray [1]byte
					copy(byteArray[:], bytes[:1])
					output = append(output, byteArray)
				case 2:
					var byteArray [2]byte
					copy(byteArray[:], bytes[:2])
					output = append(output, byteArray)
				case 3:
					var byteArray [3]byte
					copy(byteArray[:], bytes[:3])
					output = append(output, byteArray)
				case 4:
					var byteArray [4]byte
					copy(byteArray[:], bytes[:4])
					output = append(output, byteArray)
				case 5:
					var byteArray [5]byte
					copy(byteArray[:], bytes[:5])
					output = append(output, byteArray)
				case 6:
					var byteArray [6]byte
					copy(byteArray[:], bytes[:6])
					output = append(output, byteArray)
				case 7:
					var byteArray [7]byte
					copy(byteArray[:], bytes[:7])
					output = append(output, byteArray)
				case 8:
					var byteArray [8]byte
					copy(byteArray[:], bytes[:8])
					output = append(output, byteArray)
				case 9:
					var byteArray [9]byte
					copy(byteArray[:], bytes[:9])
					output = append(output, byteArray)
				case 10:
					var byteArray [10]byte
					copy(byteArray[:], bytes[:10])
					output = append(output, byteArray)
				case 11:
					var byteArray [11]byte
					copy(byteArray[:], bytes[:11])
					output = append(output, byteArray)
				case 12:
					var byteArray [12]byte
					copy(byteArray[:], bytes[:12])
					output = append(output, byteArray)
				case 13:
					var byteArray [13]byte
					copy(byteArray[:], bytes[:13])
					output = append(output, byteArray)
				case 14:
					var byteArray [14]byte
					copy(byteArray[:], bytes[:14])
					output = append(output, byteArray)
				case 15:
					var byteArray [15]byte
					copy(byteArray[:], bytes[:15])
					output = append(output, byteArray)
				case 16:
					var byteArray [16]byte
					copy(byteArray[:], bytes[:16])
					output = append(output, byteArray)
				case 17:
					var byteArray [17]byte
					copy(byteArray[:], bytes[:17])
					output = append(output, byteArray)
				case 18:
					var byteArray [18]byte
					copy(byteArray[:], bytes[:18])
					output = append(output, byteArray)
				case 19:
					var byteArray [19]byte
					copy(byteArray[:], bytes[:19])
					output = append(output, byteArray)
				case 20:
					var byteArray [20]byte
					copy(byteArray[:], bytes[:20])
					output = append(output, byteArray)
				case 21:
					var byteArray [21]byte
					copy(byteArray[:], bytes[:21])
					output = append(output, byteArray)
				case 22:
					var byteArray [22]byte
					copy(byteArray[:], bytes[:22])
					output = append(output, byteArray)
				case 23:
					var byteArray [23]byte
					copy(byteArray[:], bytes[:23])
					output = append(output, byteArray)
				case 24:
					var byteArray [24]byte
					copy(byteArray[:], bytes[:24])
					output = append(output, byteArray)
				case 25:
					var byteArray [25]byte
					copy(byteArray[:], bytes[:25])
					output = append(output, byteArray)
				case 26:
					var byteArray [26]byte
					copy(byteArray[:], bytes[:26])
					output = append(output, byteArray)
				case 27:
					var byteArray [27]byte
					copy(byteArray[:], bytes[:27])
					output = append(output, byteArray)
				case 28:
					var byteArray [28]byte
					copy(byteArray[:], bytes[:28])
					output = append(output, byteArray)
				case 29:
					var byteArray [29]byte
					copy(byteArray[:], bytes[:29])
					output = append(output, byteArray)
				case 30:
					var byteArray [30]byte
					copy(byteArray[:], bytes[:30])
					output = append(output, byteArray)
				case 31:
					var byteArray [31]byte
					copy(byteArray[:], bytes[:31])
					output = append(output, byteArray)
				case 32:
					var byteArray [32]byte
					copy(byteArray[:], bytes[:32])
					output = append(output, byteArray)
				default:
					errStr := fmt.Sprintf("Error parsing fixed size byte array. Array of size %d incompatible", argument.Type.Size)
					return nil, errors.New(errStr)
				}
				return output, nil
			}

		}

		array := strings.Split(input, ",")
		argSize := argument.Type.Size
		size := len(array)
		if argSize != 0 {
			for size != argSize {
				c.Printf("Please enter %i comma-separated list of elements:\n", argSize)
				input = c.ReadLine()
				array = strings.Split(input, ",")
				size = len(array)
			}
		}

		size = len(array)

		elementType := argument.Type.Elem

		// Elements cannot be kind slice                                        only mean slice
		if elementType.Kind == reflect.Array {
			// Is 2D byte array
			/* Nightmare to implement, have to account for:
			   * Slice of fixed byte arrays; bytes32[] in solidity for example, generally bytesn[]
			   * Fixed array of fixed byte arrays; bytes32[10] in solidity for example bytesn[m]

			   Since the upper bound of elements in an array in solidity is 2^256-1, and each fixed byte array
			   has a limit of bytes32 (bytes1, bytes2, ..., bytes31, bytes32), and Golang array creation takes
			   constant length values, we would have to paste the switch-case containing 1-32 fixed byte arrays
			   2^256-1 times to handle every possibility. Since arrays of arrays in seldom used, we have not
			   implemented it.
			*/

			return nil, errors.New("2D Arrays unsupported. Use \"bytes\" instead.")

			/*
			   slice := make([]interface{}, 0, size)
			   err = addFixedByteArrays(array, elementType.Size, slice)
			   if err != nil {
			       return nil, err
			   }
			   output = append(output, slice)
			   continue
			*/
		} else {
			switch elementType.Type {
			case reflect.TypeOf(bool(false)):
				convertedArray := make([]bool, 0, size)
				for _, item := range array {
					b, err := ConvertToBool(item)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, b)
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(int8(0)):
				convertedArray := make([]int8, 0, size)
				for _, item := range array {
					i, err := strconv.ParseInt(item, 10, 8)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, int8(i))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(int16(0)):
				convertedArray := make([]int16, 0, size)
				for _, item := range array {
					i, err := strconv.ParseInt(item, 10, 16)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, int16(i))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(int32(0)):
				convertedArray := make([]int32, 0, size)
				for _, item := range array {
					i, err := strconv.ParseInt(item, 10, 32)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, int32(i))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(int64(0)):
				convertedArray := make([]int64, 0, size)
				for _, item := range array {
					i, err := strconv.ParseInt(item, 10, 64)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, int64(i))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(uint8(0)):
				convertedArray := make([]uint8, 0, size)
				for _, item := range array {
					u, err := strconv.ParseUint(item, 10, 8)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, uint8(u))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(uint16(0)):
				convertedArray := make([]uint16, 0, size)
				for _, item := range array {
					u, err := strconv.ParseUint(item, 10, 16)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, uint16(u))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(uint32(0)):
				convertedArray := make([]uint32, 0, size)
				for _, item := range array {
					u, err := strconv.ParseUint(item, 10, 32)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, uint32(u))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(uint64(0)):
				convertedArray := make([]uint64, 0, size)
				for _, item := range array {
					u, err := strconv.ParseUint(item, 10, 64)
					if err != nil {
						return nil, err
					}
					convertedArray = append(convertedArray, uint64(u))
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(&big.Int{}):
				convertedArray := make([]*big.Int, 0, size)
				for _, item := range array {
					newInt := new(big.Int)
					newInt, ok := newInt.SetString(item, 10)
					if !ok {
						return nil, errors.New("Could not convert string to big.int")
					}
					convertedArray = append(convertedArray, newInt)
				}
				output = append(output, convertedArray)
			case reflect.TypeOf(common.Address{}):
				convertedArray := make([]common.Address, 0, size)
				for _, item := range array {
					a := common.HexToAddress(item)
					convertedArray = append(convertedArray, a)
				}
				output = append(output, convertedArray)
			default:
				errStr := fmt.Sprintf("Type %s not found", elementType.Type)
				return nil, errors.New(errStr)
			}
		}
	} else {
		switch argument.Type.Kind {
		case reflect.String:
			output = append(output, input)
		case reflect.Bool:
			b, err := ConvertToBool(input)
			if err != nil {
				return nil, err
			}
			output = append(output, b)
		case reflect.Int8:
			i, err := strconv.ParseInt(input, 10, 8)
			if err != nil {
				return nil, err
			}
			output = append(output, int8(i))
		case reflect.Int16:
			i, err := strconv.ParseInt(input, 10, 16)
			if err != nil {
				return nil, err
			}
			output = append(output, int16(i))
		case reflect.Int32:
			i, err := strconv.ParseInt(input, 10, 32)
			if err != nil {
				return nil, err
			}
			output = append(output, int32(i))
		case reflect.Int64:
			i, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return nil, err
			}
			output = append(output, int64(i))
		case reflect.Uint8:
			u, err := strconv.ParseUint(input, 10, 8)
			if err != nil {
				return nil, err
			}
			output = append(output, uint8(u))
		case reflect.Uint16:
			u, err := strconv.ParseUint(input, 10, 16)
			if err != nil {
				return nil, err
			}
			output = append(output, uint16(u))
		case reflect.Uint32:
			u, err := strconv.ParseUint(input, 10, 32)
			if err != nil {
				return nil, err
			}
			output = append(output, uint32(u))
		case reflect.Uint64:
			u, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				return nil, err
			}
			output = append(output, uint64(u))
		case reflect.Ptr:
			newInt := new(big.Int)
			newInt, ok := newInt.SetString(input, 10)
			if !ok {
				return nil, errors.New("Could not convert string to big.int")
			}
			output = append(output, newInt)
		case reflect.Array:
			if argument.Type.Type == reflect.TypeOf(common.Address{}) {
				address := common.HexToAddress(input)
				output = append(output, address)
			} else {
				return nil, errors.New("Conversion failed. Item is array type, cannot parse")
			}
		default:
			errStr := fmt.Sprintf("Error, type not found: %s", argument.Type.Kind)
			return nil, errors.New(errStr)
		}
	}
	return output, nil
}

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
	if size%8 > 0 {
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
