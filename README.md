http server with following functionalites
1)Register
2)Login
3)Unregister

Steps to run 

1)git clone https://github.com/maddy1991/httpserver.git
2)go get github.com/go-sql-driver/mysql
3)go get golang.org/x/crypto/bcrypt
4)go run server.go

configutaion 

driver : mysql (currently supported database)
datasource : root:root123@tcp(127.0.0.1:3306)/ (mysql credentials)
serverip : 192.168.204.130 (server ip to listen )
serverport : 8080 (server port to listen)
ipformat : IPV4 (IP format)
