package claims

import "github.com/rkojedzinszky/reprepro-uploader/pkg/token"

type Claims struct {
	token.EmbedJWTClaims

	Distributions []string `json:"distributions,omitempty"`
}
