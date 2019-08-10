package utils

import (
	"io"
)

type ioUtil struct {
}

var IOUtil = ioUtil{}

func (io *ioUtil) ReadFully(conn io.ReadCloser, b []byte, length uint32) error {

	var buf []byte
	if length < 4096 {
		buf = make([]byte, length)
	} else {
		buf = make([]byte, 4096)
	}
	var byte_readed uint32

	buf_len := uint32(len(buf))
	var err error
	var reqLen int
	var _tmp []byte
	for byte_readed < length {

		if length-byte_readed < buf_len {
			_tmp = make([]byte, length-byte_readed)
			reqLen, err = conn.Read(_tmp)
		} else {
			reqLen, err = conn.Read(buf)
			_tmp = buf
		}

		if err != nil {
			return err
		}
		var i uint32
		buf_len := uint32(reqLen)

		for i < buf_len {
			b[byte_readed+i] = _tmp[i]
			i++
		}
		byte_readed += buf_len
	}
	return nil
}

func (io *ioUtil) WriteFully(conn io.WriteCloser, b []byte) error {

	total_length := uint32(len(b))

	var byte_writed uint32
	for byte_writed < total_length {

		if n, err := conn.Write(b[byte_writed:]); err != nil {
			return err
		} else {
			byte_writed += uint32(n)
		}
	}

	return nil
}
