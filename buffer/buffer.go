package buffer

import (
	"fmt"
)

type circ_buf struct {
	end int
	buffer []byte
	ret_buffer []byte
	read_pos int
	write_pos int
	count int
	cr_count int
}


func Make(size int, ret_size int) *circ_buf{
	p := circ_buf {
		end: size -1,
		buffer: make([]byte, size),
		ret_buffer: make([]byte, ret_size),
		read_pos: 0,
		write_pos: 0,
		count: 0,
		cr_count: 0,
	}
	return &p
}

func (cb *circ_buf) Write_byte(b byte) {
	cb.buffer[cb.write_pos] = b
	cb.write_pos++
	if cb.write_pos >= cb.end {
		cb.write_pos = 0
	}
	cb.count++
	if b == 13 {
		cb.cr_count++
	}
}

func (cb *circ_buf) ReadString() string {
	if cb.cr_count > 0 {
		for i:=0; i < cb.count; i++ {
			b, _ := cb.Read_byte()
			if b != 13 {
				cb.ret_buffer[i] = b
			} else {
				cb.cr_count--
				return string(cb.ret_buffer[:i])
			} 
		}
		return ""

	} else {
		return ""
	}

}

func (cb *circ_buf) Read_byte() (byte, error) {
	if cb.count == 0 {
		return 0, fmt.Errorf("Empty")
	}
	if cb.read_pos >= cb.end {
		cb.read_pos = 0
	}
	b := cb.buffer[cb.read_pos] 
	cb.read_pos++
	cb.count--
	return b, nil
}

