package encryption

func phi(n uint64) uint64 {
	result := n
	for i := uint64(2); i*i <= n; i++ {
		if n%i == 0 {
			for n%i == 0 {
				n /= i
			}
			result -= result / i
		}
	}
	if n > 1 {
		result -= result / n
	}
	return result
}

func pow(x, y, m uint64) uint64 {
	if y == 0 {
		return 1
	}
	p := pow(x, y/2, m)
	p = (p * p) % m
	if y%2 == 0 {
		return p
	} else {
		return (x * p) % m
	}
}

func modularInverse(a, m uint64) uint64 {
	return pow(a, phi(m)-1, m)
}

func isPrime(x uint64) bool {
	if x < 2 {
		return false
	}
	if x == 2 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	for i := uint64(3); i*i <= x; i += 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(x uint64) uint64 {
	for {
		x++
		if isPrime(x) {
			return x
		}
	}
}

func gcd(a, b uint64) uint64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
