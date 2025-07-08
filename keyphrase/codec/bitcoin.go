package codec

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/restartfu/keyphrase/keyphrase/internal"
	"log"
)

const btcPrivKeyLen = 32
const btcChecksumLen = 4

type bitcoin struct{}

func (bitcoin) Encode(hexKey string, wordlist []string) ([]string, error) {
	privKey, err := btcutil.DecodeWIF(hexKey)
	if err != nil {
		log.Fatalf("Invalid WIF: %v", err)
	}
	rawKey := privKey.PrivKey.Serialize()

	if len(rawKey) != btcPrivKeyLen {
		return nil, fmt.Errorf("bitcoin private key must be %d bytes", btcPrivKeyLen)
	}
	checksum := internal.ChecksumSHA256(rawKey, btcChecksumLen)
	data := append(rawKey, checksum...)

	return internal.BytesToWords(data, wordlist)
}

func (bitcoin) Decode(words []string, wordlist []string) (string, error) {
	expectedLen := btcPrivKeyLen + btcChecksumLen
	data, err := internal.WordsToBytes(words, wordlist, expectedLen)
	if err != nil {
		return "", err
	}

	key := data[:btcPrivKeyLen]
	check := data[btcPrivKeyLen:]

	calcChecksum := internal.ChecksumSHA256(key, btcChecksumLen)
	if !internal.EqualBytes(check, calcChecksum) {
		return "", fmt.Errorf("checksum mismatch - wordlist or word order may be wrong")
	}

	privKey, _ := btcec.PrivKeyFromBytes(key)
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return "", fmt.Errorf("failed to create WIF: %v", err)
	}

	return wif.String(), nil
}
