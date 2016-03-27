package main

import (
	"database/sql"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	//_ "github.com/go-sql-driver/mysql"
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
					Name:  "adduser",
					Usage: "add a user to group",
					Action: addUserToGroup,
				},
				{
					Name:  "removeuser",
					Usage: "remove a user from group",
					Action: removeUserFromGroup,
				},
				{
					Name:  "list",
					Usage: "list groups",
					Action: listGroups,
				},
								{
					Name:  "list_members",
					Usage: "list groups",
					Action: listMembers,
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
					Aliases: []string{"a"},
					Usage: "turn user to active",
					Action: activateAccountById,
				},
				{
					Name:  "deactivate",
					Aliases: []string{"d"},
					Usage: "turn user to inactive",
					Action: deactivateAccountById,
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
