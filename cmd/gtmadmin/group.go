package main

import (
	"fmt"
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
	args := c.Args()

	groupid, err := strconv.Atoi(args[1])
	if err != nil {
		println("Error : Not an number")
	}
	userid, err := strconv.Atoi(args[0])
	if err != nil {
		println("Error : Not an number")
	}

	db := connectDB(c)
	admin := gotham_admin.GothamDB{DB: db}
	admin.SetGroup(userid, groupid)
	fmt.Printf("added %d to group: %d\n", userid, groupid)



}

func removeUserFromGroup (c *cli.Context){
	args := c.Args()

	groupid, err := strconv.Atoi(args[1])
	if err != nil {
		println("Error : Not an number")
	}
	userid, err := strconv.Atoi(args[0])
	if err != nil {
		println("Error : Not an number")
	}
	db := connectDB(c)
	admin := gotham_admin.GothamDB{DB: db}
	admin.DeleteUserFromGroup(userid, groupid)
	fmt.Printf("removed %d from group: %d\n", userid, groupid)}