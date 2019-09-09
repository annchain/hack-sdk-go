package hackSDK

import (
	"encoding/hex"
	"fmt"
)

func HexToBytes(s string) ([]byte, error) {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		return nil, fmt.Errorf("incorrect hex length, should be EVEN.")
	}
	return hex.DecodeString(s)
}

func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}
