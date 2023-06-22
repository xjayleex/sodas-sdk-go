package environ

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/xjayleex/sodas-sdk-go/runtime/property"
	"github.com/xjayleex/sodas-sdk-go/tools"
)

// Retrieve retreives all properties in AppProperty fields from
// runtime environemt variables.
// Usage :
// prop := property.HelloSodasProperties{}
// property.Retrieve(&prop)
// token := prop.Input.RefreshToken
func Retrieve(prop property.AppProperties) error {
	v := reflect.ValueOf(prop).Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("prop interface must be a struct")
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		kind := field.Kind()
		tagValue, ok := v.Type().Field(i).Tag.Lookup(prop.RootFieldTag())
		if !ok {
			return fmt.Errorf("AppProperties tag lookup error : field %s does not have tag %s", field.Type().Name(), prop.RootFieldTag())
		}
		envValue := os.Getenv(tagValue)
		if envValue == "" {
			return fmt.Errorf("runtime does not has os environment for key %s ", tagValue)
		}

		switch kind {
		case reflect.Struct:
			envValue = tools.ReplaceQuotes(envValue)
			if err := json.Unmarshal([]byte(envValue), field.Addr().Interface()); err != nil {
				return err
			}
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			alpha, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(alpha))
		case reflect.Float64:
			f, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		case reflect.Float32:
			f, err := strconv.ParseFloat(envValue, 32)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		case reflect.Bool:
			// Accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
			b, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			field.SetBool(b)
		case reflect.Slice:
			return errors.New("parsing logic for slice type is not implemented yet")

		default:
			return errors.New("environ value has unexpected type")
		}
	}
	return nil
}
