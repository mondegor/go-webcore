package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	Validator struct {
		decoder   mrserver.RequestDecoder
		validator mrview.Validator
	}
)

// Make sure the Validator conforms with the mrserver.RequestParserValidate interface
var _ mrserver.RequestParserValidate = (*Validator)(nil)

func NewValidator(
	decoder mrserver.RequestDecoder,
	validator mrview.Validator,
) *Validator {
	return &Validator{
		decoder:   decoder,
		validator: validator,
	}
}

func (p *Validator) Validate(r *http.Request, structPointer any) error {
	if err := p.decoder.ParseToStruct(r, structPointer); err != nil {
		return err
	}

	return p.validator.Validate(r.Context(), structPointer)
}
