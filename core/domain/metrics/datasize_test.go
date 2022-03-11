package metrics

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestByteSize_String(t *testing.T) {
	table := []struct {
		in  ByteSize
		out string
	}{
		{0, "0 B"},
		{B, "1 B"},
		{KB, "1.0 KB"},
		{MB, "1.0 MB"},
		{GB, "1.0 GB"},
		{TB, "1.0 TB"},
		{PB, "1.0 PB"},
		{EB, "1.0 EB"},
		{400 * TB, "400.0 TB"},
		{2048 * MB, "2.0 GB"},
		{B + KB, "1.0 KB"},
		{MB + 20*KB, "1.0 MB"},
		{100*MB + KB, "100.0 MB"},
	}

	for _, tt := range table {
		assert.Equal(t, tt.in.String(), tt.out)
	}
}
