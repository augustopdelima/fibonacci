package bottomup

func FibonacciBottomUp(numero int) int {
	if numero <= 1 {
		return numero
	}

	cache := make(map[int]int)

	cache[0] = 0
	cache[1] = 1

	for i := 2; i <= numero; i++ {
		ultimoIndice := i - 1
		penultimoIndice := i - 2

		cache[i] = cache[ultimoIndice] + cache[penultimoIndice]
	}

	return cache[numero]
}
