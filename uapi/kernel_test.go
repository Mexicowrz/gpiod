// SPDX-License-Identifier: MIT
//
// Copyright © 2019 Kent Gibson <warthog618@gmail.com>.

// +build linux

package uapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/gpiod/mockup"
	"github.com/warthog618/gpiod/uapi"
	"golang.org/x/sys/unix"
)

func TestRepeatedLines(t *testing.T) {
	t.Skip("leaves line as output as of 5.4-rc1")
	mockupRequired(t)
	c, err := mock.Chip(0)
	require.Nil(t, err)
	require.NotNil(t, c)
	f, err := os.Open(c.DevPath)
	require.Nil(t, err)

	hr := uapi.HandleRequest{
		Lines: 2,
	}
	hr.Offsets[0] = 1
	hr.Offsets[1] = 1

	// input
	err = uapi.GetLineHandle(f.Fd(), &hr)
	assert.NotNil(t, err)

	// output
	hr.Flags = uapi.HandleRequestOutput
	hr.DefaultValues[0] = 0
	hr.DefaultValues[1] = 1
	err = uapi.GetLineHandle(f.Fd(), &hr)
	assert.NotNil(t, err)

}

func TestOutputSets(t *testing.T) {
	t.Skip("contains known failures as of 5.4-rc1")
	mockupRequired(t)
	patterns := []struct {
		name string
		flag uapi.HandleFlag
	}{
		{"o", uapi.HandleRequestOutput},
		{"od", uapi.HandleRequestOutput | uapi.HandleRequestOpenDrain},
		{"os", uapi.HandleRequestOutput | uapi.HandleRequestOpenSource},
	}
	c, err := mock.Chip(0)
	require.Nil(t, err)
	line := 0
	for _, p := range patterns {
		for initial := 0; initial <= 1; initial++ {
			for toggle := 0; toggle <= 1; toggle++ {
				for activeLow := 0; activeLow <= 1; activeLow++ {
					final := initial
					if toggle == 1 {
						final ^= 0x01
					}
					flags := p.flag
					if activeLow == 1 {
						flags |= uapi.HandleRequestActiveLow
					}
					label := fmt.Sprintf("%s-%d-%d-%d(%d)", p.name, initial^1, initial, final, activeLow)
					tf := func(t *testing.T) {
						testLine(t, c, line, flags, initial, toggle)
					}
					t.Run(label, tf)
				}
			}
		}
	}
}

func testLine(t *testing.T, c *mockup.Chip, line int, flags uapi.HandleFlag, initial, toggle int) {
	t.Helper()
	// set mock initial - opposing default
	c.SetValue(line, initial^0x01)
	f, err := os.Open(c.DevPath)
	require.Nil(t, err)
	defer f.Close()
	// request line output
	hr := uapi.HandleRequest{
		Flags: flags,
		Lines: uint32(1),
	}
	hr.Offsets[0] = uint32(line)
	hr.DefaultValues[0] = uint8(initial)
	err = uapi.GetLineHandle(f.Fd(), &hr)
	require.Nil(t, err)
	if toggle != 0 {
		var hd uapi.HandleData
		hd[0] = uint8(initial ^ 0x01)
		err = uapi.SetLineValues(uintptr(hr.Fd), hd)
		assert.Nil(t, err, "can't set value 1")
		err = uapi.GetLineValues(uintptr(hr.Fd), &hd)
		assert.Nil(t, err, "can't get value 1")
		assert.Equal(t, uint8(initial^1), hd[0], "get value 1")
		hd[0] = uint8(initial)
		err = uapi.SetLineValues(uintptr(hr.Fd), hd)
		assert.Nil(t, err, "can't set value 2")
		err = uapi.GetLineValues(uintptr(hr.Fd), &hd)
		assert.Nil(t, err, "can't get value 2")
		assert.Equal(t, uint8(initial), hd[0], "get value 2")
		hd[0] = uint8(initial ^ 0x01)
		err = uapi.SetLineValues(uintptr(hr.Fd), hd)
		assert.Nil(t, err, "can't set value 3")
		err = uapi.GetLineValues(uintptr(hr.Fd), &hd)
		assert.Nil(t, err, "can't get value 3")
		assert.Equal(t, uint8(initial^1), hd[0], "get value 3")
	}
	// release
	unix.Close(int(hr.Fd))
}
