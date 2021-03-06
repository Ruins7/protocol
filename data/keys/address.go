package keys

import (
	"bytes"
	"encoding/hex"
	"strings"
)

//Address to be used as to reference a key-pair.
type Address []byte

func (a Address) String() string {
	return "0x" + hex.EncodeToString(a)
}

func (a Address) Bytes() []byte {
	return a
}

func (a Address) Equal(b Address) bool {
	return bytes.Equal(a, b)
}

func (a Address) Humanize() string {
	return strings.ToUpper(hex.EncodeToString(a))
}
