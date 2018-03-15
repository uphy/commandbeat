package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/uphy/commandbeat/command"
)

type (
	Publisher interface {
		Publish(spec *command.Spec, v common.MapStr)
	}
	elasticsearchPublisher struct {
		client beat.Client
	}
)

func newElasticsearchPublisher(client beat.Client) Publisher {
	return &elasticsearchPublisher{client}
}

func (e *elasticsearchPublisher) Publish(spec *command.Spec, v common.MapStr) {
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
	spec.LogDebug("<event>%#v", event)
	if !spec.Debug {
		e.client.Publish(event)
	}
}
