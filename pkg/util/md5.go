package util

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "dousound"

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(secret))

	return hex.EncodeToString(m.Sum([]byte(value)))
}
