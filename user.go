package gotham_admin

import (
	"database/sql"
	"fmt"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)



type RhizaUser struct {
	id       int
	username string
	email    string
	// is_active bool
	// created time.Time
	// lastLogin time.Time
	groups   map[int]bool
}



func (r GothamDB) GetUsers() []RhizaUser {

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

func (r GothamDB) GetUserByEmail(s string) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE email=?", s)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.id, &user.email, &user.username, &group)

		checkErr(err)
		user.groups[group] = true
	}
	return user

}

func (r GothamDB) GetUserById(id int) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE account.id=?", id)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.id, &user.email, &user.username, &group)

		checkErr(err)
		user.groups[group] = true
	}
	return user

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


