package gotham_admin

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type RhizaUserDB struct {
	DB *sql.DB
}

type RhizaUser struct {
	id       int
	username string
	email    string
	groups   map[int]bool
}

type RhizaGroup struct {
	group_id int
	title    string
}

type RhizaGroupMember struct {
	id int
	username string
	email string
}

func (r RhizaUserDB) GetUsers() []RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, group_id FROM account LEFT join account_to_account_group ON account.id = account_to_account_group.account_id ")
	checkErr(err)

	var userlist []RhizaUser

	var lastUser RhizaUser

	for rows.Next() {
		var i sql.NullInt64
		var group int

		var currentUser RhizaUser

		err = rows.Scan(&currentUser.id, &currentUser.email, &currentUser.username, &i)

		if i.Valid {
			group = int(i.Int64)
			if currentUser.id == lastUser.id {
				lastUser.groups[group] = true
			} else {
				currentUser.groups = make(map[int]bool)
				currentUser.groups[group] = true
				userlist = append(userlist, lastUser)
				lastUser = currentUser

			}
		}

	}

	return userlist

}

func (r RhizaUserDB) GetUserByEmail(s string) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE email=?", s)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int
		var group_title string

		err = rows.Scan(&user.id, &user.email, &user.username, &group)

		checkErr(err)
		user.groups[group] = true
		fmt.Println(user.id, group, group_title)
	}
	return user

}

func (r RhizaUserDB) GetUserById(id int) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE account.id=?", id)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int
		var group_title string

		err = rows.Scan(&user.id, &user.email, &user.username, &group)

		checkErr(err)
		user.groups[group] = true
		fmt.Println(user.id, group, group_title)
	}
	return user

}

func (r RhizaUserDB) GetGroupNames() []RhizaGroup {
	rows, err := r.DB.Query("SELECT id, title FROM account_group")
	checkErr(err)

	var groups []RhizaGroup
	var s sql.NullString

	for rows.Next() {
		var group RhizaGroup

		err = rows.Scan(&group.group_id, &s)

		if s.Valid {
			group.title = s.String
		} else {
			group.title = ""
		}

		checkErr(err)
		fmt.Println(group)

		groups = append(groups, group)
	}
	return groups

}

func (r RhizaUserDB) GetGroupName(groupid int) RhizaGroup {
	rows, err := r.DB.Query("SELECT id, title FROM account_group WHERE id=?", groupid)
	checkErr(err)

	var group RhizaGroup
	var s sql.NullString

	for rows.Next() {

		err = rows.Scan(&group.group_id, &s)

		if s.Valid {
			group.title = s.String
		} else {
			group.title = ""
		}

		checkErr(err)
		fmt.Println(group)

	}
	return group

}

func (r RhizaUserDB) SetGroup(userid int, groupid int) {

	user := r.GetUserById(userid)

	if user.groups[groupid] {
		fmt.Println("Already a member of this group")
	} else {

		stmt, err := r.DB.Prepare("INSERT account_to_account_group SET account_id=?, group_id=?")

		checkErr(err)

		res, err := stmt.Exec(userid, groupid)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)
	}

}

func (r RhizaUserDB) GetGroupMembers(groupid int) []RhizaGroupMember{

	rows, err := r.DB.Query("SELECT account_to_account_group.account_id, username, email FROM account JOIN account_to_account_group ON account_id = account_to_account_group.account_id WHERE account_to_account_group.group_id=? and account.id = account_to_account_group.account_id", groupid)
	checkErr(err)

	var members []RhizaGroupMember

	for rows.Next() {
		var user RhizaGroupMember

		err = rows.Scan(&user.id, &user.username, &user.email)

		checkErr(err)

		members = append(members, user)
	}
	return members

}


func (r RhizaUserDB) DeleteUserFromGroup(userid int, groupid int) {
	user := r.GetUserById(userid)

	if !user.groups[groupid] {
		fmt.Println("Not a member of this group")
	} else {

		stmt, err := r.DB.Prepare("DELETE FROM account_to_account_group WHERE account_id=? AND group_id=?")

		checkErr(err)

		res, err := stmt.Exec(userid, groupid)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)
	}
}


func (r RhizaUser) DisplayUser(){
	fmt.Printf("%d\t%s\t%s\t groups:", r.id, r.username, r.email)
	for key := range r.groups {
	 	fmt.Printf(" %d ",  key)
     }
	fmt.Printf("\n")
}




func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


