package compression

import (
	"bytes"
	"context"
	"io"

	"github.com/anicoll/rp-plugin/pkg/compression"
	"github.com/redpanda-data/benthos/v4/public/service"
)

var compressionConfigSpec = service.NewConfigSpec().
	Summary("compression processor that compresses the payload of the message.")

func newCompressionProcessor(_ *service.ParsedConfig, logger *service.Logger) (*compressionProcessor, error) {
	return &compressionProcessor{
		log: logger,
	}, nil
}

func init() {
	err := service.RegisterProcessor(
		"compression", compressionConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
			return newCompressionProcessor(conf, mgr.Logger())
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type compressionProcessor struct {
	log *service.Logger
}

func (p *compressionProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
	data, err := m.AsBytes()
	if err != nil {
		p.log.Errorf("failed to get message data AsBytes: %s", err.Error())
		return nil, err
	}

	buffer := new(bytes.Buffer)
	wc, err := compression.NewWriteCloser(buffer)
	if err != nil {
		p.log.Errorf("zstd.NewWriter: " + err.Error())
		return nil, err
	}
	if _, err := io.Copy(wc, bytes.NewReader(data)); err != nil {
		p.log.Errorf("io.Copy: " + err.Error())
		return nil, err
	}
	if err := wc.Close(); err != nil {
		p.log.Errorf("zstdEncoder.Close: " + err.Error())
		return nil, err
	}

	msg := service.NewMessage(buffer.Bytes())

	// append metadata to downstream message.
	if err := m.MetaWalk(func(k, v string) error {
		msg.MetaSet(k, v)
		return nil
	}); err != nil {
		return nil, err
	}

	return []*service.Message{msg}, nil
}

func (p *compressionProcessor) Close(_ context.Context) error {
	return nil
}
