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
	is_active int
	// created time.Time
	// lastLogin time.Time
	groups map[int]bool
}

func (r GothamDB) GetUsers() []RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, is_active, group_id FROM account LEFT join account_to_account_group ON account.id = account_to_account_group.account_id ")
	checkErr(err)

	var userlist []RhizaUser

	var lastUser RhizaUser

	for rows.Next() {
		var nullableid sql.NullInt64
		var group int

		var currentUser RhizaUser

		err = rows.Scan(&currentUser.id, &currentUser.email, &currentUser.username, &currentUser.is_active, &nullableid)

		if nullableid.Valid {
			group = int(nullableid.Int64)
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

	fmt.Printf("Number of users: %d\n", len(userlist))
	return userlist

}

func (r GothamDB) GetUserByEmail(s string) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, is_active, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE email=?", s)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.id, &user.email, &user.username, &user.is_active, &group)

		checkErr(err)
		user.groups[group] = true
	}
	return user

}

func (r GothamDB) GetUserById(id int) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, is_active, group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE account.id=?", id)
	checkErr(err)
	var user RhizaUser
	user.groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.id, &user.email, &user.username, &user.is_active, &group)

		checkErr(err)
		user.groups[group] = true
	}
	return user

}

func (r RhizaUser) DisplayUser() {
	fmt.Printf("%d\t%s\t%s\t", r.id, r.username, r.email)

	if r.is_active == 1{
		fmt.Printf("\tactive\t")		
	} else {
		fmt.Printf("\tinactive\t")
	}
	fmt.Printf(" groups:")
	for key := range r.groups {
		fmt.Printf(" %d ", key)
	}
	fmt.Printf("\n")
}

func (r GothamDB) DeactivateAccount(userid int) {

	stmt, err := r.DB.Prepare("UPDATE account SET is_active=0 WHERE id=?")

	checkErr(err)

	res, err := stmt.Exec(userid)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Printf("Number of records updated: %d\n",affect)
}

func (r GothamDB) ActivateAccount(userid int) {

	stmt, err := r.DB.Prepare("UPDATE account SET is_active=1 WHERE id=?")

	checkErr(err)

	res, err := stmt.Exec(userid)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Printf("Number of records updated: %d\n",affect)
}

func (r GothamDB) DeactivateAccountByEmail(email string) {
	user := r.GetUserByEmail(email)
	r.DeactivateAccount(user.id)

}

func (r GothamDB) ActivateAccountByEmail(email string) {
	user := r.GetUserByEmail(email)
	r.ActivateAccount(user.id)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
