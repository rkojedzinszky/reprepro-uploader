package token

import (
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// EncoderOpts represents options applied to each token
type EncoderOpts func(*jwt.Claims)

// EncoderWithIssuer will set Issuer
func EncoderWithIssuer(issuer string) EncoderOpts {
	return func(c *jwt.Claims) {
		c.Issuer = issuer
	}
}

// EncoderWithSubject will set Issuer
func EncoderWithSubject(subject string) EncoderOpts {
	return func(c *jwt.Claims) {
		c.Subject = subject
	}
}

// EncoderWithAge will set expiry according to current time and given duration
func EncoderWithAge(age time.Duration) EncoderOpts {
	return func(c *jwt.Claims) {
		c.Expiry = jwt.NewNumericDate(time.Now().Add(age))
	}
}

// Encoder encodes DataWithClaims into a token
type Encoder interface {
	Encode(DataWithClaims, ...EncoderOpts) (string, error)
}

type jweencoder struct {
	builder jwt.Builder

	opts []EncoderOpts
}

// NewEncoder returns a JWE encoder
// expects a 16, 24 or 32 byte length key
func NewEncoder(key []byte, opts ...EncoderOpts) (Encoder, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return nil, ErrInvalidKey
	}

	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       key,
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	return jweencoder{
		builder: jwt.Encrypted(encrypter),
		opts:    opts,
	}, nil
}

// MustNewEncoder returns a JWE encoder
// panics on wrong key length
func MustNewEncoder(key []byte, opts ...EncoderOpts) Encoder {
	enc, err := NewEncoder(key, opts...)

	if err != nil {
		panic(err)
	}

	return enc
}

// Encode encodes data into token
func (j jweencoder) Encode(data DataWithClaims, opts ...EncoderOpts) (token string, err error) {
	claims := data.GetClaims()

	for _, opt := range append(j.opts, opts...) {
		opt(claims)
	}

	token, err = j.builder.Claims(data).Serialize()

	return
}
