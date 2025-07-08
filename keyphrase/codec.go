package keyphrase

type Codec interface {
	Encode(hexKey string, wordlist []string) ([]string, error)
	Decode(words []string, wordlist []string) (string, error)
}
