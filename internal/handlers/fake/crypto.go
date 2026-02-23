package fake

import (
	stdcrypto "crypto"
	"crypto/ed25519"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

type CryptoHandler interface {
	RSACache() *cache.RSACache
	ECDSACache() *cache.ECDSACache
	ED25519Cache() *cache.ED25519Cache
}

func LoadHandlerKey(h CryptoHandler, meta *CryptoMeta) (pub stdcrypto.PublicKey, priv stdcrypto.PrivateKey, err error) {
	switch meta.Algorithm {
	case crypto.AlgorithmECDSA:
		if key, err2 := h.ECDSACache().Load(meta.Subject, meta.ECDSACurve); err2 != nil {
			err = err2
		} else {
			pub = &key.PublicKey
			priv = key
		}
	case crypto.AlgorithmED25519:
		if key, err2 := h.ED25519Cache().Load(meta.Subject); err2 != nil {
			err = err2
		} else {
			pub = key.Public().(ed25519.PublicKey)
			priv = key
		}
	default:
		if key, err2 := h.RSACache().Load(meta.Subject, meta.Length); err2 != nil {
			err = err2
		} else {
			pub = &key.PublicKey
			priv = key
		}
	}

	return
}
