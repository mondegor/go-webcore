package mrparser

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	// Validator - comment struct.
	Validator struct {
		decoder   mrserver.RequestDecoder
		validator mrview.Validator
	}
)

// Make sure the Validator conforms with the mrserver.RequestParserValidate interface.
var _ mrserver.RequestParserValidate = (*Validator)(nil)

// NewValidator - создаёт объект Validator.
func NewValidator(decoder mrserver.RequestDecoder, validator mrview.Validator) *Validator {
	return &Validator{
		decoder:   decoder,
		validator: validator,
	}
}

// Validate - comment method.
func (p *Validator) Validate(r *http.Request, structPointer any) error {
	return p.validate(r.Context(), r.Body, structPointer)
}

// ValidateContent - comment method.
func (p *Validator) ValidateContent(ctx context.Context, content []byte, structPointer any) error {
	return p.validate(ctx, bytes.NewReader(content), structPointer)
}

func (p *Validator) validate(ctx context.Context, r io.Reader, structPointer any) error {
	if err := p.decoder.ParseToStruct(ctx, r, structPointer); err != nil {
		return err
	}

	return p.validator.Validate(ctx, structPointer)
}
