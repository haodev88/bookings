package main

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

type Person struct {
	name string
	job  string
}

type Test struct {
	firstName string
	lastName  string
}

func main() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	rootPath := filepath.Dir(d)
	fileEnv := rootPath + "\\.env"
	log.Println(fileEnv)
	err := godotenv.Load(fileEnv)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")
	log.Println(s3Bucket, secretKey)
	os.Exit(1)
	/*
		myArr:=make(map[string]string)
		myArr["name"] = "haotran"
		log.Print(myArr["name"])
	*/

	names := []string{"Tran", "Vi", "Hao"}
	log.Println(names)

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

	printData := []Person{
		{name: "Hao tran", job: "develop"},
		{name: "Dieu Truong", job: "Designer"},
	}

	for i, item := range printData {
		log.Println(i, item)
	}

	p1 := Person{
		name: "Haotran",
		job:  "Developer",
	}

	p2 := Person{
		name: "Dieu Truong",
		job:  "Designer",
	}

	var pArr []Person
	pArr = append(pArr, p1)
	pArr = append(pArr, p2)

	for _, item := range pArr {
		val := reflect.Indirect(reflect.ValueOf(item))
		keyName := val.Type().Field(0).Name
		keyJob := val.Type().Field(1).Name
		log.Println(keyName, item.name)
		log.Println(keyJob, item.job)
	}

	//data:= make(map[string][]string)
	//data["test"] = []string{"one", "two"}
	//fmt.Println(data["test"][0])

	data := make(map[string][][]Test)
	data["test"] = append(data["test"], []Test{
		{"Hao", "Tran"},
		{"Dieu", "Truong"},
	})

	/*
		var data [][]Person
			data = append(data, []Person{
				{"Hao", "Tran"},
				{"Dieu", "Truong"},
			})
			data = append(data, []Person{
				{"Hao1", "Tran1"},
				{"Dieu1", "Truong1"},
			})
			fmt.Println(data)
	*/

}
