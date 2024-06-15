package mrcrypt

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	PassVowels      PassCharsKinds = 1  // PassVowels  - гласные буквы
	PassConsonants  PassCharsKinds = 2  // PassConsonants  - согласные буквы
	PassNumerals    PassCharsKinds = 4  // PassNumerals  - цифры
	PassSigns       PassCharsKinds = 8  // PassSigns  - знаки
	PassAbc         PassCharsKinds = 3  // PassAbc = PassVowels + PassConsonants
	PassAbcNumerals PassCharsKinds = 7  // PassAbcNumerals = PassVowels + PassConsonants + PassNumerals
	PassAll         PassCharsKinds = 15 // PassAll = PassVowels + PassConsonants + PassNumerals + PassSigns

	pwCharSetLen = 4
)

type (
	// PassCharsKinds  - comment type.
	PassCharsKinds uint8

	pwCharSet struct {
		kind            PassCharsKinds
		successivelyMax int
		firstOrLast     bool
		lettersLen      int
		letters         []byte
	}
)

var pwCharSets = [pwCharSetLen]pwCharSet{
	{PassVowels, 2, true, 10, []byte("aeiuyAEIUY")}, // oO - символы удалены, чтобы не перепутать с нулём
	{PassConsonants, 2, true, 40, []byte("bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ")},
	{PassNumerals, 1, false, 9, []byte("123456789")}, // 0 - символ удалён, чтобы не перепутать с oO
	{PassSigns, 1, false, 6, []byte(".!?@$&")},
}

// GenPassword  - comment func.
func GenPassword(length int, charsKinds PassCharsKinds) string {
	if length < 1 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d < 1", length)).Send()
		length = 1
	}

	if length > 128 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d > 128", length)).Send()
		length = 128
	}

	if charsKinds == 0 {
		charsKinds = PassAll
		mrlog.Default().Warn().Err(errors.New("param 'charsKinds' is zero")).Send() //nolint:wsl
	}

	var (
		abc    [pwCharSetLen]pwCharSet
		abcLen int
	)

	for i := 0; i < pwCharSetLen; i++ {
		if (pwCharSets[abcLen].kind & charsKinds) > 0 {
			abc[abcLen] = pwCharSets[i]
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
		charSetIndex           int
		countSuccessivelySigns int
	}{}

	for i := 0; i < length; i++ {
		var abcIndex int

		for {
			abcIndex = getRandValue(abcLen)

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
		result[i] = abc[abcIndex].letters[getRandValue(abc[abcIndex].lettersLen)]
	}

	return string(result)
}

func getRandValue(max int) int {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		mrlog.Default().Error().Err(err).Send()

		return 0
	}

	return int(value.Int64())
}
