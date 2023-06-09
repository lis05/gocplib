type factorsEngine struct {
	n      int
	primes []int
	sieve  []int
}

func (f *factorsEngine) init(n int) {
	f.n = n
	f.primes = make([]int, 0)
	f.sieve = make([]int, n+1)

	f.sieve[1] = 1
	for i := 2; i <= n; i++ {
		if f.sieve[i] != 0 {
			continue
		}
		f.primes = append(f.primes, i)
		f.sieve[i] = i
		for ii := i * i; ii <= n; ii += i {
			f.sieve[ii] = i
		}
	}
}

func (f factorsEngine) isPrime(n int) bool {
	if n <= f.n {
		return f.sieve[n] == n
	} else {
		return isPrime(n)
	}
}

func (f factorsEngine) factorize(n int) []pair[int] {
	var res []pair[int]
	for div := 2; div*div <= n; div++ {
		if n <= f.n {
			break
		}

		cnt := 0
		for n%div == 0 {
			cnt++
			n /= div
		}
		if cnt == 0 {
			continue
		}
		res = append(res, pair[int]{div, cnt})
	}

	if n > f.n {
		res = append(res, pair[int]{n, 1})
		return res
	}

	for n > 1 {
		cnt := 0
		div := f.sieve[n]
		for n%div == 0 {
			cnt++
			n /= div
		}
		if cnt == 0 {
			continue
		}
		res = append(res, pair[int]{div, cnt})
	}

	return res
}
