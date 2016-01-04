package server

import (
	"fmt"
	"github.com/quipo/statsd"
	"golang.org/x/net/context"
	"time"
)

type DBEvent int

const (
	DBStart DBEvent = iota
	DBStop
)

type Instrumentation struct {
	MetricID           string
	startTime          time.Time
	marshalStartTime   time.Time
	marshalDuration    time.Duration
	unmarshalStartTime time.Time
	unmarshalDuration  time.Duration
}

var StatsdClient *statsd.StatsdClient

func Next(ctx context.Context, s string) bool {
	instrumentation, ok := ctx.Value("instrumentation").(*Instrumentation)
	if !ok {
		return false
	}

	instrumentation.Next(s)

	return true
}

func NewInstrumentation(metricID string) *Instrumentation {
	return &Instrumentation{MetricID: metricID}
}

func (i *Instrumentation) StartHandler() {
	i.startTime = time.Now()
}

func (i *Instrumentation) Next(name string) {
	i.MetricID += "." + name
}

func (i *Instrumentation) StartMarshal() {
	i.marshalStartTime = time.Now()
}

func (i *Instrumentation) StopMarshal() {
	i.StopMarshalWithDuration(time.Duration(time.Now().UnixNano() - i.marshalStartTime.UnixNano()))
}

func (i *Instrumentation) StopMarshalWithDuration(d time.Duration) {
	i.marshalDuration += d
}

func (i *Instrumentation) StartUnmarshal() {
	i.unmarshalStartTime = time.Now()
}

func (i *Instrumentation) StopUnmarshal() {
	i.StopUnmarshalWithDuration(time.Duration(time.Now().UnixNano() - i.unmarshalStartTime.UnixNano()))
}

func (i *Instrumentation) StopUnmarshalWithDuration(d time.Duration) {
	i.unmarshalDuration += d
}

func shouldInstrument() bool {
	return StatsdClient != nil
}

func (i *Instrumentation) StopHandler(status int) {
	if !shouldInstrument() {
		return
	}

	d := time.Duration(time.Now().UnixNano() - i.startTime.UnixNano())

	StatsdClient.Timing(i.MetricID+".unmarshal_time", int64(i.unmarshalDuration))
	StatsdClient.Timing(i.MetricID+".marshal_time", int64(i.marshalDuration))
	StatsdClient.Timing(i.MetricID+".total_time", int64(d))
	StatsdClient.Incr(fmt.Sprintf("%s.status.%v", i.MetricID, status), 1)
}
