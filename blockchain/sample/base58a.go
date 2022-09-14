package sample

import (
	"bytes"
	"fmt"
	"math/big"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 对x取余数
		result = append(result, b58Alphabet[mod.Int64()])
	}
	ReverseBytes(result)
	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b == '1' {
			zeroBytes++
		} else {
			break
		}
	}
	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)     //反推出余数
		result.Mul(result, big.NewInt(58))               //之前的结果乘以58
		result.Add(result, big.NewInt(int64(charIndex))) //加上这个余数
	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{0x00}, zeroBytes), decoded...)
	return decoded
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func Base58aSample() {
	//org := []byte("qwerty")
	//fmt.Println(string(org))
	//ReverseBytes(org)
	//fmt.Println(string(org))
	strOrg := "hello jack, how are you!"
	byte58 := Base58Encode([]byte(strOrg))
	str58 := string(byte58)
	fmt.Printf("%s\n", string(str58))
	fmt.Printf("%s", string(Base58Decode([]byte(str58))))
}
