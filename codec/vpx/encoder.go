package vpx

import (
	"fmt"
	"image"

	"github.com/edaniels/gostream"

	"github.com/edaniels/golog"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/prop"
)

type encoder struct {
	codec  codec.ReadCloser
	img    image.Image
	debug  bool
	logger golog.Logger
}

type VCodec string

const (
	CodecVP8 VCodec = "V_VP8"
	CodecVP9 VCodec = "V_VP9"
)

func NewEncoder(codecType VCodec, width, height int, debug bool, logger golog.Logger) (gostream.Encoder, error) {
	enc := &encoder{debug: debug, logger: logger}

	var builder codec.VideoEncoderBuilder
	switch codecType {
	case CodecVP8:
		params, err := vpx.NewVP8Params()
		if err != nil {
			return nil, err
		}
		builder = &params

	case CodecVP9:
		params, err := vpx.NewVP9Params()
		if err != nil {
			return nil, err
		}
		builder = &params
	default:
		return nil, fmt.Errorf("[WARN] unsupported VPX codec: %s", codecType)
	}

	codec, err := builder.BuildVideoEncoder(enc, prop.Media{
		Video: prop.Video{
			Width:  width,
			Height: height,
		},
	})
	if err != nil {
		return nil, err
	}
	enc.codec = codec

	return enc, nil
}

func (v *encoder) Read() (img image.Image, release func(), err error) {
	return v.img, nil, nil
}

func (v *encoder) Encode(img image.Image) ([]byte, error) {
	v.img = img
	data, _, err := v.codec.Read()
	return data, err
}
