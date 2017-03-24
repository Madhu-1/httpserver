package controller

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"log"
	"validator"
	"utils"
)
//DB global database variable
var(
	DB *sql.DB

)
func errorHandler(resp http.ResponseWriter,message string,code int){
	http.Error(resp,message,code)
	return
}
//Register func for user signup , user information is validator and stored in database
func Register(res http.ResponseWriter, req *http.Request) {
	panicRecover(res)
	if req.Method != "POST" {
		http.ServeFile(res, req, "views/signup.html")
		return
	}
	var(
		user string
		username string
		password string
	)

	username = req.FormValue("username")
	password = req.FormValue("password")
	err:=validator.ValidateUserInfo(username,password)
	if nil != err{
		errorHandler(res,err.Error(),utils.InvalidRespCode)
		return

	}

	err = DB.QueryRow("SELECT username FROM loginuser WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("generate password error",err.Error())
			errorHandler(res, "internal server error, unable to create your account.", utils.ServerErrRespCode)
			return
		}
		stmt, err := DB.Prepare("INSERT INTO loginuser(username,password) VALUES(?,?)")
		if err != nil {
			log.Println("database error ",err.Error())
			errorHandler(res, "Server error, unable to create your account.",  utils.ServerErrRespCode)
			return
		}
		_, err = stmt.Exec(username, hashedPassword)
		if err != nil {
			log.Println(err.Error())
			errorHandler(res, "Server error, unable to create your account.",  utils.ServerErrRespCode)
			return
		}
		res.WriteHeader(200)
		res.Write([]byte("User created!"))

	case err != nil:
		log.Println("database error",err.Error())
		errorHandler(res, "internal server error, unable to create your account.", utils.ServerErrRespCode)

	default:
		http.Redirect(res, req, "/", 301)
	}
}
// LoginPage  provided information is  authenticated
func LoginPage(res http.ResponseWriter, req *http.Request) {
	panicRecover(res)
	if req.Method != "POST" {
		http.ServeFile(res, req, "views/login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	var databaseUsername,databasePassword string


	err := DB.QueryRow("SELECT username, password FROM loginuser WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	switch {
	case err==sql.ErrNoRows:
		errorHandler(res,fmt.Sprintf("username %s doesnot exist",username),404)
		return
	case nil!=err:
		log.Println(err.Error())
		errorHandler(res,utils.InternalServerErrorMsg,utils.ServerErrRespCode)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if nil != err {
		errorHandler(res,"password doesnot match",utils.InvalidRespCode)
		return
	}
	res.WriteHeader(200)
	res.Write([]byte("Login successful \n Hello " + databaseUsername))

}
// DeleteUser verify user information is valid and delete from the storage
func DeleteUser(res http.ResponseWriter, req *http.Request) {
	panicRecover(res)
	if req.Method != "POST" {
		http.ServeFile(res, req, "views/deleteuser.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername,databasePassword string

	err := DB.QueryRow("SELECT username, password FROM loginuser WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	switch {
	case err==sql.ErrNoRows:
		errorHandler(res,fmt.Sprintf("username %s doesnot exist",username),utils.NotFound)
		return
	case nil!=err:
		log.Println(err.Error())
		errorHandler(res,utils.InternalServerErrorMsg,utils.ServerErrRespCode)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if nil != err {
		errorHandler(res,"password doesnot match",400)
		return
	}
	_, err = DB.Exec("DELETE FROM loginuser WHERE username=?",username)
	if nil!=err{
		log.Println(err.Error())
		errorHandler(res,utils.InternalServerErrorMsg,utils.ServerErrRespCode)
		return
	}
	res.Write([]byte( databaseUsername+ " successfully deleted"))

}
// HomePage serves the home page of the server
func HomePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "views/home.html")
}

func panicRecover(res http.ResponseWriter){
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			errorHandler(res,utils.InternalServerErrorMsg,utils.ServerErrRespCode)
		}
	}()
}
