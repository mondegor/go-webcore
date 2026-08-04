package parser

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	// Validator - парсер полей структур с валидацией указанной в тегах полей.
	Validator struct {
		decoder   request.ParserDecode
		validator mrview.Validator
	}
)

// NewValidator - создаёт объект Validator.
func NewValidator(
	decoder request.ParserDecode,
	validator mrview.Validator,
) *Validator {
	return &Validator{
		decoder:   decoder,
		validator: validator,
	}
}

// Validate - возвращает в structPointer распарсеный внешний запрос или ошибку,
// если валидация запроса не прошла.
func (p *Validator) Validate(r *http.Request, structPointer any) error {
	return p.parseAndValidate(r.Context(), r.Body, structPointer)
}

// ValidateContent - возвращает в structPointer распарсенный []byte или ошибку,
// если валидация запроса не прошла.
func (p *Validator) ValidateContent(ctx context.Context, content []byte, structPointer any) error {
	return p.parseAndValidate(ctx, bytes.NewReader(content), structPointer)
}

// ValidateStruct - возвращает результат проверки заранее подготовленной структуры с данными.
func (p *Validator) ValidateStruct(ctx context.Context, structPointer any) error {
	return p.validator.Validate(ctx, structPointer)
}

func (p *Validator) parseAndValidate(ctx context.Context, r io.Reader, structPointer any) error {
	if err := p.decoder.ParseToStruct(ctx, r, structPointer); err != nil {
		return err
	}

	return p.validator.Validate(ctx, structPointer)
}
