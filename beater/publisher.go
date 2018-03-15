package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/uphy/commandbeat/command"
)

type (
	Publisher interface {
		Publish(spec *command.Spec, v common.MapStr)
	}
	elasticPublisher struct {
		client beat.Client
	}
)

func newElasticsearchPublisher(client beat.Client) *elasticPublisher {
	return &elasticPublisher{client}
}

func (e *elasticPublisher) Publish(spec *command.Spec, v common.MapStr) {
	var timestamp time.Time
	if t, ok := v["@timestamp"]; ok {
		timestamp = t.(time.Time)
		delete(v, "@timestamp")
	} else {
		timestamp = time.Now()
	}
	v["type"] = spec.Name
	event := beat.Event{
		Timestamp: timestamp,
		Fields:    v,
	}
	e.LogDebug(spec, "<event>%#v", event)
	if !spec.Debug {
		e.client.Publish(event)
	}
}

func (e *elasticPublisher) LogDebug(spec *command.Spec, msg string, args ...interface{}) {
	if spec.Debug {
		logp.Info("[%s] %s", spec.Name, fmt.Sprintf(msg, args...))
	}
}
