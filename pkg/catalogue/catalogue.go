package catalogue

import (
	"hash/crc32"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Generate
func Generate() string {
	crc32InUint32 := crc32.ChecksumIEEE([]byte(strconv.Itoa(rand.Int())))
	crc32InString := strconv.FormatUint(uint64(crc32InUint32), 16)
	return crc32InString
}
