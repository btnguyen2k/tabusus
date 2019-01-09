package main
//
//import (
//	"fmt"
//	"github.com/mongodb/mongo-go-driver/bson"
//	"tabusus"
//)
//
//func main() {
//	app := tabusus.NewApp("system")
//	fmt.Println("Application:", app)
//	data, err := app.ToJson()
//	fmt.Println("Error:", err)
//	fmt.Println("JSON :", string(data[:]))
//	var json bson.M
//	bson.UnmarshalExtJSON(data, false, &json)
//	fmt.Println("Data :", json)
//	fmt.Println("Created:", app.GetTimeCreated())
//	fmt.Println("Updated:", app.GetTimeUpdated())
//
//	fmt.Println()
//
//	app2 := tabusus.NewAppFromJson(json)
//	fmt.Println("Application:", app2)
//	data, err = app2.ToJson()
//	fmt.Println("Error:", err)
//	fmt.Println("JSON :", string(data[:]))
//	bson.UnmarshalExtJSON(data, false, &json)
//	fmt.Println("Data :", json)
//	fmt.Println("Created:", app2.GetTimeCreated())
//	fmt.Println("Updated:", app2.GetTimeUpdated())
//}
