package topdown

func FibonacciTopDown(numero int) int {

	casosBases := map[int]int{
		0: 0,
		1: 1,
	}

	return computarCache(numero, casosBases)
}

func computarCache(numero int, cache map[int]int) int {
	if valor, encontrado := cache[numero]; encontrado {
		return valor
	}

	cache[numero] = computarCache(numero-1, cache) + computarCache(numero-2, cache)
	return cache[numero]
}
