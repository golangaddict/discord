package discord

import (
	"errors"
	"fmt"
	en2 "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vdtr "github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strconv"
	"strings"
)

var (
	validate   *vdtr.Validate
	translator ut.Translator
)

func init() {
	en := en2.New()
	uni := ut.New(en, en)

	translator, _ = uni.GetTranslator("en")
	validate = vdtr.New()
	en_trans.RegisterDefaultTranslations(validate, translator)
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		v, ok := field.Tag.Lookup("name")
		if ok {
			return v
		}

		return field.Name
	})
}

func Bind(args []string, out any) error {
	if err := bind(args, out); err != nil {
		return err
	}

	if err := validate.Struct(out); err != nil {
		if v, ok := err.(vdtr.ValidationErrors); ok {
			return errors.New(v[0].Translate(translator))
		}

		return err
	}

	return nil
}

// TODO: remove index slicing with : since we use csv parser now.
func bind(args []string, out any) error {
	t := reflect.TypeOf(out).Elem()
	v := reflect.ValueOf(out).Elem()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		if tagValue, ok := fieldType.Tag.Lookup("index"); ok {
			tagData := strings.Split(tagValue, ":")
			startIndex, err := strconv.Atoi(tagData[0])
			if err != nil {
				return err
			}

			fieldValue := v.Field(i)

			if len(args) < startIndex+1 {
				fieldValue.Set(reflect.Zero(fieldType.Type))
				continue
			}

			if len(tagData) > 1 {
				endIndex, _ := strconv.Atoi(tagData[1])

				if endIndex == 0 {
					fieldValue.Set(reflect.ValueOf(strings.Join(args[startIndex:], " ")))
					continue
				}

				fieldValue.Set(reflect.ValueOf(strings.Join(args[startIndex:endIndex], " ")))
				continue
			}

			arg, err := parseArg(args[startIndex], fieldValue.Kind())
			if err != nil {
				return fmt.Errorf("validate: failed to parse argument %s: %w", fieldType.Tag.Get("name"), err)
			}
			fieldValue.Set(reflect.ValueOf(arg).Convert(fieldValue.Type()))
		}
	}

	return nil
}

func parseArg(arg string, t reflect.Kind) (any, error) {
	if t == reflect.String {
		return arg, nil
	}

	switch t {
	case reflect.Int:
		i, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}

		return i, nil
	}

	return nil, fmt.Errorf("unknown field type: %s", t)
}
