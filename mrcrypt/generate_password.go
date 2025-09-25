package mrcrypt

import (
	"crypto/rand"
	"math"
)

// Символы используемые в пароле.
const (
	PassVowels      PassCharsKinds = 1  // гласные буквы
	PassConsonants  PassCharsKinds = 2  // согласные буквы
	PassNumerals    PassCharsKinds = 4  // цифры
	PassSigns       PassCharsKinds = 8  // знаки
	PassAbc         PassCharsKinds = 3  // PassVowels + PassConsonants
	PassAbcNumerals PassCharsKinds = 7  // PassVowels + PassConsonants + PassNumerals
	PassAll         PassCharsKinds = 15 // PassVowels + PassConsonants + PassNumerals + PassSigns
)

const (
	pwCharSetLen = 4
)

type (
	// PassCharsKinds - вид символов используемых в пароле.
	PassCharsKinds uint8

	pwCharSet struct {
		kind            PassCharsKinds
		successivelyMax int
		firstOrLast     bool
		lettersLen      uint8
		letters         []byte
	}

	// PasswordGenerator - библиотека для генерации стоковых последовательностей.
	PasswordGenerator struct {
		pwCharSets [pwCharSetLen]pwCharSet
	}
)

// NewPasswordGenerator - создаёт объект PasswordGenerator.
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{
		pwCharSets: [pwCharSetLen]pwCharSet{
			{PassVowels, 2, true, 10, []byte("aeiuyAEIUY")}, // oO - символы удалены, чтобы не перепутать с нулём
			{PassConsonants, 2, true, 40, []byte("bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ")},
			{PassNumerals, 1, false, 9, []byte("123456789")}, // 0 - символ удалён, чтобы не перепутать с символами oO
			{PassSigns, 1, false, 12, []byte("!$%&.<=>?@_~")},
		},
	}
}

// Generate - comment func.
func (pg *PasswordGenerator) Generate(length int, charsKinds PassCharsKinds) string {
	if length < 1 {
		length = 1
	}

	if charsKinds == 0 || charsKinds > PassAll {
		charsKinds = PassAll
	}

	var (
		abc    [pwCharSetLen]pwCharSet
		abcLen uint8
	)

	for i := 0; i < pwCharSetLen; i++ {
		if (pg.pwCharSets[abcLen].kind & charsKinds) > 0 {
			abc[abcLen] = pg.pwCharSets[i]
			abcLen++
		}
	}

	// если указан только один набор символов
	if abcLen == 1 {
		abc[0].successivelyMax = length // максимальная длина совпадает с длиной пароля
		abc[0].firstOrLast = true       // первый и последний символ не проверяется
	}

	result := make([]byte, length)

	lastAbc := struct {
		charSetIndex           uint8
		countSuccessivelySigns int
	}{}

	for i := 0; i < length; i++ {
		var abcIndex uint8

		for {
			abcIndex = pg.getRandValue(abcLen)

			// если выбранный тип можно использовать для генерации первого и последнего символа
			// или если символ не первый и не последний
			if abc[abcIndex].firstOrLast || (i != 0 && i != (length-1)) {
				// если предыдущий символ такого же типа
				if abcIndex != lastAbc.charSetIndex {
					lastAbc.charSetIndex = abcIndex
					lastAbc.countSuccessivelySigns = 1

					break
				}

				// если подряд идущих символов не превышает
				if lastAbc.countSuccessivelySigns < abc[abcIndex].successivelyMax {
					lastAbc.countSuccessivelySigns++

					break
				}
			}
		}

		// обращение к случайному символу типа
		result[i] = abc[abcIndex].letters[pg.getRandValue(abc[abcIndex].lettersLen)]
	}

	return string(result)
}

func (pg *PasswordGenerator) getRandValue(maxValue uint8) uint8 {
	var tmp [1]byte

	bits100 := uint64(math.Log2(float64(maxValue)) * 100)
	bits := bits100 / 100

	if bits100%100 != 0 {
		bits++
	}

	mask := uint8(1<<bits) - 1

	for {
		if _, err := rand.Read(tmp[:]); err != nil {
			return 0
		}

		rnd := tmp[0] & mask

		if rnd < maxValue {
			return rnd
		}
	}
}
