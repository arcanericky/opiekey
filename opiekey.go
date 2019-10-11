package opiekey

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/md4"
)

// Algorithm is the hashing algorithm
//type Algorithm int8

const (
	algMD5 = iota
	algMD4
	algSHA1
)

type Algorithm struct {
	id   int8
	text string
}

// String returns the algorithm name
func (a Algorithm) String() string {
	return a.text
}

func (a Algorithm) hasher() func(d []byte) []byte {
	var f func(d []byte) []byte

	switch a.id {
	case algMD4:
		f = func(d []byte) []byte {
			h := md4.New()
			h.Write(d)
			return h.Sum(nil)
		}
	case algMD5:
		f = func(d []byte) []byte {
			h := md5.Sum(d)
			return h[:]
		}
	case algSHA1:
		f = func(d []byte) []byte {
			h := sha1.Sum(d)
			return h[:]
		}
	}

	return f
}

// MD5 is the MD5 algorithm
var MD5 = Algorithm{algMD5, "MD5"}

// MD4 is the MD4 algorithm
var MD4 = Algorithm{algMD4, "MD4"}

// SHA1 is the SHA1 algorithm
var SHA1 = Algorithm{algSHA1, "SHA1"}

func hashLen(seqNum int, seed, passphrase string, alg func([]byte) []byte) []byte {
	hash := []byte(seed + passphrase)
	r := make([]uint32, 5)
	q := make([]byte, 8)

	for i := 0; i <= seqNum; i++ {
		checksum := alg(hash)
		checksumLen := len(checksum)

		for j := 0; j < checksumLen/4; j++ {
			r[j] = binary.LittleEndian.Uint32(checksum[j*4:])
		}

		if checksumLen > 16 {
			binary.LittleEndian.PutUint32(q[0:], r[0]^r[2]^r[4])
		} else {
			binary.LittleEndian.PutUint32(q[0:], r[0]^r[2])
		}
		binary.LittleEndian.PutUint32(q[4:], r[1]^r[3])

		hash = q
	}

	return hash
}

func hashToUint64(hash []byte) uint64 {
	var sum uint64

	for _, i := range hash {
		sum = sum << 8
		sum |= uint64(i & 0xff)
	}

	return sum
}

func binaryToHex(sum uint64) string {
	h := fmt.Sprintf("%016X", sum)

	s := ""
	for i := 0; i < len(h); i += 4 {
		if len(s) > 0 {
			s = s + " "
		}

		s = s + h[i:i+4]
	}

	return s
}

func binaryToWords(sum uint64) string {
	var p uint64
	workSum := sum

	for i := 0; i < 32; i++ {
		p = p + workSum&3
		workSum = workSum >> 2
	}

	s := ""
	for i := 4; i >= 0; i-- {
		s = s + opieWords[(sum>>uint((i*11+9))&0x7ff)] + " "
	}

	s = s + opieWords[((sum<<2)&0x7fc)|(p&3)]

	return s
}

func getOTPNumber(seqNum int, seed, passphrase string, alg Algorithm) uint64 {
	hash := hashLen(seqNum, seed, passphrase, alg.hasher())
	sum := hashToUint64(hash)

	return sum
}

// ComputeHexResponse computes the response to an OTP challenge
func ComputeHexResponse(seqNum int, seed, passphrase string, alg Algorithm) string {
	sum := getOTPNumber(seqNum, seed, passphrase, alg)
	s := binaryToHex(sum)

	return s
}

// ComputeWordResponse computes the response to an OTP challenge
func ComputeWordResponse(seqNum int, seed, passphrase string, alg Algorithm) string {
	sum := getOTPNumber(seqNum, seed, passphrase, alg)
	s := binaryToWords(sum)

	return s
}
