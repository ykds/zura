package random

import (
	"crypto/rand"
	"math/big"
)

var (
	strSeed = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	numSeed = []byte("1234567890")
	mixSeed = append(strSeed, numSeed...)
)

func RandStr(len int) string {
	return random(strSeed, len)

}

func RandNum(len int) string {
	return random(numSeed, len)
}

func RandMix(len int) string {
	return random(mixSeed, len)
}

func random(seed []byte, length int) string {
	str := ""
	bigInt := big.NewInt(int64(len(seed)))
	for i := 0; i < length; i++ {
		i, _ := rand.Int(rand.Reader, bigInt)
		str += string(seed[i.Int64()])
	}
	return str
}
