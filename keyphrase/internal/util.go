package internal

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func LoadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

func SaveWordsToFile(words []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, w := range words {
		_, err := file.WriteString(w + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func BitLength(n int) int {
	bits := 0
	for (1 << bits) < n {
		bits++
	}
	return bits
}

func BytesToWords(data []byte, wordlist []string) ([]string, error) {
	bitsPerWord := BitLength(len(wordlist))
	totalBits := len(data) * 8

	if totalBits%bitsPerWord != 0 {
		return nil, fmt.Errorf("wordlist size %d incompatible: %d bits not divisible by %d bits/word",
			len(wordlist), totalBits, bitsPerWord)
	}

	num := new(big.Int).SetBytes(data)
	binStr := fmt.Sprintf("%0*b", totalBits, num)

	numWords := totalBits / bitsPerWord
	words := make([]string, 0, numWords)

	for i := 0; i < len(binStr); i += bitsPerWord {
		chunk := binStr[i : i+bitsPerWord]
		idx := new(big.Int)
		idx.SetString(chunk, 2)
		words = append(words, wordlist[idx.Int64()])
	}

	return words, nil
}

func WordsToBytes(words []string, wordlist []string, expectedLen int) ([]byte, error) {
	wordMap := make(map[string]int)
	for i, w := range wordlist {
		wordMap[w] = i
	}

	bitsPerWord := BitLength(len(wordlist))
	totalBits := len(words) * bitsPerWord

	if totalBits != expectedLen*8 {
		return nil, fmt.Errorf("invalid number of words: got %d, expected length %d bytes", len(words), expectedLen)
	}

	var binStr strings.Builder
	for i, w := range words {
		idx, ok := wordMap[w]
		if !ok {
			return nil, fmt.Errorf("word #%d not in wordlist: %q", i+1, w)
		}
		binStr.WriteString(fmt.Sprintf("%0*b", bitsPerWord, idx))
	}

	num := new(big.Int)
	num.SetString(binStr.String(), 2)

	data := num.Bytes()
	if len(data) < expectedLen {
		padded := make([]byte, expectedLen)
		copy(padded[expectedLen-len(data):], data)
		data = padded
	}

	return data, nil
}

func ChecksumSHA256(data []byte, length int) []byte {
	sum := sha256.Sum256(data)
	return sum[:length]
}

func EqualBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
