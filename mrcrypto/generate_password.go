package mrcrypto

import (
    "crypto/rand"
    "math/big"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    PassVowels PassCharsKinds = 1
    PassConsonants = 2
    PassNumerals = 4
    PassSigns = 8
    PassAbc = 3 // PassVowels + PassConsonants
    PassAbcNumerals = 7 // PassVowels + PassConsonants + PassNumerals
    PassAll = 15 // PassVowels + PassConsonants + PassNumerals + PassSigns

    pwCharSetLen = 4
)

type (
    PassCharsKinds uint8

    pwCharSet struct {
        kind PassCharsKinds
        successivelyMax int
        firstOrLast bool
        lettersLen int
        letters []byte
    }
)

var (
    pwCharSets = [pwCharSetLen]pwCharSet{
        {PassVowels, 2, true, 10, []byte("aeiuyAEIUY")}, // oO - символы удалены, чтобы не перепутать с нулём
        {PassConsonants, 2, true, 40, []byte("bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ")},
        {PassNumerals, 1, false, 9, []byte("123456789")}, // 0 - символ удалён, чтобы не перепутать с oO
        {PassSigns, 1, false, 6, []byte(".!?@$&")},
    }
)

func GenPassword(length int, charsKinds PassCharsKinds) string {
    if length < 1 {
        mrcore.LogError("param 'length': %d < 1", length)
        length = 1
    }

    if length > 128 {
        mrcore.LogError("param 'length': %d > 128", length)
        length = 128
    }

    if charsKinds == 0 {
        mrcore.LogError("param 'charsKinds' is zero", length)
        charsKinds = PassAll
    }

    var abc [pwCharSetLen]pwCharSet
    var abcLen int

    for i := 0 ; i < pwCharSetLen; i++  {
        if (pwCharSets[abcLen].kind & charsKinds) > 0 {
            abc[abcLen] = pwCharSets[i]
            abcLen++
        }
    }

    // если указан только один набор символов
    if abcLen == 1 {
        abc[0].successivelyMax = length // максимальная длина совпадает с длиной пароля
        abc[0].firstOrLast = true // первый и последний символ не проверяется
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
            if abc[abcIndex].firstOrLast || (i != 0 && i != (length - 1)) {
                // если предыдущий символ такого же типа
                if abcIndex == lastAbc.charSetIndex {
                    // если подряд идущих символов не превышает
                    if lastAbc.countSuccessivelySigns < abc[abcIndex].successivelyMax {
                        lastAbc.countSuccessivelySigns++
                        break
                    }
                } else {
                    lastAbc.charSetIndex = abcIndex
                    lastAbc.countSuccessivelySigns = 1
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
        mrcore.LogErr(err)
        return 0
    }

    return int(value.Int64())
}
