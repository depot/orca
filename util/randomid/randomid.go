package randomid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

func RandomID() string {
	now := time.Now()
	var buf [15]byte
	rand.Read(buf[:])
	return fmt.Sprintf("%d-%s", now.Nanosecond(), base64.URLEncoding.EncodeToString(buf[:]))
}
