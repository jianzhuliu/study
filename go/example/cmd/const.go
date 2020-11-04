/*
在同一个 const group 中，如果常量定义与前一行的定义一致，则可以省略类型和值。
编译时，会按照前一行的定义自动补全。
*/
package main

import "fmt"

func main() {
	const (
		a, b = "golang", 100
		d, e
		f bool = true
		g
	)
	fmt.Println(d, e, g)
}