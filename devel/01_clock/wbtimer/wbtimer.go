package wbtimer

import (
	"time"

	"github.com/beevik/ntp"
)

const (
	DefaultHost = "0.beevik-ntp.pool.ntp.org"
)

type WbTimer struct {
	response *ntp.Response
	host     string
}

func New(host string) (*WbTimer, error) {
	response, err := ntp.Query(host)
	if err != nil {
		return nil, err
	}
	return &WbTimer{
		response: response,
		host:     host,
	}, nil
}

func (c *WbTimer) CurrentTime() (time.Time, time.Time) {
	prec := time.Now().Add(c.response.ClockOffset)
	timeNow := time.Now()
	return prec, timeNow
}
