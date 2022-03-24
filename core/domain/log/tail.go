package log

import (
	"github.com/hpcloud/tail"
	"github.com/icza/backscanner"
	"io"
	"os"
)

type TailOptions struct {
	Stream   bool
	Backfill int
}

func Tail(path string, opts TailOptions) (<-chan string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	scanner := backscanner.New(file, int(stat.Size()))

	offset := 0
	n := opts.Backfill

	for {
		if n == 0 {
			break
		}
		_, pos, err := scanner.LineBytes()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		n--
		offset = pos
	}

	stream, err := tail.TailFile(path, tail.Config{
		Location: &tail.SeekInfo{
			Offset: int64(offset),
			Whence: io.SeekStart,
		},
		ReOpen:    opts.Stream,
		MustExist: true,
		Follow:    opts.Stream,
		Logger:    tail.DiscardingLogger,
	})
	if err != nil {
		return nil, err
	}

	defer stream.Cleanup()

	ch := make(chan string)

	go func() {
		for line := range stream.Lines {
			ch <- line.Text
		}

		close(ch)
	}()

	return ch, nil
}
