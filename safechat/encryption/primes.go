package encryption

import (
	"math/rand"
	"time"
)

func generatePrimes(lower, upper uint64) (uint64, uint64) {
	p := nextPrime(rand.Uint64()%(upper-lower) + lower)
	q := nextPrime(rand.Uint64()%(upper-lower) + lower)
	if p == q {
		q = nextPrime(q)
	}
	return p, q
}

func generateKeys(p, q uint64) (PrivateKey, PublicKey) {
	n := p * q
	e := uint64(3)
	for {
		if gcd(p-1, e) == 1 && gcd(q-1, e) == 1 {
			break
		}
		e += 2
	}
	d := modularInverse(e, phi(n))
	return PrivateKey{n, d}, PublicKey{n, e}
}

func GenerateKeyPair() (PublicKey, PrivateKey) {
	rand.Seed(time.Now().Unix())
	bound := uint64(1 << 15)
	p, q := generatePrimes(bound, bound*2)
	priv, pub := generateKeys(p, q)
	return pub, priv
}
