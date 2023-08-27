package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

type StringValid struct {
	*Valid
	key   string
	value string
}

var compiledAlphaRegex = regexp.MustCompile("^[a-zA-Z ]*$")
var compiledAlphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9 ]*$")
var compiledNumericRegex = regexp.MustCompile("^[0-9]*$")

func (v *Valid) Text(value, key string) *StringValid {
	if v.cancelled {
		return &StringValid{Valid: v}
	}
	return &StringValid{
		Valid: v,
		key:   key,
		value: value,
	}
}

func (v *StringValid) Email() *StringValid {
	if v.cancelled {
		return v
	}
	if !isEmail(v.value) {
		v.addError(v.key, "is not a valid email")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (v *StringValid) Blacklist(invalid string) *StringValid {
	if v.cancelled {
		return v
	}
	for _, char := range v.value {
		if contains(invalid, string(char)) {
			v.addError(v.key, "contains invalid characters")
			if v.lazy {
				v.cancelled = true
			}
			break
		}
	}
	return v
}

func (v *StringValid) Alpha() *StringValid {
	if v.cancelled {
		return v
	}
	if !compiledAlphaRegex.MatchString(v.value) {
		v.addError(v.key, "must contain only letters")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *StringValid) AlphaNumeric() *StringValid {
	if v.cancelled {
		return v
	}
	if !compiledAlphaNumericRegex.MatchString(v.value) {
		v.addError(v.key, "must contain only letters and numbers")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *StringValid) Numeric() *StringValid {
	if v.cancelled {
		return v
	}
	if !compiledNumericRegex.MatchString(v.value) {
		v.addError(v.key, "must contain only numbers")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *StringValid) Whitelist(valid string) *StringValid {
	if v.cancelled {
		return v
	}
	for _, char := range v.value {
		if !contains(valid, string(char)) {
			v.addError(v.key, "contains invalid characters")
			if v.lazy {
				v.cancelled = true
			}
			break
		}
	}
	return v
}

func contains(valid string, char string) bool {
	for _, validChar := range valid {
		if string(validChar) == char {
			return true
		}
	}
	return false
}

func (v *StringValid) Required() *StringValid {
	if v.cancelled {
		return v
	}
	if v.value == "" {
		v.addError(v.key, "is required")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *StringValid) Min(min int) *StringValid {
	if v.cancelled {
		return v
	}
	if len(v.value) < min {
		v.addError(v.key, "must have at least "+fmt.Sprint(min)+" characters")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *StringValid) Max(max int) *StringValid {
	if v.cancelled {
		return v
	}
	if len(v.value) > max {
		v.addError(v.key, "must have at most "+fmt.Sprint(max)+" characters")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}
