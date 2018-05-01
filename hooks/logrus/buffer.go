package logrus

import (
	"bytes"
)

func bufWriteKV(buf *bytes.Buffer, k string, v []byte, sep, l, r byte) (
	err error,
) {
	_, err = buf.WriteString(k)
	if err != nil {
		return
	}
	err = buf.WriteByte(sep)
	if err != nil {
		return
	}
	if l > 0 {
		err = buf.WriteByte(l)
		if err != nil {
			return
		}
	}
	for _, c := range v {
		if c == '\\' {
			err = buf.WriteByte('\\')
			if err != nil {
				return
			}
			err = buf.WriteByte('\\')
			if err != nil {
				return
			}
			continue
		}
		if c == '"' {
			err = buf.WriteByte('\\')
			if err != nil {
				return
			}
			err = buf.WriteByte('"')
			if err != nil {
				return
			}
			continue
		}
		err = buf.WriteByte(c)
		if err != nil {
			return
		}
	}

	if r > 0 {
		err = buf.WriteByte(r)
	}

	return
}
