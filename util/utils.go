package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
