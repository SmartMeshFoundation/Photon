package utils

import (
	"math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"

	"bytes"
	"encoding/gob"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kataras/iris/utils"
)

// BytesToString accepts bytes and returns their string presentation
// instead of string() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

// StringToBytes accepts string and returns their []byte presentation
// instead of byte() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// Random takes a parameter (int) and returns random slice of byte
// ex: var randomstrbytes []byte; randomstrbytes = utils.Random(32)
func Random(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// RandomString accepts a number(10 for example) and returns a random string using simple but fairly safe random algorithm
func RandomString(n int) string {
	return string(Random(n))
}
func NewRandomAddress() common.Address {
	hash := Sha3([]byte(utils.Random(10)))
	return common.BytesToAddress(hash[12:])
}

func StringInterface(i interface{}, depth int) string {
	c := spew.Config
	spew.Config.DisableMethods = false
	spew.Config.MaxDepth = depth
	s := spew.Sdump(i)
	spew.Config = c
	return s
}
func StringInterface1(i interface{}) string {
	c := spew.Config
	spew.Config.DisableMethods = false
	spew.Config.MaxDepth = 1
	s := spew.Sdump(i)
	spew.Config = c
	return s
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// Exists returns true if directory||file exists
func Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func RandomGenerator() common.Hash {
	return Sha3(utils.Random(32))
}
