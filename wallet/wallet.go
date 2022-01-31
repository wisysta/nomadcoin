package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.io/wisysta/nomadcoin/utils"
)

// 공개키 암호화의 비밀키는 '암호문을 복호화'하기 위한 것이고, 공개키는 '평문을 암호화'하는 것입니다.
// 전자서명의 비밀키는 '평문을 사용하여 서명'하기 위한 것이고, 공개키는 '서명이 올바른지 검증'하는 것입니다. 즉 역할이 전혀 다릅니다.

// 1) hash msg
// 2) generate key pair -> keyPair (privateK, publicK) (save priv to a file)
// 3) sign the hash -> ("hashed_msg" + privateKey) => "signature"
// 4) verify -> ("hashed_msg" + "signature" + "publicK") -> true

// const (
// 	signature     string = "8abc60763cb92df0a9edc7c46817dd2eb94b848d4b41d36154b6ebb5a7a5decad5cb8b755af96336d2e8ec391c311f2cdf283a44e9d3c7ae6811d3aab1d09f09"
// 	privateKey    string = "30770201010420c6fe4bd4a221124e9a0a9d84670aebff1987c9ba4eb5cc2d2df473c5e3b3805ea00a06082a8648ce3d030107a1440342000434ad4365b49ab702d4547a3a1e73f0e343e5628a3c52023e8f1cdb05d71dad58865cd5100835437db6e5aa3adc9fab025b1570df259052a6cccef6d91eb37778"
// 	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
// )

// func Start() {
// 	privBytes, err := hex.DecodeString(privateKey)
// 	utils.HandleErr(err)

// 	private, err := x509.ParseECPrivateKey(privBytes)
// 	utils.HandleErr(err)

// 	sigBytes, err := hex.DecodeString(signature)
// 	rBytes := sigBytes[:len(sigBytes)/2]
// 	sBytes := sigBytes[len(sigBytes)/2:]

// 	var bigR, bigS = big.Int{}, big.Int{}

// 	bigR.SetBytes(rBytes)
// 	bigS.SetBytes(sBytes)

// 	hashBytes, err := hex.DecodeString(hashedMessage)

// 	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)

// 	fmt.Println(ok)
// }

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat("nomadcoin.wallet")
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	sigBytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := sigBytes[:len(sigBytes)/2]
	secondHalfBytes := sigBytes[len(sigBytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)

	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
