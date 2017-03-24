package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"validator"
	"controller"
	"log"
	"utils"
	"net/http"
)

func main() {
	var err error
	validator.ValidateConf()

	controller.DB,err= sql.Open(validator.NewConfig().GetDriver(), validator.NewConfig().GetDataSource())
	if nil != err {
		log.Fatal(err)
	}
	_,err = controller.DB.Exec("CREATE DATABASE IF NOT EXISTS "+utils.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	_,err = controller.DB.Exec("USE "+utils.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}

	_,err = controller.DB.Exec("CREATE TABLE IF NOT EXISTS  loginuser ( username varchar(32), password varchar(128) )")
	if err != nil {
		log.Fatal(err)
	}
	controller.DB.Close()
	controller.DB, err = sql.Open(validator.NewConfig().GetDriver(), validator.NewConfig().GetDataSource() + utils.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer controller.DB.Close()
	err = controller.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server Listening on ",validator.NewConfig().GetServerADDR())
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/login", controller.LoginPage)
	http.HandleFunc("/unregister", controller.DeleteUser)
	http.HandleFunc("/", controller.HomePage)
	http.ListenAndServe(validator.NewConfig().GetServerADDR(), nil)
}
