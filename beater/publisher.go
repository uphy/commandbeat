package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type (
	// Publisher is a interface responsible for publishing events to anywhere.
	Publisher interface {
		Publish(spec *commandSpec, v common.MapStr)
	}
	// elasticPublisher is a publisher publishes events to elasticsearch.
	elasticPublisher struct {
		client beat.Client
	}
)

func newElasticPublisher(client beat.Client) *elasticPublisher {
	return &elasticPublisher{client}
}

func (e *elasticPublisher) Publish(spec *commandSpec, v common.MapStr) {
	var timestamp time.Time
	if t, ok := v["@timestamp"]; ok {
		timestamp = t.(time.Time)
		delete(v, "@timestamp")
	} else {
		timestamp = time.Now()
	}
	v["type"] = spec.name
	event := beat.Event{
		Timestamp: timestamp,
		Fields:    v,
	}
	e.LogDebug(spec, "<event>%#v", event)
	if !spec.debug {
		e.client.Publish(event)
	}
}

func (e *elasticPublisher) LogDebug(spec *commandSpec, msg string, args ...interface{}) {
	if spec.debug {
		logp.Info("[%s] %s", spec.name, fmt.Sprintf(msg, args...))
	}
}
