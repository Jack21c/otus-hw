package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	var errs ValidationErrors
	valueV := reflect.ValueOf(v)
	for i := 0; i < valueV.NumField(); i++ {
		st := valueV.Type().Field(i)
		val := valueV.Field(i)
		errs = append(errs, validateEl(st, val)...)
	}
	// Place your code here.
	return nil
}

func validateEl(st reflect.StructField, val reflect.Value) ValidationErrors {
	var errs ValidationErrors
	go func() {
		switch val.Kind() {
		case reflect.String:
			errs = append(errs, validateString(st, val)...)
		case reflect.Int:
			errs = append(errs, validateInt(st, val)...)
		case reflect.Struct:
			for i := 0; i < val.NumField(); i++ {
				st2 := val.Type().Field(i)
				val2 := val.Field(i)
				errs = append(errs, validateEl(st2, val2)...)
			}
		case reflect.Slice:

		}
	}()

	return errs
}

func validateInt(st reflect.StructField, val reflect.Value) ValidationErrors {
	var errs ValidationErrors
	value, ok := st.Tag.Lookup("min")
	if ok {
		min, err := strconv.Atoi(value)
		if err != nil {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   err,
			})
		} else if val.Int() < int64(min) {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   errors.New("число меньше минимального значения"),
			})
		}
	}

	value, ok = st.Tag.Lookup("max")
	if ok {
		max, err := strconv.Atoi(value)
		if err != nil {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   err,
			})
		} else if val.Int() > int64(max) {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   errors.New("число больше максимального значения"),
			})
		}
	}

	value, ok = st.Tag.Lookup("in")
	if ok {
		in := false
		strs := strings.Split(value, ",")
		for _, s := range strs {
			num, err := strconv.Atoi(s)
			if err != nil {
				errs = append(errs, ValidationError{
					Field: st.Name,
					Err:   err,
				})
			} else if val.Equal(reflect.ValueOf(num)) {
				in = true
				break
			}
		}
		if !in {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   errors.New("число не входит в указанное множество чисел"),
			})
		}
	}

	return errs
}

func validateString(st reflect.StructField, val reflect.Value) ValidationErrors {
	var errs ValidationErrors
	value, ok := st.Tag.Lookup("len")
	if ok {
		len, err := strconv.Atoi(value)
		if err != nil {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   err,
			})
		} else if len != val.Len() {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   errors.New("длина строки не совпадает с указанной"),
			})
		}
	}

	value, ok = st.Tag.Lookup("in")
	if ok {
		in := false
		strs := strings.Split(value, ",")
		for _, s := range strs {
			if val.Equal(reflect.ValueOf(s)) {
				in = true
				break
			}
		}
		if !in {
			errs = append(errs, ValidationError{
				Field: st.Name,
				Err:   errors.New("строка не входит в указанное множество строк"),
			})
		}
	}

	return errs
}
