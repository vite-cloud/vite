package util

import (
	"fmt"
	"os"
)

type Color [3]uint8

var Blue Color = [3]uint8{59, 130, 246}

var Red Color = [3]uint8{239, 68, 68}
var Yellow Color = [3]uint8{245, 158, 11}
var Green Color = [3]uint8{16, 185, 129}

var White Color = [3]uint8{255, 255, 255}
var Gray Color = [3]uint8{147, 148, 153}

var AnsiEnabled = true

func checkAnsi(args []string) bool {
	if os.Getenv("GOOS") == "windows" {
		return false
	}

	if os.Getenv("TERM") == "dumb" {
		return false
	}

	for _, arg := range args {
		if arg == "--no-ansi" {
			return false
		}
	}

	return true
}

func init() {
	AnsiEnabled = checkAnsi(os.Args)
}

func (c Color) String() string {
	return c.Fg()
}

func (c Color) Fg() string {
	if !AnsiEnabled {
		return ""
	}

	return fmt.Sprintf("\x1b[1m\x1b[38;2;%d;%d;%dm", c[0], c[1], c[2])
}

func (c Color) Bg() string {
	if !AnsiEnabled {
		return ""
	}

	return fmt.Sprintf("\x1b[1m\x1b[48;2;%d;%d;%dm", c[0], c[1], c[2])
}

func Reset() string {
	if !AnsiEnabled {
		return ""
	}

	return "\x1b[0m"
}
