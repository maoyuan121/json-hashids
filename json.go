package jsonhashids

import (
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
)

// 定义一个错误。指出输入不是整型
var ErrNotInteger = errors.New("not integer")

// 为 jsoniter 注册扩展
func NewConfigWithHashIDs(salt string, minLength int) jsoniter.API {

	e := NewHashIDsExtension(salt, minLength)
	config := jsoniter.ConfigCompatibleWithStandardLibrary
	config.RegisterExtension(e)

	return config
}

// 实例化扩展
// @param salt 加密盐
// @param minLength 生成的 id 的最小长度
// return 返回扩展
func NewHashIDsExtension(salt string, minLength int) *HashIDsExtension {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	return &HashIDsExtension{
		hashid: h,
	}
}

// 扩展结构
type HashIDsExtension struct {
	hashid *hashids.HashID
	jsoniter.DummyExtension
}

// 更新结构解析器
// 实现扩展接口
func (extension *HashIDsExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		// 如果字段有 tag hashids 且值为 true 继续， 否则执行下一次迭代
		tag := binding.Field.Tag().Get("hashids")
		if tag != "true" {
			continue
		}

		// 如果字段类型为整型则继续，否则执行下一次迭代
		switch binding.Field.Type().Kind() {
		case reflect.Int:
		case reflect.Uint:
		case reflect.Int8:
		case reflect.Uint8:
		case reflect.Int16:
		case reflect.Uint16:
		case reflect.Int32:
		case reflect.Uint32:
		case reflect.Int64:
		case reflect.Uint64:
		default:
			continue
		}

		typeName := binding.Field.Type().String()
		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {

			i, err := int64Val(typeName, ptr)
			if err != nil {
				stream.Error = err
				return
			}

			hashed, err := extension.hashid.EncodeInt64([]int64{i})
			if err != nil {
				stream.Error = err
				return
			}
			stream.WriteString(hashed)
		}}

		binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {

			str := iter.ReadString()
			if str == "" {
				return
			}

			ints := extension.hashid.DecodeInt64(str)
			if len(ints) != 1 {
				return
			}

			setIntValue(typeName, ptr, ints[0])
		}}
	}
}

// 实现 ValDecoder 接口
type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

// 实现 ValEncoder 接口
type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}

// 转换为 int64
// @param typeName 字段的类型
// @param ptr 指向字段的指针
func int64Val(typeName string, ptr unsafe.Pointer) (int64, error) {
	var i int64

	switch typeName {
	case "int":
		ip := (*int)(ptr)
		i = int64(*ip)
	case "uint":
		ip := (*uint)(ptr)
		i = int64(*ip)
	case "int8":
		ip := (*int8)(ptr)
		i = int64(*ip)
	case "uint8":
		ip := (*uint8)(ptr)
		i = int64(*ip)
	case "int16":
		ip := (*int16)(ptr)
		i = int64(*ip)
	case "uint16":
		ip := (*uint16)(ptr)
		i = int64(*ip)
	case "int32":
		ip := (*int32)(ptr)
		i = int64(*ip)
	case "uint32":
		ip := (*uint32)(ptr)
		i = int64(*ip)
	case "int64":
		ip := (*int64)(ptr)
		i = *ip
	case "uint64":
		ip := (*uint64)(ptr)
		i = int64(*ip)
	default:
		return 0, ErrNotInteger
	}

	return i, nil
}

// 设置为整型
// @param typename 字段类型
// @param ptr 指向字段的指针
// @param val 要设置的值
func setIntValue(typeName string, ptr unsafe.Pointer, val int64) error {
	switch typeName {
	case "int":
		ip := (*int)(ptr)
		*ip = int(val)
	case "uint":
		ip := (*uint)(ptr)
		*ip = uint(val)
	case "int8":
		ip := (*int8)(ptr)
		*ip = int8(val)
	case "uint8":
		ip := (*uint8)(ptr)
		*ip = uint8(val)
	case "int16":
		ip := (*int16)(ptr)
		*ip = int16(val)
	case "uint16":
		ip := (*uint16)(ptr)
		*ip = uint16(val)
	case "int32":
		ip := (*int32)(ptr)
		*ip = int32(val)
	case "uint32":
		ip := (*uint32)(ptr)
		*ip = uint32(val)
	case "int64":
		ip := (*int64)(ptr)
		*ip = int64(val)
	case "uint64":
		ip := (*uint64)(ptr)
		*ip = uint64(val)
	default:
		return ErrNotInteger
	}

	return nil
}
