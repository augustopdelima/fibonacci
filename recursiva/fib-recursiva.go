package recursiva

func FibonacciRecursiva(numero int) int {
	if numero < 2 {
		return numero
	}
	return FibonacciRecursiva(numero-1) + FibonacciRecursiva(numero-2)
}
