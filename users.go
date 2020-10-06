package main 

import(

)

type User struct {
	ID        int
	Uname  string
	Upass string
	Utype  string
}

// Query User for login Confirmation
func QueryUser(username string) User {
	db := DbConn()
	var Users = User{}
	row, err := db.Query(`
		SELECT user_id, 
		user_name, 
		user_pass, 
		user_type
		FROM edc_db.users WHERE BINARY user_name=?`, username)
	CheckErr(err)
	for row.Next(){
		err = row.Scan(&Users.ID, &Users.Uname, &Users.Upass, &Users.Utype)
		CheckErr(err)
	}
	defer db.Close()
	return Users
}

// Get the usernames
func Usernames() []string{
	db := DbConn()
	usr_row, err := db.Query("SELECT user_name FROM users WHERE `user_type` = 'registrar' OR `user_type` = 'admin' ")
	CheckErr(err)
	var user string
	var users []string
	for usr_row.Next(){
		err = usr_row.Scan(&user) 
		CheckErr(err)
		users = append(users, user)
	}
	defer db.Close()
	return users
}

// Display Users
func UsersDisp() []User{
	db := DbConn()
	usr_row, err := db.Query("SELECT * FROM users")
	CheckErr(err)
	user := User{}
	users := []User{}
	for usr_row.Next(){
		err = usr_row.Scan(&user.ID, &user.Uname, &user.Upass, &user.Utype) 
		CheckErr(err)
		users = append(users, user)
	}
	defer db.Close()
	return users
}

// Add User
func AddUser(user User){
	db := DbConn()
	row, err := db.Query("INSERT INTO `edc_db`.`users` (`user_name`,`user_pass`, `user_type`) VALUES(?, ?, ?);",user.Uname, user.Upass, user.Utype)
	CheckErr(err)
	row.Close()
	defer db.Close()
}

// Remove User
func DelUser(id int){
	db := DbConn()
	row, err := db.Query("DELETE FROM `edc_db`.`users` WHERE user_id=?", id)
	CheckErr(err)
	row.Close()
	defer db.Close()
}

func Editpass(user User){
	db := DbConn()
	row1, err := db.Query("UPDATE `edc_db`.`users` SET `user_pass` = ? WHERE `user_id` = ?", user.Upass, user.ID)
	CheckErr(err)
	row1.Close()
	defer db.Close()
}

// Select Usernames from Actions Table
func ActionsUsernames() []string{
	db := DbConn()
	usr_row, err := db.Query("SELECT DISTINCT user_name FROM actions")
	CheckErr(err)
	var user string
	var users []string
	for usr_row.Next(){
		err = usr_row.Scan(&user) 
		CheckErr(err)
		users = append(users, user)
	}
	defer db.Close()
	return users
}

func IsEmpty(data string) bool{
	if len(data) == 0 {
		return true
	}else {
		return false
	}
}

// Check Username
func IsExist(usr string) bool{
	db := DbConn()
	var name string
	err := db.QueryRow("SELECT user_name FROM users WHERE  user_name=?", usr).Scan(&name)
	defer db.Close()
	if err == nil{
		return true
	}else{
		return false
	}
}