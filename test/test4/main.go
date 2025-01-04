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
}
