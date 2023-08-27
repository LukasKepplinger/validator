package validator

import (
	"fmt"
	"time"
)

type TimeValid struct {
	*Valid
	value time.Time
	key   string
}

func (v *Valid) Time(value time.Time, key string) *TimeValid {
	if v.cancelled {
		return &TimeValid{Valid: v}
	}
	return &TimeValid{
		Valid: v,
		key:   key,
		value: value,
	}
}

func (v *TimeValid) Past() *TimeValid {
	if v.value.After(time.Now()) {
		v.addError(v.key, "must be in the past")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *TimeValid) Future() *TimeValid {
	if v.value.Before(time.Now()) {
		v.addError(v.key, "must be in the future")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *TimeValid) Before(t time.Time) *TimeValid {
	if v.value.After(t) {
		v.addError(v.key, "must be before "+t.String())
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *TimeValid) After(t time.Time) *TimeValid {
	if v.value.Before(t) {
		v.addError(v.key, "must be after "+t.String())
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}

func (v *TimeValid) Just(seconds int) *TimeValid {

	min := time.Now().Add(time.Duration(seconds) * time.Second)
	max := time.Now().Add(time.Duration(seconds) * time.Second)

	if v.value.Before(min) || v.value.After(max) {
		v.addError(v.key, "must be just "+fmt.Sprint(seconds)+" seconds from now")
		if v.lazy {
			v.cancelled = true
		}
	}

	return v
}

func (v *TimeValid) Required() *TimeValid {
	if v.value.IsZero() {
		v.addError(v.key, "is required")
		if v.lazy {
			v.cancelled = true
		}
	}
	return v
}
