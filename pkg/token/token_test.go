package token_test

import (
	"math/rand"
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/rkojedzinszky/reprepro-uploader/pkg/token"
)

func genKey() (key []byte) {
	key = make([]byte, 32)

	rand.Read(key)

	return
}

type testtoken struct {
	token.EmbedJWTClaims

	Userid   int
	Username string
}

func TestData(t *testing.T) {
	key := genKey()

	enc := token.MustNewEncoder(key)
	te := testtoken{Userid: 5, Username: "user"}
	tokenstring, err := enc.Encode(&te)
	if err != nil {
		t.Fatal(err)
	}

	dec := token.MustNewDecoder(key)
	td := testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != nil {
		t.Error(err)
	}

	if td.Userid != te.Userid {
		t.Error("Expected decoded Userid to match")
	}

	if td.Username != te.Username {
		t.Error("Expected decoded Username to match")
	}
}

func TestExpiry(t *testing.T) {
	key := genKey()

	// Have to go backwards enough to fail VaidateWithLeeway
	enc := token.MustNewEncoder(key, token.EncoderWithAge(-jwt.DefaultLeeway+time.Second))
	te := testtoken{}
	tokenstring, err := enc.Encode(&te)
	if err != nil {
		t.Fatal(err)
	}

	dec := token.MustNewDecoder(key, token.DecoderWithTime())
	td := testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	td = testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != jwt.ErrExpired {
		t.Error("expected jwt.ErrExpired")
	}
}

func TestIssuer(t *testing.T) {
	key := genKey()

	// Have to go backwards enough to fail VaidateWithLeeway
	enc := token.MustNewEncoder(key, token.EncoderWithIssuer("test"))
	te := testtoken{}
	tokenstring, err := enc.Encode(&te)
	if err != nil {
		t.Fatal(err)
	}

	dec := token.MustNewDecoder(key, token.DecoderWithIssuer("test"))
	td := testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != nil {
		t.Error(err)
	}

	dec = token.MustNewDecoder(key, token.DecoderWithIssuer("test2"))
	td = testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != jwt.ErrInvalidIssuer {
		t.Error("expected jwt.ErrInvalidIssuer")
	}
}

func TestSubject(t *testing.T) {
	key := genKey()

	// Have to go backwards enough to fail VaidateWithLeeway
	enc := token.MustNewEncoder(key, token.EncoderWithSubject("test"))
	te := testtoken{}
	tokenstring, err := enc.Encode(&te)
	if err != nil {
		t.Fatal(err)
	}

	dec := token.MustNewDecoder(key, token.DecoderWithSubject("test"))
	td := testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != nil {
		t.Error(err)
	}

	dec = token.MustNewDecoder(key, token.DecoderWithSubject("test2"))
	td = testtoken{}
	err = dec.Decode(tokenstring, &td)
	if err != jwt.ErrInvalidSubject {
		t.Error("expected jwt.ErrInvalidSubject")
	}
}
