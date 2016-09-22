package rate

import (
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/shumkovdenis/currency-server/load"
)

type Single struct {
	length  int
	timeout time.Duration
	service load.Service
	rates   Rates
	getTime time.Time
	err     error
}

func NewSingle(length int, timeout time.Duration, service load.Service) *Single {
	return &Single{
		length:  length,
		timeout: timeout,
		service: service,
		rates:   make(Rates, 0, length),
	}
}

func (s *Single) Start() {
	go func() {
		tick := time.Tick(1 * time.Second)
		rates := s.next(tick)
		for rate := range rates {
			s.addRate(rate)
		}
	}()
}

func (s *Single) LastRate() (Rate, error) {
	l := len(s.rates)
	if l == 0 {
		return Rate{}, ErrNoRates
	}
	return s.rates[l-1], nil
}

func (s *Single) Rates() (Rates, error) {
	l := len(s.rates)
	if l == 0 {
		return nil, ErrNoRates
	}
	return s.rates, nil
}

func (s *Single) Error() error {
	return s.err
}

func (s *Single) addRate(rate Rate) {
	if len(s.rates) == s.length {
		s.rates = s.rates[1:]
	}
	s.rates = append(s.rates, rate)
	sort.Sort(s.rates)
	s.err = rate.err
}

func (s *Single) next(in <-chan time.Time) <-chan Rate {
	out := make(chan Rate)
	go func() {
		for now := range in {
			go s.get(out, now)
		}
	}()
	return out
}

func (s *Single) get(in chan<- Rate, t time.Time) {
	value, getErr := s.service.GetRate()
	isCopy := false
	if getErr != nil {
		log.Error(getErr)
		last, err := s.LastRate()
		if err == nil && t.Sub(s.getTime) < s.timeout {
			getErr = nil
			isCopy = true
			value = last.Value
		}
	} else {
		s.getTime = t
	}

	msg := " rate"
	if isCopy {
		msg = "Copy" + msg
	} else {
		msg = "New" + msg
	}
	if getErr == nil {
		log.WithField("value", value).Debug(msg)
	}

	in <- Rate{t.Unix() * 1000, value, getErr}
}
