package fake

import (
	stdcrypto "crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"io"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

type CryptoHandler interface {
	RSACache() cache.Cacher[*rsa.PrivateKey]
	ECDSACache() cache.Cacher[*ecdsa.PrivateKey]
	ED25519Cache() cache.Cacher[ed25519.PrivateKey]
}

func LoadHandlerKey(h CryptoHandler, meta *CryptoMeta, rnd io.Reader) (pub stdcrypto.PublicKey, priv stdcrypto.PrivateKey, err error) {
	switch meta.Algorithm {
	case crypto.AlgorithmECDSA:
		req := &cache.ECDSALoader{
			Hostname: meta.Subject,
			Curve:    meta.ECDSACurve,
			Random:   rnd,
		}
		if key, err2 := h.ECDSACache().Load(req); err2 != nil {
			err = err2
		} else {
			pub = &key.PublicKey
			priv = key
		}
	case crypto.AlgorithmED25519:
		req := &cache.ED25519Loader{
			Hostname: meta.Subject,
			Random:   rnd,
		}
		if key, err2 := h.ED25519Cache().Load(req); err2 != nil {
			err = err2
		} else {
			pub = key.Public().(ed25519.PublicKey)
			priv = key
		}
	default:
		req := &cache.RSALoader{
			Hostname: meta.Subject,
			Length:   meta.Length,
			Random:   rnd,
		}
		if key, err2 := h.RSACache().Load(req); err2 != nil {
			err = err2
		} else {
			pub = &key.PublicKey
			priv = key
		}
	}

	return
}
