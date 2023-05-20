// SPDX-FileCopyrightText: 2019 Kent Gibson <warthog618@gmail.com>
//
// SPDX-License-Identifier: MIT

// Package bananapi provides convenience mappings from Banana Pi pin names to
// offsets.
package orangepi

import (
	"errors"
	"strconv"
	"strings"
)

var (
	sunxi_GPA = 0
	// sunxi_GPB = 32
	// sunxi_GPC = 64
	// sunxi_GPD = 96
	// sunxi_GPE = 128
	// sunxi_GPF = 160
	sunxi_GPG = 192
	// BY CHOW  ADD
	// sunxi_GPH = 224
	// sunxi_GPI = 256
	// sunxi_GPJ = 288
	// sunxi_GPK = 320
	sunxi_GPL = 352
	// sunxi_GPM = 384
	// sunxi_GPN = 448
	// sunxi_GPO = 448 + 32
)

// GPIO aliases to offsets
var (
	POWER_LED  = GPIO_TO_OFFSET[1]
	STATUS_LED = GPIO_TO_OFFSET[2]

	PA12 = GPIO_TO_OFFSET[3]
	PA11 = GPIO_TO_OFFSET[5]
	PA6  = GPIO_TO_OFFSET[7]
	PA1  = GPIO_TO_OFFSET[11]
	PA0  = GPIO_TO_OFFSET[13]
	PA3  = GPIO_TO_OFFSET[15]
	PA15 = GPIO_TO_OFFSET[19]
	PA16 = GPIO_TO_OFFSET[21]
	PA14 = GPIO_TO_OFFSET[23]
	PG6  = GPIO_TO_OFFSET[8]
	PG7  = GPIO_TO_OFFSET[10]
	PA7  = GPIO_TO_OFFSET[12]
	PA19 = GPIO_TO_OFFSET[16]
	PA18 = GPIO_TO_OFFSET[18]
	PA2  = GPIO_TO_OFFSET[22]
	PA13 = GPIO_TO_OFFSET[24]
	PA10 = GPIO_TO_OFFSET[26]
)

var GPIO_TO_OFFSET = map[int]int{
	1:  sunxi_GPL + 10,
	2:  sunxi_GPA + 17,
	3:  sunxi_GPA + 12,
	5:  sunxi_GPA + 11,
	7:  sunxi_GPA + 6,
	11: sunxi_GPA + 1,
	13: sunxi_GPA + 0,
	15: sunxi_GPA + 3,
	19: sunxi_GPA + 15,
	21: sunxi_GPA + 16,
	23: sunxi_GPA + 14,
	8:  sunxi_GPG + 6,
	10: sunxi_GPG + 7,
	12: sunxi_GPA + 7,
	16: sunxi_GPA + 19,
	18: sunxi_GPA + 18,
	22: sunxi_GPA + 2,
	24: sunxi_GPA + 13,
	26: sunxi_GPA + 10,
}

// ErrInvalid indicates the pin name does not match a known pin.
var ErrInvalid = errors.New("invalid pin number")

func rangeCheck(p int) (int, error) {
	if p < 2 || p >= 27 {
		return 0, ErrInvalid
	}
	return p, nil
}

// Pin maps a pin string name to a pin number.
//
// Pin names are case insensitive and may be of the form GPIOX, or X.
func Pin(s string) (int, error) {
	s = strings.ToLower(s)
	switch {
	case strings.HasPrefix(s, "gpio"):
		v, err := strconv.ParseInt(s[4:], 10, 8)
		if err != nil {
			return 0, err
		}
		return rangeCheck(int(v))
	default:
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return 0, err
		}
		return rangeCheck(int(v))
	}
}

// MustPin converts the string to the corresponding pin number or panics if that
// is not possible.
func MustPin(s string) int {
	v, err := Pin(s)
	if err != nil {
		panic(err)
	}
	return v
}
