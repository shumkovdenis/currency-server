package load

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Truefx struct {
	url string
}

func NewTruefx(url string) *Truefx {
	return &Truefx{url}
}

func (s Truefx) GetRate() (float64, error) {
	resp, err := http.Get(s.url)
	if err != nil {
		log.Error(err)
		return 0, ErrRateFailed
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0, ErrRateFailed
	}

	_, line, err := bufio.ScanLines(body, true)
	if err != nil {
		log.Error(err)
		return 0, ErrRateFailed
	}

	arr := strings.Split(string(line), ",")
	value, err := strconv.ParseFloat(arr[2]+arr[3], 64)
	if err != nil {
		log.Error(err)
		return 0, ErrRateFailed
	}

	return value, nil
}
