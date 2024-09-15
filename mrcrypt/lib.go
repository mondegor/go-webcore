package mrcrypt

import (
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// Lib - библиотека для генерации стоковых последовательностей.
	Lib struct {
		pwCharSets [pwCharSetLen]pwCharSet
		logger     mrlog.Logger
	}
)

// New - создаёт объект Lib.
func New(logger mrlog.Logger) *Lib {
	return &Lib{
		pwCharSets: [pwCharSetLen]pwCharSet{
			{PassVowels, 2, true, 10, []byte("aeiuyAEIUY")}, // oO - символы удалены, чтобы не перепутать с нулём
			{PassConsonants, 2, true, 40, []byte("bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ")},
			{PassNumerals, 1, false, 9, []byte("123456789")}, // 0 - символ удалён, чтобы не перепутать с oO
			{PassSigns, 1, false, 6, []byte(".!?@$&")},
		},
		logger: logger,
	}
}
