package validator

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"utils"
)

var config Config

//Config configuration struct
type Config struct {
	Driver     string `yml:"driver"`
	DataSource string `yml:"datasource"`
	ServerIP   string `yml:"serverip"`
	ServerPort int    `yml:"serverport"`
	IPFormat   string `yml:"ipformat"`
}

//ConfigOps interface to get configurations
type ConfigOps interface {
	GetDriver() string
	GetDataSource() string
	GetServerADDR() string
}
type configst struct {
	*Config
}

var _ ConfigOps = configst{}

// NewConfig return the interface to retrieve configuration info
func NewConfig() ConfigOps {
	return configst{}
}
func (c configst) GetDriver() string {
	return config.Driver
}
func (c configst) GetDataSource() string {
	return config.DataSource
}
func (c configst) GetServerADDR() string {
	return config.ServerIP + ":" + strconv.Itoa(config.ServerPort)
}

func checkError(err error) {
	if nil != err {
		log.Fatal(err.Error())
	}
}

// ValidateConf validates the data in configuration file
func ValidateConf() {
	var (
		isErr bool
		conf  Config
	)
	var confErr = make([]string, 0)
	filename, err := filepath.Abs("./conf/config.yml")
	checkError(err)
	file, err := ioutil.ReadFile(filename)
	checkError(err)

	err = yaml.Unmarshal(file, &conf)
	checkError(err)
	if conf.DataSource == "" {
		isErr = true
		confErr = append(confErr, "Datasource is not present in conf file")

	}
	if conf.Driver == "" {
		isErr = true
		confErr = append(confErr, "Driver is not present in conf file")

	}
	if conf.ServerPort < 1 || conf.ServerPort > 65535 {
		isErr = true
		confErr = append(confErr, "provide valid server port to listen")

	}
	switch conf.IPFormat {
	case "IPV4":
		{
			if nil == net.ParseIP(strings.TrimSpace(conf.ServerIP)).To4() {
				isErr = true
				confErr = append(confErr, fmt.Sprintf("provided IPV4 address %s is invalid", conf.ServerIP))

			}
		}
	case "IPV6":
		{
			if nil == net.ParseIP(strings.TrimSpace(conf.ServerIP)).To16() {
				isErr = true
				confErr = append(confErr, fmt.Sprintf("provided IPV6 address %s is invalid", conf.ServerIP))

			}
		}
	default:
		isErr = true
		confErr = append(confErr, "check ip type IPV4/IPV6 in config.yml")
	}
	if isErr {
		for _, msg := range confErr {
			log.Println(msg)
		}
		os.Exit(1)
	}
	config = conf
}

func validateUserName(name string) error {
	if len(name) < 5 || len(name) > 15 {
		return errors.New("username should be min 5 char or max 15 char")
	}
	reName := regexp.MustCompile(utils.UsernameRegex)
	valid := reName.Match([]byte(name))

	if !valid {
		return errors.New("username should container alphanumaric")
	}
	return nil
}
func validatePassword(pwd string) error {
	if len(pwd) < 8 || len(pwd) > 15 {
		return errors.New("password should be min 8 char or max 15 char")
	}
	rePassword := regexp.MustCompile(utils.PasswordRegex)
	valid := rePassword.Match([]byte(pwd))

	if !valid {
		return errors.New("password should contain alphanumaric")
	}
	return nil
}

// ValidateUserInfo validates the username and password
func ValidateUserInfo(username string, password string) error {
	var err error
	err = validateUserName(username)
	if nil != err {

		return err
	}
	err = validatePassword(password)
	if nil != err {
		return err

	}
	return nil
}
