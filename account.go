package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io"
	"sync"

	"github.com/btcsuite/btcd/btcec"
)

type OgAccount struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func NewAccount(privHex string) (*OgAccount, error) {
	priv, err := hex.DecodeString(privHex)
	if err != nil {
		return nil, fmt.Errorf("decode hex private error: %v", err)
	}

	_, ecdsapub := btcec.PrivKeyFromBytes(btcec.S256(), priv)
	pub := FromECDSAPub((*ecdsa.PublicKey)(ecdsapub))

	addr := Keccak256(pub)[12:]

	a := OgAccount{}
	a.PrivateKey = fmt.Sprintf("%x", priv)
	a.PublicKey = fmt.Sprintf("%x", pub)
	a.Address = fmt.Sprintf("%x", addr)

	return &a, nil
}

func GenerateAccount() OgAccount {
	priv, pub := randomKeyPair()

	a := OgAccount{}
	a.PrivateKey = fmt.Sprintf("%x", priv)
	a.PublicKey = fmt.Sprintf("%x", pub)
	a.Address = fmt.Sprintf("%x", Keccak256(pub)[12:])

	return a
}

func randomKeyPair() (priv, pub []byte) {
	privBytes := [32]byte{}
	copy(privBytes[:], CRandBytes(32))

	priv = privBytes[:]

	_, ecdsapub := btcec.PrivKeyFromBytes(btcec.S256(), priv)
	pub = FromECDSAPub((*ecdsa.PublicKey)(ecdsapub))

	return priv, pub
}

func CRandBytes(numBytes int) []byte {
	gRandInfo := &randInfo{}
	gRandInfo.MixEntropy(RandBytes(32))

	b := make([]byte, numBytes)
	_, err := gRandInfo.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

type randInfo struct {
	mtx          sync.Mutex
	seedBytes    [32]byte
	cipherAES256 cipher.Block
	streamAES256 cipher.Stream
	reader       io.Reader
}

func (ri *randInfo) MixEntropy(seedBytes []byte) {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()
	// Make new ri.seedBytes
	hashBytes := Sha256(seedBytes)
	hashBytes32 := [32]byte{}
	copy(hashBytes32[:], hashBytes)
	ri.seedBytes = xorBytes32(ri.seedBytes, hashBytes32)
	// Create new cipher.Block
	var err error
	ri.cipherAES256, err = aes.NewCipher(ri.seedBytes[:])
	if err != nil {
		panic("Error creating AES256 cipher: " + err.Error())
	}
	// Create new stream
	ri.streamAES256 = cipher.NewCTR(ri.cipherAES256, RandBytes(aes.BlockSize))
	// Create new reader
	ri.reader = &cipher.StreamReader{S: ri.streamAES256, R: crand.Reader}
}

func (ri *randInfo) Read(b []byte) (n int, err error) {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()
	return ri.reader.Read(b)
}

func xorBytes32(bytesA [32]byte, bytesB [32]byte) (res [32]byte) {
	for i, b := range bytesA {
		res[i] = b ^ bytesB[i]
	}
	return res
}

func Sha256(bytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

func RandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return btcec.S256()
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
