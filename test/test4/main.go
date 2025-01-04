package main

func main() {
	var n int;
	n = 5;
	var res int;
	res = 1;
	for n >= 1 {
		res = res * n;
		n = n - 1;
	}

	for {
		println(res);
		break;
	}

	println(factorial(5));

	return;
}

func factorial(n int) int {
	if n <= 1 {
		return 1;
	}

	return n * factorial(n - 1);
}
