package validator

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/fatih/color"
)

type Valid struct {
	lazy        bool
	cancelled   bool
	errors      map[string][]string
	makroErrors map[string]map[string][]string
}

// func New() Valid {
// 	return Valid{errors: make(map[string][]string)}
// }

func (v Valid) Lazy() Valid {
	v.lazy = true
	return v
}

type Validatable interface {
	Validate() Valid
}

func (v *Valid) Struct(validatable Validatable, key string) {
	if v.cancelled {
		return
	}
	validation := validatable.Validate()
	if len(validation.errors) > 0 {
		if v.makroErrors == nil {
			v.makroErrors = make(map[string]map[string][]string)
		}
		if validation.makroErrors[key] != nil {
			key = key + "|"
		}
		v.makroErrors[key] = validation.errors
	}
}

func (v *Valid) addError(key string, message string) {
	if v.errors == nil {
		v.errors = make(map[string][]string)
	}
	if v.errors[key] == nil {
		v.errors[key] = make([]string, 0)
	}
	errorsForKey := v.errors[key]
	errorsForKey = append(errorsForKey, message)
	v.errors[key] = errorsForKey

	if v.lazy {
		v.cancelled = true
	}
}

func (v Valid) Log() Valid {
	if len(v.errors) > 0 || len(v.makroErrors) > 0 {
		redBackground := color.New(color.BgRed).SprintFunc()
		log.Println(redBackground(" VALIDATION ERROR "))
		printErrorsOnField(v.errors)
		for item, makroErrors := range v.makroErrors {
			log.Println(item + ":")
			printErrorsOnFieldIndented(makroErrors)
		}
		log.Println(redBackground(" END VALIDATION ERROR "))
	}
	return v
}

func printErrorsOnField(errors map[string][]string) {
	for field, errorsOnField := range errors {
		log.Println(field + ":")
		for _, err := range errorsOnField {
			log.Println(" - " + err)
		}
	}
}

func printErrorsOnFieldIndented(errors map[string][]string) {
	for field, errorsOnField := range errors {
		log.Println(" - " + field + ":")
		for _, err := range errorsOnField {
			log.Println("    - " + err)
		}
	}
}

func (v Valid) LogStructured() Valid {
	if len(v.errors) > 0 || len(v.makroErrors) > 0 {
		redBackground := color.New(color.BgRed).SprintFunc()
		log.Println(redBackground(" VALIDATION ERROR "))

		errs, err := json.Marshal(v.errors)
		if err != nil {
			log.Println(err)
		}
		makroErrs, err := json.Marshal(v.makroErrors)
		if err != nil {
			log.Println(err)
		}
		errorsJson := append(errs, makroErrs...)

		log.Println(string(errorsJson))
		log.Println(redBackground(" END VALIDATION ERROR "))
	}
	return v
}

func (v Valid) HandleResponse(w http.ResponseWriter) error {
	if len(v.errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		errs, err := json.Marshal(v.errors)
		if err != nil {
			log.Println(err)
		}
		makroErrs, err := json.Marshal(v.makroErrors)
		if err != nil {
			log.Println(err)
		}
		errorsJson := append(errs, makroErrs...)

		if err != nil {
			log.Println(err)
		}
		w.Write(errorsJson)
		return errors.New("validation error")
	}
	return nil
}

func (v Valid) Handle() error {
	if len(v.errors) > 0 {
		valErrors, err := json.Marshal(v.errors)
		if err != nil {
			log.Println(err)
		}
		return errors.New(string(valErrors))
	}
	return nil
}
