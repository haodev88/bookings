package main

import (
	"log"
	"reflect"
)

type Person struct {
	name string
	job string
}

func main()  {
	/*
	myArr:=make(map[string]string)
	myArr["name"] = "haotran"
	log.Print(myArr["name"])
	 */

	names:=[]string{"Tran", "Vi", "Hao"}
	log.Println(names)
	log.Print(names)

	/*
	person := Person{
		"Haotran",
		"Developer",
	}
	var myArr map[string]Person
	myArr = make(map[string]Person)
	myArr["me"] = person
	log.Print(myArr)
	*/

	/*
	person:= Person{
		"haotran",
		"developer",
	}
	var printArr map[string]Person
	printArr = make(map[string]Person)
	printArr["me"] = person
	log.Println(printArr)
	*/

	/*
	data:= make(map[string]string)
	data["name"] = "haotran"
	data["job"]  = "developer"

	printData:=make(map[string]interface{})
	printData["me"] = data
	log.Print(printData)
	 */

	printData:= []Person{
		{name:"Hao tran", job:"develop"},
		{name:"Dieu Truong", job:"Designer"},
	}


	for i,item:= range printData {
		log.Println(i, item)
	}

	p1 := Person{
		name:"Haotran",
		job:"Developer",
	}

	p2:= Person{
		name:"Dieu Truong",
		job:"Designer",
	}

	var pArr []Person
	pArr = append(pArr, p1)
	pArr = append(pArr, p2)

	for _,item:= range pArr {
		val := reflect.Indirect(reflect.ValueOf(item))
		keyName:= val.Type().Field(0).Name
		keyJob:= val.Type().Field(1).Name
		log.Println(keyName, item.name)
		log.Println(keyJob, item.job)
	}
}

