package main

import (
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/rhizajay/gotham_admin"
)

func listUsers(c *cli.Context) {
	var db = connectDB(c)
	var admin = gotham_admin.GothamDB{db}
	for _, value := range admin.GetUsers() {
		value.DisplayUser()
	}
}

func userByEmail(c *cli.Context) {
	var db = connectDB(c)
	var admin = gotham_admin.GothamDB{db}
	admin.GetUserByEmail(c.Args().First()).DisplayUser()
}

func userByID(c *cli.Context) {
	var db = connectDB(c)
	var admin = gotham_admin.GothamDB{db}
	userid, err := strconv.Atoi(c.Args().First())
	if err != nil {
		println("Error : Not an number")
	}
	admin.GetUserById(userid).DisplayUser()
}

func activateAccountById(c *cli.Context){
	db := connectDB(c)
	admin := gotham_admin.GothamDB{db}
	userid, err := strconv.Atoi(c.Args().First())
	if err != nil {
		println("Error : Not an number")
	}
	admin.ActivateAccount(userid)

}

func deactivateAccountById(c *cli.Context){
	db := connectDB(c)
	admin := gotham_admin.GothamDB{db}
	userid, err := strconv.Atoi(c.Args().First())
	if err != nil {
		println("Error : Not an number")
	}
	admin.DeactivateAccount(userid)

}