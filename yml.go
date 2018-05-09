package main

import (
	"bytes"
)

func main() {
	print(bytes.Equal(nil, []byte{1}))
}
