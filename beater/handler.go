package beater

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/uphy/commandbeat/command"
)

type (
	// publishHandler is command result handler which parses and publish to elasticsearch.
	publishHandler struct {
		publisher *elasticPublisher
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
