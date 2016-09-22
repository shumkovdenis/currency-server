package rate

import "errors"

type Service interface {
	Start()
	LastRate() (Rate, error)
	Rates() (Rates, error)
	Error() error
}

type Rate struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	err       error
}

type Rates []Rate

var (
	ErrNoRates = errors.New("No rates")
)

func (r Rates) Len() int {
	return len(r)
}

func (r Rates) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rates) Less(i, j int) bool {
	return r[i].Timestamp < r[j].Timestamp
}
