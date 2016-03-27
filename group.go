package gotham_admin

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type RhizaGroup struct {
	group_id int
	title    string
}

func (r GothamDB) GetGroupNames() []RhizaGroup {
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

		groups = append(groups, group)
	}
	return groups

}

func (r GothamDB) GetGroupName(groupid int) RhizaGroup {
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

	}
	return group

}

func (r GothamDB) SetGroup(userid int, groupid int) {

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

func (r GothamDB) GetGroupMembersByGroupId(groupid int) []RhizaUser {

	rows, err := r.DB.Query("SELECT account_to_account_group.account_id, email, username, is_active FROM account JOIN account_to_account_group ON account_id = account_to_account_group.account_id WHERE account_to_account_group.group_id=? and account.id = account_to_account_group.account_id", groupid)
	checkErr(err)

	var members []RhizaUser

	for rows.Next() {
		var user RhizaUser

		err = rows.Scan(&user.id, &user.email, &user.username, &user.is_active)
		user.groups = make(map[int]bool)
		user.groups[groupid] = true
		checkErr(err)

		members = append(members, user)
	}
	return members

}

func (r GothamDB) DeleteUserFromGroup(userid int, groupid int) {
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

// type RhizaGroup struct {
// 	group_id int
// 	title    string
// }

func (r RhizaGroup) DisplayGroup() {
	fmt.Printf("%d\t%s\n", r.group_id, r.title)
}
