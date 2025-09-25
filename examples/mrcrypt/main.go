package main

import (
	"fmt"
	"os"

	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-webcore/mrcrypt"
)

func main() {
	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	value, _ := mrcrypt.GenerateDigits(16)
	logger.Debug("GenerateDigits", "value", value)

	value, _ = mrcrypt.GenerateHex(16)
	logger.Info("GenTokenHex", "value", value)

	value, _ = mrcrypt.GenerateToken(64)
	logger.Info("GenerateToken", "value", value)

	valueBytes, _ := mrcrypt.GenerateBytes([]byte("abc123.,:"), 16)
	logger.Info("GenerateBytes", "value", string(valueBytes))

	pwgen := mrcrypt.NewPasswordGenerator()

	logger.Info("GenPassword", "password", pwgen.Generate(16, mrcrypt.PassAll))

	pw := pwgen.Generate(12, mrcrypt.PassAbc)
	logger.Info("PasswordStrength 12 abc", "password", pw, "strength", mrcrypt.PasswordStrength(pw))

	pw = pwgen.Generate(9, mrcrypt.PassAbcNumerals)
	logger.Info("PasswordStrength 9 abc+num", "password", pw, "strength", mrcrypt.PasswordStrength(pw))

	pw = pwgen.Generate(12, mrcrypt.PassAll)
	logger.Info("PasswordStrength 12 all", "password", pw, "strength", mrcrypt.PasswordStrength(pw))

	fmt.Println(mrcrypt.PasswordStrength("<rin>24zD*~"))
	fmt.Println(mrcrypt.PasswordStrength("<rin>24xX.vD"))
	fmt.Println(mrcrypt.PasswordStrength("12345aAlowD"))
	fmt.Println(mrcrypt.PasswordStrength("12345aAl.D"))
	fmt.Println(mrcrypt.PasswordStrength("123eeeeddggDDll"))
	fmt.Println(mrcrypt.PasswordStrength("1234567890a"))
	fmt.Println(mrcrypt.PasswordStrength("12345678.a"))
	fmt.Println(mrcrypt.PasswordStrength("123456D.a"))
	fmt.Println(mrcrypt.PasswordStrength("12345678D"))
	fmt.Println(mrcrypt.PasswordStrength("123456.D"))
	fmt.Println(mrcrypt.PasswordStrength("1234s.D"))
}
