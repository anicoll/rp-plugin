package compression

import (
	"context"

	"github.com/anicoll/rp-plugin/pkg/compression"
	"github.com/redpanda-data/benthos/v4/public/service"
)

var deCompressionConfigSpec = service.NewConfigSpec().
	Summary("de_compression processor that decompresses the payload of the message.")

func newDeCompressionProcessor(_ *service.ParsedConfig, logger *service.Logger) (*deCompressionProcessor, error) {
	rc, err := compression.NewReadCloser()
	if err != nil {
		return nil, err
	}
	return &deCompressionProcessor{
		readCloser: rc,
		log:        logger,
	}, nil
}

func init() {
	err := service.RegisterProcessor(
		"de_compression", deCompressionConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
			return newDeCompressionProcessor(conf, mgr.Logger())
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type deCompressionProcessor struct {
	readCloser compression.ReadCloser
	log        *service.Logger
}

func (p *deCompressionProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
	data, err := m.AsBytes()
	if err != nil {
		p.log.Errorf("m.AsBytes: " + err.Error())
		return nil, err
	}

	decompressedData, err := p.readCloser.DecodeAll(data)
	if err != nil {
		p.log.Errorf("p.readCloser: " + err.Error())
		return nil, err
	}

	msg := service.NewMessage(decompressedData)

	// append metadata to downstream message.
	if err := m.MetaWalk(func(k, v string) error {
		msg.MetaSet(k, v)
		return nil
	}); err != nil {
		return nil, err
	}

	return []*service.Message{msg}, nil
}

func (p *deCompressionProcessor) Close(ctx context.Context) error {
	return p.readCloser.Close()
}
