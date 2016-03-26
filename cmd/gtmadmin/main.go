package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rhizajay/gotham_admin"
)

func main() {
	app := cli.NewApp()
	app.Name = "gtmadmin"
	app.Usage = "administrator for gotham!"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dbhost, d",
			Value:  "127.0.0.1:13308",
			Usage:  "database_host:port",
			EnvVar: "RHIZA_USERDB",
		},
		cli.StringFlag{
			Name:   "user, u",
			Value:  "username",
			Usage:  "username",
			EnvVar: "RHIZA_ADMIN_USER",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "password",
			EnvVar: "RHIZA_ADMIN_PASS",
		},
		cli.StringFlag{
			Name:  "customer, c",
			Usage: "customer",
		},
	}

	app.Commands = []cli.Command{

		{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a user to group",
					Action: func(c *cli.Context) {
						println("a1 new task template: ", c.Args().First())
					},
				},
				{
					Name:  "remove",
					Usage: "remove a user from group",
					Action: func(c *cli.Context) {
						println("a2 aremoved task template: ", c.Args().First())
					},
				},
				{
					Name:  "list",
					Usage: "list groups",
					Action: func(c *cli.Context) {
						println("a3 removed task template: ", c.Args().First())
					},
				},
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "users commands",
			Subcommands: []cli.Command{
				{
					Name:  "activate",
					Usage: "turn user to active",
					Action: func(c *cli.Context) {
						println("b1 Activated: ", c.Args().First())
					},
				},
				{
					Name:  "deactivate",
					Usage: "turn user to inactive",
					Action: func(c *cli.Context) {
						println("b2 removed task template: ", c.GlobalString("customer"))
					},
				},
				{
					Name:   "list",
					Usage:  "list users",
					Action: listUsers,
				},
				{
					Name:   "email",
					Usage:  "search user by email",
					Action: userByEmail,
				},
				{
					Name:   "id",
					Usage:  "search user by id",
					Action: userByID,
				},
			},
		},
	}

	app.Run(os.Args)
}

func listUsers(c *cli.Context) {

	fmt.Println("got to listUsers")
	//println(c.FlagNames)
	println("display values")
	println(c.GlobalString("customer"))
	println(c.GlobalString("user"))

	var db = connectDB(c)
	var admin = gotham_admin.GothamDB{db}
	//fmt.Println(admin.GetUsers())
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

func connectDB(c *cli.Context) *sql.DB {

	println("connectDB")

	var customer = c.GlobalString("customer")
	var dbhost = c.GlobalString("dbhost")
	var username = c.GlobalString("user")
	var password = c.GlobalString("password")
	var db *sql.DB

	s := []string{customer, "_user"}
	customerdb := strings.Join(s, "")

	s = []string{username, ":", password, "@tcp(", dbhost, ")/", customerdb}
	connectString := strings.Join(s, "")

	var err error
	db, err = sql.Open("mysql", connectString)
	if err != nil {
		panic(err)
	}
	return db

}
