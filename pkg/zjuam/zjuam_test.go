package zjuam

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestEncryption(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	modulus := "9d2c7ce4add43b2caf5a6f49b2fef2ef38c008c0a8132133b85da9f577c39a49dd87c076c7087ab2fbc8c810bc8be5730ab81c18b87ef87004ef1fb8fc628c23"
	exponent := "10001"
	password := "aaaa0123"

	pk, err := newPubKey(modulus, exponent)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	p := pk.encrypt(password)

	if p != "3f5312e265def9ed1da2a5d3c5dcf5d6e0af31d1e8561af8d5c01326b66bb59fc1c85632fe9d8d70932d71a425b58f2e5e0c21bd9885d7feea383e6b21a974bf" {
		t.Fail()
	}
}
