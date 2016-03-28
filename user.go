package gotham_admin

import (
	"database/sql"
	"fmt"

	//"time"

	_ "github.com/go-sql-driver/mysql"
)

type RhizaUser struct {
	Id       int
	Username string
	Email    string
	Is_active int
	// created time.Time
	// lastLogin time.Time
	Groups map[int]bool
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

		err = rows.Scan(&currentUser.Id, &currentUser.Email, &currentUser.Username, &currentUser.Is_active, &nullableid)

		if nullableid.Valid {
			group = int(nullableid.Int64)
			if currentUser.Id == lastUser.Id {
				lastUser.Groups[group] = true
			} else {
				currentUser.Groups = make(map[int]bool)
				currentUser.Groups[group] = true
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
	user.Groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.Id, &user.Email, &user.Username, &user.Is_active, &group)

		checkErr(err)
		user.Groups[group] = true
	}
	return user

}

func (r GothamDB) GetUserById(id int) RhizaUser {

	rows, err := r.DB.Query("SELECT account.id, email, username, is_active, Group_id FROM account join account_to_account_group ON account.id = account_to_account_group.account_id WHERE account.id=?", id)
	checkErr(err)
	var user RhizaUser
	user.Groups = make(map[int]bool)

	for rows.Next() {
		var group int

		err = rows.Scan(&user.Id, &user.Email, &user.Username, &user.Is_active, &group)

		checkErr(err)
		user.Groups[group] = true
	}
	return user

}

func (r RhizaUser) DisplayUser() {
	fmt.Printf("%d\t%s\t%s\t", r.Id, r.Username, r.Email)

	if r.Is_active == 1{
		fmt.Printf("\tactive\t")		
	} else {
		fmt.Printf("\tinactive\t")
	}
	fmt.Printf(" groups:")
	for key := range r.Groups {
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
	r.DeactivateAccount(user.Id)

}

func (r GothamDB) ActivateAccountByEmail(email string) {
	user := r.GetUserByEmail(email)
	r.ActivateAccount(user.Id)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
