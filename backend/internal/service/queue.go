package service

import (
	"errors"
	"fmt"
)

var ErrQueueOutOfRange = errors.New("queue index out of range")

func FormatQueueNumber(index int) (string, error) {
	if index < 0 || index >= 260 {
		return "", ErrQueueOutOfRange
	}

	letter := rune('A' + index/10)
	number := index % 10

	return fmt.Sprintf("%c%d", letter, number), nil
}
