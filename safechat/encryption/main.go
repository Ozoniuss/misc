package encryption

// import (
// 	"math/rand"
// )

// func generatePrimes(lower, upper uint64) (uint64, uint64) {
// 	p := nextPrime(rand.Uint64()%(upper-lower) + lower)
// 	q := nextPrime(rand.Uint64()%(upper-lower) + lower)
// 	if p == q {
// 		q = nextPrime(q)
// 	}
// 	return p, q
// }

// func generateKeys(p, q uint64) (PrivateKey, PublicKey) {
// 	n := p * q
// 	e := uint64(3)
// 	for {
// 		if gcd(p-1, e) == 1 && gcd(q-1, e) == 1 {
// 			break
// 		}
// 		e += 2
// 	}
// 	d := modularInverse(e, phi(n))
// 	return PrivateKey{n, d}, PublicKey{n, e}
// }

// func main() {
// 	rand.Seed(time.Now().Unix())
// 	bound := uint64(1 << 15)
// 	p, q := generatePrimes(bound, bound*2)
// 	priv, pub := generateKeys(p, q)
// 	s := []byte("babuinul unor coaie")
// 	c := pub.EncryptString(s)
// 	m2 := priv.DecryptString(c)
// 	fmt.Printf("Private key: %s\n", priv.String())
// 	fmt.Printf("Public key: %s\n", pub.String())
// 	fmt.Printf("Encryted message: %s\n", c)
// 	fmt.Printf("Decrypted message: %s\n", m2)
// }
