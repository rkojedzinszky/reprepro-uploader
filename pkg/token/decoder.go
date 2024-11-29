package token

import (
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// DecoderOpts represents options applied to each token
type DecoderOpts func(*jwt.Expected)

// Decoder decodes a token into DataWithClaims
type Decoder interface {
	Decode(string, DataWithClaims, ...DecoderOpts) error
}

// DecoderWithIssuer will expect Issuer to match
func DecoderWithIssuer(issuer string) DecoderOpts {
	return func(e *jwt.Expected) {
		e.Issuer = issuer
	}
}

// DecoderWithSubject will expect Subject to match
func DecoderWithSubject(subject string) DecoderOpts {
	return func(e *jwt.Expected) {
		e.Subject = subject
	}
}

// DecoderWithTime will expect claim to be valid in current time
func DecoderWithTime() DecoderOpts {
	return func(e *jwt.Expected) {
		e.Time = time.Now()
	}
}

type jwedecoder struct {
	key interface{}

	opts []DecoderOpts
}

// NewDecoder returns a new JWE decoder
func NewDecoder(key []byte, opts ...DecoderOpts) (Decoder, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return nil, ErrInvalidKey
	}

	return jwedecoder{
		key:  key,
		opts: opts,
	}, nil
}

// MustNewDecoder returns a new JWE decoder, panics on error
func MustNewDecoder(key []byte, opts ...DecoderOpts) Decoder {
	enc, err := NewDecoder(key, opts...)

	if err != nil {
		panic(err)
	}

	return enc
}

// Decode validates the token and decodes into data
// data should only be trusted if returned error is nil
func (j jwedecoder) Decode(token string, data DataWithClaims, opts ...DecoderOpts) error {
	object, err := jwt.ParseEncrypted(token, []jose.KeyAlgorithm{jose.A256GCMKW, jose.DIRECT}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return err
	}

	if err := object.Claims(j.key, data); err != nil {
		return err
	}

	expected := jwt.Expected{}

	for _, opt := range append(j.opts, opts...) {
		opt(&expected)
	}

	return data.GetClaims().Validate(expected)
}
