package main

import (
	"Lumino/config"
	"Lumino/model"
	"Lumino/store"
	"flag"
	"fmt"
)

func main() {
	config.LoadConfig()
	configName := flag.String("dbConfigName", "", "database config name")
	shouldDrop := flag.Bool("shouldDrop", false, "drop tables if existing before migration")
	flag.Parse()
	fmt.Println("configName: ", *configName)
	fmt.Println("shouldDrop: ", *shouldDrop)
	var db *store.DB
	if *configName == "" {
		db = store.NewPgDB()
	} else {
		db = store.NewPgDBWithConfig(*configName)
	}
	if err := db.AutoMigrate(model.Transaction{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(model.User{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(model.Category{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(model.Account{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(model.AccountBook{}); err != nil {
		fmt.Println(err)
	}
}
