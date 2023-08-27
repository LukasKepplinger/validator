package validator

import (
	"fmt"
)

type NumberValid struct {
	*Valid
	key   string
	value int
}

func (v *Valid) Number(value int, key string) *NumberValid {
	if v.cancelled {
		return &NumberValid{Valid: v}
	}
	return &NumberValid{
		Valid: v,
		key:   key,
		value: value,
	}
}

func (v *NumberValid) Required() *NumberValid {
	if v.value == 0 {
		v.addError(v.key, "is required")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *NumberValid) Min(min int) *NumberValid {
	if v.value < min {
		v.addError(v.key, "must be at least "+fmt.Sprint(min))
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *NumberValid) Max(max int) *NumberValid {
	if v.value > max {
		v.addError(v.key, "must be at most "+fmt.Sprint(max))
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}
