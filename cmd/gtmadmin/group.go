package main

import (
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/rhizajay/gotham_admin"
)

func listMembers(c *cli.Context){
	db := connectDB(c)
	admin := gotham_admin.GothamDB{DB: db}
	groupid, err := strconv.Atoi(c.Args().First())
	if err != nil {
		println("Error : Not an number")
	}

	for _, value := range admin.GetGroupMembersByGroupId(groupid) {
		value.DisplayUser()
	}

}

func listGroups(c *cli.Context){
	db := connectDB(c)
	admin := gotham_admin.GothamDB{DB: db}
	for _, value := range admin.GetGroupNames() {
		value.DisplayGroup()
	}
}

func addUserToGroup (c *cli.Context){
	println("added")
}

func removeUserFromGroup (c *cli.Context){
	println("removed")
}