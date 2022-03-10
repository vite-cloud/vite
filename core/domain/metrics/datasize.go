package metrics

import (
	"fmt"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
	MB          = KB << 10
	GB          = MB << 10
	TB          = GB << 10
	PB          = TB << 10
	EB          = PB << 10
)

func (b ByteSize) KB() float64 {
	v := b / KB
	r := b % KB
	return float64(v) + float64(r)/float64(KB)
}

func (b ByteSize) MB() float64 {
	v := b / MB
	r := b % MB
	return float64(v) + float64(r)/float64(MB)
}

func (b ByteSize) GB() float64 {
	v := b / GB
	r := b % GB
	return float64(v) + float64(r)/float64(GB)
}

func (b ByteSize) TB() float64 {
	v := b / TB
	r := b % TB
	return float64(v) + float64(r)/float64(TB)
}

func (b ByteSize) PB() float64 {
	v := b / PB
	r := b % PB
	return float64(v) + float64(r)/float64(PB)
}

func (b ByteSize) EB() float64 {
	v := b / EB
	r := b % EB
	return float64(v) + float64(r)/float64(EB)
}

func (b ByteSize) String() string {
	switch {
	case b >= EB:
		return fmt.Sprintf("%.1f EB", b.EB())
	case b >= PB:
		return fmt.Sprintf("%.1f PB", b.PB())
	case b >= TB:
		return fmt.Sprintf("%.1f TB", b.TB())
	case b >= GB:
		return fmt.Sprintf("%.1f GB", b.GB())
	case b >= MB:
		return fmt.Sprintf("%.1f MB", b.MB())
	case b >= KB:
		return fmt.Sprintf("%.1f KB", b.KB())
	default:
		return fmt.Sprintf("%d B", b)
	}
}
