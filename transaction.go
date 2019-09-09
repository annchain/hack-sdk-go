package hackSDK

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/annchain/OG/common/crypto/secp256k1"
)

var (
	EmptyAddress = "0000000000000000000000000000000000000000"
	EmptyHash    = "0000000000000000000000000000000000000000000000000000000000000000"
	EmptyUint64  = uint64(0)
	EmptyBigInt  = big.NewInt(0)
)

type Transaction struct {
	Parents   []string
	From      string
	To        string
	Nonce     uint64
	Guarantee *big.Int
	Value     *big.Int
}

func (tx *Transaction) SignatureTarget() ([]byte, error) {
	msg := &bytes.Buffer{}

	// write parents
	for _, parentHex := range tx.Parents {
		pBytes, err := HexToBytes(parentHex)
		if err != nil {
			return nil, fmt.Errorf("invalid parent: %v", err)
		}
		binary.Write(msg, binary.BigEndian, pBytes)
	}

	// write nonce
	binary.Write(msg, binary.BigEndian, tx.Nonce)

	// write from, to
	fromBytes, err := HexToBytes(tx.From)
	if err != nil {
		return nil, fmt.Errorf("invalid FROM: %v", err)
	}
	binary.Write(msg, binary.BigEndian, fromBytes)

	if tx.To == "" {
		tx.To = EmptyAddress
	}
	toBytes, err := HexToBytes(tx.To)
	if err != nil {
		return nil, fmt.Errorf("invalid TO: %v", err)
	}
	binary.Write(msg, binary.BigEndian, toBytes)

	// write value
	value := tx.Value.Bytes()
	if tx.Value.Int64() == 0 {
		value = []byte{0}
	}
	binary.Write(msg, binary.BigEndian, value)

	// write guarantee
	guarantee := tx.Guarantee.Bytes()
	if tx.Guarantee.Int64() == 0 {
		guarantee = []byte{0}
	}
	binary.Write(msg, binary.BigEndian, guarantee)

	return msg.Bytes(), nil
}

func (tx *Transaction) Sign(privBytes []byte) ([]byte, error) {
	priv, err := toECDSA(privBytes, true)
	if err != nil {
		return nil, fmt.Errorf("ToECDSA error: %v. priv bytes: %x", err, privBytes)
	}
	msg, err := tx.SignatureTarget()
	if err != nil {
		return nil, fmt.Errorf("get signature target error: %v", err)
	}

	hash := Sha256(msg)
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	seckey := PaddedBigBytes(priv.D, priv.Params().BitSize/8)
	return secp256k1.Sign(hash, seckey)
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)
