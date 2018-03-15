package beater

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/uphy/commandbeat/command"

	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

type (
	publishHandler struct {
		publisher *elasticPublisher
	}
	Publisher interface {
		Publish(spec *commandSpec, v common.MapStr)
	}
	elasticPublisher struct {
		client beat.Client
	}
)

func newPublishHandler(publisher *elasticPublisher) *publishHandler {
	return &publishHandler{publisher}
}

func (p *publishHandler) BeforeStart(spec command.Spec) error {
	s := spec.(*commandSpec)
	logp.Info("Running command... (name=%s)", s.name)
	return nil
}

func (p *publishHandler) HandleStdOut(spec command.Spec, out string) error {
	s := spec.(*commandSpec)
	p.publisher.LogDebug(s, "<stdout>%s", out)
	v, err := s.parser.Parse(out)
	if err != nil {
		logp.Err("failed to parse the line got from stdin. (command=%s, args=%v, line=%s, err=%s)", spec.Command, spec.Args, out, err)
		return nil
	}
	p.publisher.LogDebug(s, "<parsed>%#v", v)
	p.publisher.Publish(s, v)
	return nil
}

func (p *publishHandler) HandleStdErr(spec command.Spec, err string) error {
	return nil
}

func (p *publishHandler) AfterExit(spec command.Spec, status int) error {
	s := spec.(*commandSpec)
	logp.Info("Finished command. (name=%s, status=%d)", s.name, status)
	return nil
}

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
