package main

import (
	"log"
	"myproto/protopackage"
	"github.com/golang/protobuf/proto"
)

func showErr(msg string, err error){
	log.Fatal(msg, err)
}

func main(){
	test := &protopackage.Student{
		Name : "student1",
		Male : true,
		Scores : []int32{98,99,87},
	}
	
	data, err := proto.Marshal(test)
	if err != nil {
		showErr("proto.Marshal()",err)
	}
	
	newTest := &protopackage.Student{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		showErr("proto.Unmarshal()", err)
	}
	
	log.Println(test.GetName())
	if test.GetName() != newTest.GetName(){
		log.Fatalf("data mismatch %q != %q", test.GetName(), newTest.GetName())
	}
}