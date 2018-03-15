package beater

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/uphy/commandbeat/command"
)

type (
	publishHandler struct {
		publisher *elasticPublisher
	}
)

func NewPublishHandler(publisher *elasticPublisher) *publishHandler {
	return &publishHandler{publisher}
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
