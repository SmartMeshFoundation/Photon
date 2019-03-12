package utils

import (
	"errors"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"time"
	"unsafe"

	"bytes"
	"encoding/gob"

	"crypto/ecdsa"

	"runtime/debug"

	"fmt"

	rand2 "crypto/rand"
	"io"

	"encoding/base32"

	"encoding/binary"

	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// BytesToString accepts bytes and returns their string presentation
// instead of string() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func BytesToString(b []byte) string {
	/* #nosec */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	/* #nosec */
	return *(*string)(unsafe.Pointer(&sh))
}

// StringToBytes accepts string and returns their []byte presentation
// instead of byte() this method doesn't generate memory allocations,
// BUT it is not safe to use anywhere because it points
// this helps on 0 memory allocations
func StringToBytes(s string) []byte {
	/* #nosec */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  0,
	}
	/* #nosec */
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//RandSrc random source from math
var RandSrc = rand.NewSource(time.Now().UnixNano())

func readFullOrPanic(r io.Reader, v []byte) int {
	n, err := io.ReadFull(r, v)
	if err != nil {
		panic(err)
	}
	return n
}

// Random takes a parameter (int) and returns random slice of byte
// ex: var randomstrbytes []byte; randomstrbytes = utils.Random(32)
func Random(n int) []byte {
	v := make([]byte, n)
	readFullOrPanic(rand2.Reader, v)
	return v
}

// RandomString accepts a number(10 for example) and returns a random string using simple but fairly safe random algorithm
func RandomString(n int) string {
	s := base32.StdEncoding.EncodeToString(Random(n))
	return s[:n]
}

//NewRandomInt generate a random int ,not more than n
func NewRandomInt(n int) int {
	return rand.New(RandSrc).Intn(n)
}

//NewRandomInt64 generate a random int64
func NewRandomInt64() int64 {
	return rand.New(RandSrc).Int63()
}

//NewRandomAddress generate a address,there maybe no corresponding priv key
func NewRandomAddress() common.Address {
	hash := Sha3([]byte(Random(10)))
	return common.BytesToAddress(hash[12:])
}

//NewRandomHash generate random hash,for testonly
func NewRandomHash() common.Hash {
	return Sha3(Random(64))
}

//MakePrivateKeyAddress generate a private key and it's address
func MakePrivateKeyAddress() (*ecdsa.PrivateKey, common.Address) {
	//#nosec
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

//StringInterface use spew to string any object with max `depth`,it's not thread safe.
func StringInterface(i interface{}, depth int) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	c := spew.Config
	spew.Config.DisableMethods = false
	//spew.Config.ContinueOnMethod = false
	spew.Config.MaxDepth = depth
	s := spew.Sdump(i)
	spew.Config = c
	return s
}

//StringInterface1 use spew to string any object with depth 1,it's not thread safe.
func StringInterface1(i interface{}) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	c := spew.Config
	spew.Config.DisableMethods = false
	spew.Config.MaxDepth = 1
	s := spew.Sdump(i)
	spew.Config = c
	return s
}

//DeepCopy use gob to copy
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

// GetHomePath returns the user's $HOME directory
func GetHomePath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

//SystemExit quit and print stack
func SystemExit(code int) {
	if code != 0 {
		debug.PrintStack()
	}
	time.Sleep(time.Second * 2)
	os.Exit(code)
}

// ReadVarInt reads a variable length integer from r and returns it as a uint64.
func ReadVarInt(r io.Reader) (uint64, error) {
	var discriminant uint8
	err := binary.Read(r, binary.LittleEndian, &discriminant)
	if err != nil {
		return 0, err
	}
	var rv uint64
	switch discriminant {
	case 0xff:
		var sv uint64
		err := binary.Read(r, binary.LittleEndian, &sv)
		if err != nil {
			return 0, err
		}
		rv = sv

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x100000000)
		if rv < min {
			return 0, errors.New("ReadVarInt wrong")
		}

	case 0xfe:
		var sv uint32
		err := binary.Read(r, binary.LittleEndian, &sv)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x10000)
		if rv < min {
			return 0, errors.New("ReadVarInt wrong")
		}

	case 0xfd:
		var sv uint16
		err := binary.Read(r, binary.LittleEndian, &sv)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0xfd)
		if rv < min {
			return 0, errors.New("ReadVarInt wrong")
		}

	default:
		rv = uint64(discriminant)
	}

	return rv, nil
}

// WriteVarInt serializes val to w using a variable number of bytes depending
// on its value.
func WriteVarInt(w io.Writer, val uint64) error {
	if val < 0xfd {
		return binary.Write(w, binary.LittleEndian, uint8(val))
	}

	if val <= math.MaxUint16 {
		err := binary.Write(w, binary.LittleEndian, uint8(0xfd))
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(val))
	}

	if val <= math.MaxUint32 {
		err := binary.Write(w, binary.LittleEndian, uint8(0xfe))
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(val))
	}

	err := binary.Write(w, binary.LittleEndian, uint8(0xff))
	if err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, val)
}

// ToJSONFormat :
func ToJSONFormat(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
