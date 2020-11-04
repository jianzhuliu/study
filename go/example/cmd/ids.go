package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	fmt.Println("user id:", os.Getuid())

	var u *user.User
	u, _ = user.Current()
	fmt.Println("Group ids:")
	groupIds, _ := u.GroupIds()
	for _, i := range groupIds {
		fmt.Println(i, " ")
	}
	fmt.Println()
	
	
	var c chan int = make(chan int)
	num := 5
	
	go get(c)
	
	c <- num
	
	fmt.Println("main")
	fmt.Println("sum",<-c)
}

func get(c chan int){
	fmt.Println("c")
	
	d := <- c 
	sum := 0
	for i:=0;i<=d;i++{
		sum += i 
	}
	
	fmt.Println("c send")
	c <- sum
	
	
}