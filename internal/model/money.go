package model

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatKopecks переводит копейки в читаемую строку: 2299994 -> "22 999,94".
func FormatKopecks(kopecks int64) string {
	sign := ""
	if kopecks < 0 {
		sign = "-"
		kopecks = -kopecks
	}

	rubles := strconv.FormatInt(kopecks/100, 10)

	var b strings.Builder
	for i, r := range rubles {
		if i > 0 && (len(rubles)-i)%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}

	return fmt.Sprintf("%s%s,%02d", sign, b.String(), kopecks%100)
}
