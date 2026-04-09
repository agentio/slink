package tid

func Offset(t string, i uint16) string {
	v := uint64ForBase32String(t)
	return base32StringForUint64(v + uint64(i))
}
