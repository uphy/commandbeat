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
		Publish(spec *command.Spec, v common.MapStr)
	}
	elasticPublisher struct {
		client beat.Client
	}
)

func newPublishHandler(publisher *elasticPublisher) *publishHandler {
	return &publishHandler{publisher}
}

func (p *publishHandler) BeforeStart(spec *command.Spec) error {
	logp.Info("Running command... (name=%s)", spec.Name)
	return nil
}

func (p *publishHandler) HandleStdOut(spec *command.Spec, out string) error {
	p.publisher.LogDebug(spec, "<stdout>%s", out)
	v, err := spec.Parser.Parse(out)
	if err != nil {
		logp.Err("failed to parse the line got from stdin. (command=%s, args=%v, line=%s, err=%s)", spec.Command, spec.Args, out, err)
		return nil
	}
	p.publisher.LogDebug(spec, "<parsed>%#v", v)
	p.publisher.Publish(spec, v)
	return nil
}

func (p *publishHandler) HandleStdErr(spec *command.Spec, err string) error {
	return nil
}

func (p *publishHandler) AfterExit(spec *command.Spec, status int) error {
	logp.Info("Finished command. (name=%s, status=%d)", spec.Name, status)
	return nil
}

func newElasticPublisher(client beat.Client) *elasticPublisher {
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
