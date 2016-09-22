package load

import "errors"

type Service interface {
	GetRate() (float64, error)
}

var (
	ErrRateFailed = errors.New("Rate failed")
)
