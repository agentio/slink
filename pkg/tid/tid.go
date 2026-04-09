package tid

import (
	"slices"
	"strings"
	"time"
)

const base32Alphabet = "234567abcdefghijklmnopqrstuvwxyz"

// We cannot use encoding/base32 here because, even with the correct alphabet,
// it yields uint64 values that are >>1 (shifted one bit right) vs. the encodings
// we need for TIDs.
// e.g. 7fffffffffffffff instead of ffffffffffffffff for "jzzzzzzzzzzzz".

func uint64ForBase32String(s string) uint64 {
	var v uint64
	for i := range len(s) {
		c := strings.IndexByte(base32Alphabet, s[i])
		if c < 0 {
			return 0
		}
		v = (v << 5) | uint64(c&0x1F)
	}
	return v
}

func base32StringForUint64(v uint64) string {
	var b strings.Builder
	for range 13 { // 64 bits / 5 bits per char = 12.8, round up to 13
		b.WriteByte(base32Alphabet[v&0x1f])
		v = v >> 5
	}
	runes := []rune(b.String())
	slices.Reverse(runes)
	return string(runes)
}

func ParseTID(tid string) (time.Time, uint16, error) {
	v := uint64ForBase32String(tid)
	return time.UnixMicro(int64(v >> 10)).UTC(), uint16(v & 0x3FF), nil
}

func TID(t time.Time, clockid uint16) string {
	return base32StringForUint64(uint64(t.UnixMicro())<<10 | uint64(clockid&0x3ff))
}

const ATProtoTimeFormat = "2006-01-02T15:04:05.999999Z"
const ATProtoTimeSecondsFormat = "2006-01-02T15:04:05Z"

func FromString(s string, clockid uint16) (string, error) {
	t, err := time.Parse(ATProtoTimeFormat, s)
	if err != nil {
		return "", err
	}
	return TID(t, clockid), nil
}
