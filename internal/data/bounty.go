package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidBerriesFormat = errors.New("invalid berries format: supported format is ('300M berries', '1.5B berries', '100 berries')")

type Berries int64

func (b Berries) MarshalJSON() ([]byte, error) {
	var jsonvalue string

	if b >= 1000000000 { // 1 billion or more
		billions := float64(b) / 1000000000
		if billions == float64(int64(billions)) {
			jsonvalue = fmt.Sprintf("%.0fB berries", billions)
		} else {
			jsonvalue = fmt.Sprintf("%.1fB berries", billions)
		}
	} else if b >= 1000000 { // 1 million or more
		millions := float64(b) / 1000000
		if millions == float64(int64(millions)) {
			jsonvalue = fmt.Sprintf("%.0fM berries", millions)
		} else {
			jsonvalue = fmt.Sprintf("%.1fM berries", millions)
		}
	} else {
		jsonvalue = fmt.Sprintf("%d berries", b)
	}

	quotedJSONValue := strconv.Quote(jsonvalue)
	return []byte(quotedJSONValue), nil
}

func (b *Berries) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidBerriesFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "berries" {
		return ErrInvalidBerriesFormat
	}

	valueStr := parts[0]
	var multiplier int64 = 1

	if strings.HasSuffix(valueStr, "B") {
		multiplier = 1000000000
		valueStr = strings.TrimSuffix(valueStr, "B")
	} else if strings.HasSuffix(valueStr, "M") {
		multiplier = 1000000
		valueStr = strings.TrimSuffix(valueStr, "M")
	}

	if multiplier == 1 {
		i, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return ErrInvalidBerriesFormat
		}
		*b = Berries(i)
	} else {
		f, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return ErrInvalidBerriesFormat
		}
		*b = Berries(int64(f * float64(multiplier)))
	}

	return nil
}
