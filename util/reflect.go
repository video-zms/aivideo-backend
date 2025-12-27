package util

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func GetFieldJsonTag(f reflect.StructField) (string, string) {
	jt := f.Tag.Get("json")
	jk := jt
	jo := ""
	//tag逗号分隔的第一个为name，后面的为option
	if idx := strings.Index(jk, ","); idx != -1 {
		jk = jt[:idx]
		jo = jt[idx+1:]
	}

	if jk == "" || jk == "-" {
		jk = f.Name
	}
	return jk, jo
}

func GetStructTag(v interface{}) []string {
	var keys []string
	keysSlice := keys[:]
	ty := reflect.TypeOf(v)
	for i := 0; i < ty.NumField(); i++ {
		key := ty.Field(i).Tag.Get("db")
		if len(key) == 0 || key == "-" {
			continue
		}
		keysSlice = append(keysSlice, key)
	}
	return keysSlice
}

func GetStructTagDeep(v interface{}) []string {
	var keys []string
	keysSlice := keys[:]
	ty := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		ty = ty.Elem()
	}
	for i := 0; i < ty.NumField(); i++ {
		if value.Field(i).Kind() == reflect.Struct {
			//当使用嵌套结构体时小心非公开（小写）嵌套结构体会导致挂起，并且不会panic

			//如果结构体为非公开，直接跳过
			if !ast.IsExported(value.Field(i).Type().Name()) {
				continue
			}
			keysSlice = append(keysSlice, GetStructTagDeep(value.Field(i).Interface())...)
		}
		key := ty.Field(i).Tag.Get("db")
		if len(key) == 0 || key == "-" {
			continue
		}
		keysSlice = append(keysSlice, key)
	}
	return keysSlice
}

// use this carefully, it is designed for interface{} with no specific struct definition
func GetMapFieldsDeep(data map[string]interface{}) []string {
	var keys []string
	keysSlice := keys[:]
	for k, v := range data {
		if reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
			keysSlice = append(keysSlice, GetMapFieldsDeep(v.(map[string]interface{}))...)
			continue
		}
		keysSlice = append(keysSlice, k)
	}
	return keysSlice
}

func GetStructFieldsAndValuesDeep(s interface{}) ([]string, []string) {
	return getFieldsAndValues(GetStructTagDeep(s))
}

func getMapFieldsAndValuesDeep(s map[string]interface{}) ([]string, []string) {
	return getFieldsAndValues(GetMapFieldsDeep(s))
}

// s is a struct
func GetStructFieldsAndValues(s interface{}) ([]string, []string) {
	return getFieldsAndValues(GetStructTag(s))
}

func getFieldsAndValues(rawTags []string) ([]string, []string) {
	var fields []string
	var values []string
	for _, v := range rawTags {
		fields = append(fields, "`"+v+"`")
		values = append(values, ":"+v)
	}
	return fields, values
}

func GetMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

func GetMethodAndLine() string {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	file, line := f.FileLine(pc)
	return file + " line: " + strconv.Itoa(line) + " " + f.Name()
}

func PackErr(err error, prefix string) error {
	if err == nil {
		return nil
	}

	return errors.New(prefix + " \r\n " + err.Error())
}

func GetSimpleStackTrace() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return strings.Replace(f.Name(), "wespy-http-go/app/", "", -1)
}

func GetFullStackTrace() string {
	bs := debug.Stack()
	return string(bs)
}

/*
*
根据db配置获取结构体key
*/
func GetStructTagOnlyNeedDb(v interface{}) []string {
	var keys []string
	keysSlice := keys[:]
	ty := reflect.TypeOf(v)
	for i := 0; i < ty.NumField(); i++ {
		key := ty.Field(i).Tag.Get("db")
		if len(key) == 0 || key == "-" {
			continue
		}
		keysSlice = append(keysSlice, key)
	}
	return keysSlice
}

// 使用gob序列化实现的深拷贝，已知会把空数组拷成null
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}

// 使用json序列化实现的深拷贝
func DeepCopyJson(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return json.NewDecoder(&buf).Decode(dst)
}

func CopyContext(c *gin.Context) *gin.Context {
	if c == nil {
		return nil
	}
	return c.Copy()
}

func GetStructFieldsAndValuesOnlyNeedDb(s interface{}) ([]string, []string) {
	var fields []string
	var values []string
	tags := GetStructTagOnlyNeedDb(s)
	for _, v := range tags {
		fields = append(fields, v)
		values = append(values, ":"+v)
	}
	return fields, values
}

// 将嵌套结构体扁平化，并将所有带tag标记的字段以string形式存入map
func StructToStringMap(s interface{}, tag string, dest map[string]string) {
	//ret := make(map[string]string, 0)
	ty := reflect.TypeOf(s)
	value := reflect.ValueOf(s)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		ty = ty.Elem()
	}

	if value.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < value.NumField(); i++ {
		field := ty.Field(i)
		v := value.Field(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			StructToStringMap(v.Interface(), tag, dest)
			continue
		}

		if key := field.Tag.Get(tag); key != "" {
			dest[key] = fmt.Sprint(v.Interface())
		}
	}
}

func Struct2Map(s interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, 0)
	ty := reflect.TypeOf(s)
	value := reflect.ValueOf(s)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		ty = ty.Elem()
	}

	if value.Kind() != reflect.Struct {
		return ret
	}

	for i := 0; i < value.NumField(); i++ {
		field := ty.Field(i)
		v := value.Field(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			ret[field.Tag.Get("json")] = Struct2Map(v.Interface())
			continue
		}

		ret[field.Tag.Get("json")] = v.Interface()
	}
	return ret
}

func InterfaceMapToStringMap(m map[string]interface{}, dest map[string]string) {
	for k, v := range m {
		str, _ := json.Marshal(v)
		dest[k] = string(str)
	}
}

// 一般用于将嵌套的json信息重新序列化，存到第一层map中
func SuppressStringMap(dataMap map[string]interface{}) {
	for k := range dataMap {
		v := reflect.ValueOf(dataMap[k])
		if v.Kind() == reflect.Map {
			bytes, _ := json.Marshal(dataMap[k])
			dataMap[k] = string(bytes)
		}
	}
}

var ErrTargetFieldNoMatch = errors.New("target fields not matching root struct")
var ErrNilTarget = errors.New("empty target fields")
var ErrNotPtr = errors.New("param must be a pointer")
var ErrNil = errors.New("root struct cannot be nil")
var ErrNotStruct = errors.New("root not of struct type")

func GetDBTagName(rootStruct interface{}, targetField interface{}) (string, error) {
	tags, err := getTagNames("db", rootStruct, targetField)
	if err != nil {
		return "", err
	}
	return tags[0], nil
}

func GetDBTagNames(rootStruct interface{}, targetField ...interface{}) ([]string, error) {
	return getTagNames("db", rootStruct, targetField...)
}

func getTagNames(tag string, rootStruct interface{}, targetFields ...interface{}) ([]string, error) {
	if len(targetFields) == 0 {
		return nil, ErrNilTarget
	}
	fvs := []reflect.Value{}
	rv := reflect.ValueOf(rootStruct)
	if rv.Kind() != reflect.Ptr {
		return nil, ErrNotPtr
	}
	for _, targetField := range targetFields {
		fv := reflect.ValueOf(targetField)
		if fv.Kind() != reflect.Ptr {
			return nil, ErrNotPtr
		}
		fvs = append(fvs, fv)
	}

	rv, ok := deRefV(rv)
	if !ok {
		return nil, ErrNil
	}
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	tags := []string{}

	for i := 0; i < rv.NumField(); i++ {
		rfv := rv.Field(i)
		if rfv.CanAddr() {
			for _, fv := range fvs {
				if rfv.Addr().Pointer() == fv.Pointer() {
					tagV := rt.Field(i).Tag.Get(tag)
					tags = append(tags, tagV)
				}
			}
		}
	}
	if len(tags) != len(targetFields) {
		return nil, ErrTargetFieldNoMatch
	}
	return tags, nil
}

func getTagName(rootStruct interface{}, targetField interface{}, tag string) (string, error) {
	fv := reflect.ValueOf(targetField)
	rv := reflect.ValueOf(rootStruct)
	if fv.Kind() != reflect.Ptr || rv.Kind() != reflect.Ptr {
		return "", ErrNotPtr
	}

	rv, ok := deRefV(rv)
	if !ok {
		return "", ErrNil
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		rfv := rv.Field(i)
		if rfv.CanAddr() {
			if rfv.Addr().Pointer() == fv.Pointer() {
				tagV := rt.Field(i).Tag.Get(tag)
				return tagV, nil
			}
		}
	}
	return "", nil
}

func deRefV(v reflect.Value) (reflect.Value, bool) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v, false
		}
		v = v.Elem()
	}
	return v, true
}

func GetInt(v interface{}) int64 {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float())
	}
	return 0
}

// GetStructTagAndValueDeep 顺序是绝对的 https://groups.google.com/g/golang-nuts/c/_he7g1TL0K8/m/8BBCFB0tCgAJ
func GetStructTagAndValueDeep(v interface{}) ([]string, []interface{}) {
	if v == nil {
		return []string{}, []interface{}{}
	}
	keys := make([]string, 0)
	datas := make([]interface{}, 0)
	getStructTagAndValueDeep(v, &keys, &datas)
	return keys, datas
}

func getStructTagAndValueDeep(v interface{}, keys *[]string, datas *[]interface{}) {
	ty := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	for i := 0; i < 10; i++ {
		if value.Kind() != reflect.Ptr {
			break
		}
		value = value.Elem()
		ty = ty.Elem()
	}
	for i := 0; i < ty.NumField(); i++ {
		switch value.Field(i).Kind() {
		case reflect.Struct:
			//当使用嵌套结构体时小心非公开（小写）嵌套结构体会导致挂起，并且不会panic
			//如果结构体为非公开，直接跳过
			if !ast.IsExported(value.Field(i).Type().Name()) {
				continue
			}
			getStructTagAndValueDeep(value.Field(i).Interface(), keys, datas)
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
			reflect.String:
			key := ty.Field(i).Tag.Get("db")
			if len(key) <= 0 || key == "-" {
				continue
			}
			*keys = append(*keys, key)
			*datas = append(*datas, value.Field(i).Interface())
		}
	}
}
