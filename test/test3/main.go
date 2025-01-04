package main

func main() {
	var a int;

	if true {
		a = 1;
		var b int;
		println("true");
	} else {
		println("false");
	}

	if false {
		println("false x2");
	} else {
		println("true x2");
	}

	if false {
		println("false x3");
	} else if true {
		println("true x3");
	} else {
		println("false x3");
	}

	if false {
		println("false x4");
	} else if false {
		println("false x4");
	} else {
		println("true x4");
	}

	println(a);
	println(b);
}
