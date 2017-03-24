package utils

const(
	//DatabaseName database name
	DatabaseName ="userinfo"
	//InternalServerErrorMsg internal server error message
	InternalServerErrorMsg ="encountered internal server error"
	//UsernameRegex username regex
	UsernameRegex =`^[a-zA-Z][a-zA-Z0-9]*[._]?[a-zA-Z0-9]+$`
	//PasswordRegex password regex
	PasswordRegex =`[0-9a-zA-Z\s]*`
	//ServerErrRespCode internal server response code
	ServerErrRespCode =500
	//InvalidRespCode invalid response code
	InvalidRespCode =400
	//NotFound not found response code
	NotFound=404
)
