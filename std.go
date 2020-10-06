package main

import (
	"fmt"
	"time"
	"strconv"
    "database/sql"
)

const layoutISO = "2006-01-02"

// Student Struct
type Student struct{
	Sid int
	Sname string
	Sphone string
	Slvl string
	Sen_date time.Time  // Student Enrollment Date / First Date
	Re_date time.Time   // Receipt Date
	Stype string   // Enrollment Type
	Entime string  // 11:00-13:00, 13:00-15:00 ...etc
	Stime string  // Regular or mid-month
	D1,D2 string
}

// Generate Student uid from 8 digits (xxxxyyyy)
func (std Student) GenerateId(date time.Time) (ValidId int, lid int, err error){
	// Start DB Connection
	db := DbConn()
	var xr int
	last_id := db.QueryRow("SELECT IFNULL(MAX(`id`), 0) AS id FROM edc_db.students;")
	err = last_id.Scan(&lid)
	CheckErr(err)
	xr = lid % 99999
	x := fmt.Sprintf("%05d", xr + 1)

	y := date.Format("2006-01-02")[2:4] + date.Format("2006-01-02")[5:7]

	lid = lid + 1
	si, err := strconv.Atoi(y + x)
	ValidId = checkId(si, x, y)
	defer db.Close()
	return ValidId, lid, err
}

// Check if id is Exist if not (Note: Used with Placement Test)
func checkId(si int,x string ,y string)(int){
	// Start DB Connection
	db := DbConn()
	var checked_id,tempid int
	var nx string
	checked_id = si
	err := db.QueryRow(`SELECT std_id FROM edc_db.students WHERE std_id=?`, checked_id).Scan(&tempid)
	defer db.Close()
	if err != nil {
		if err == sql.ErrNoRows{
			return checked_id
		}else{
			panic(err)
		}
	}else {
		for i := 0; i < 99999; i++{
			nx = fmt.Sprintf("%05d", i)
			checked_id, _ = strconv.Atoi(y + nx)
			db := DbConn()
			err := db.QueryRow(`SELECT std_id FROM edc_db.students WHERE std_id=?`, checked_id).Scan(&tempid)
			defer db.Close()
			if err != nil {
				if err == sql.ErrNoRows{
					return checked_id
				}else{
					panic(err)
				}
			}
		}
	}
	
	return checked_id
}

// Get student id using his uid
func GetStdId(std_id int) (id int) {
	// Start DB Connection
	db := DbConn()
	std_wid, err := db.Query("SELECT id FROM students WHERE `std_id` = ?", std_id)
	CheckErr(err)
	err = std_wid.Scan(&id)
	CheckErr(err)
	fmt.Println(id)
	defer db.Close()
	return id
}

// Get student name using his id
func GetStdName(id int)(name string){
	// Start DB Connection
	db := DbConn()
	row, nerr := db.Query("SELECT std_name FROM students WHERE `id` = ?", id)
	CheckErr(nerr)
	for row.Next(){
		err = row.Scan(&name)
		CheckErr(err)
	}
	defer db.Close()
	return name
}

// To Add Student to The Datebase (Placement Test)
func AddStd(std Student){
	// Start DB Connection
	db := DbConn()
	//var d = time.Now()
	l := ""
	d, _ := time.Parse("2006-01-02", "1000-01-01")
	//d.Format(time.RFC3339)
	t := ""
	ti := ""
	s := ""
	row, err := db.Query("INSERT INTO `edc_db`.`students` (`std_id`, `std_name`, `std_phone`, `std_lvl`,`en_date`, `last_reg`, `en_type`, `time`, `std_time`) VALUES(?, ?, ?, ?, ? ,?, ?, ?, ?);", std.Sid, std.Sname, std.Sphone, l, d,std.Sen_date, t, ti, s)
	CheckErr(err)
	row.Close()
	// defer db.Close()
	defer db.Close()
	return
}

// Student Enrollment
func UpdateStd(std Student){
	// Start DB Connection
	db := DbConn()
	// Get en Date using std_id

	var enroll_date time.Time
	date_row := db.QueryRow("SELECT en_date from students WHERE std_id=?",std.Sid)
	err = date_row.Scan(&enroll_date)
	CheckErr(err)

	tempDate, err := time.Parse("2006-01-02","1000-01-01")
	CheckErr(err)

	if enroll_date == tempDate {
		row1, err := db.Query("UPDATE `edc_db`.`students` SET `std_lvl` = ?,`en_date` =?, `last_reg`=?, `en_type` =?, `time`=?, `std_time` = ? WHERE `std_id` = ?",std.Slvl, std.Sen_date, std.Sen_date, std.Stype, std.Entime, std.Stime, std.Sid)
		CheckErr(err)
		row1.Close()
	}else{
		row2, err := db.Query("UPDATE `edc_db`.`students` SET `std_lvl` = ?, `last_reg` =?, `en_type` =?, `time`=?, `std_time` = ? WHERE `std_id` = ?",std.Slvl, std.Sen_date, std.Stype, std.Entime, std.Stime, std.Sid)
		CheckErr(err)
		row2.Close()
	}
	defer db.Close()
	return
}

// List Of All Students... To be Displayed
func StdDisp() []Student{
	// Start DB Connection
	db := DbConn()
	std_row, err := db.Query("SELECT * FROM students ORDER BY id DESC")
	CheckErr(err)
	std := Student{}
	stds := []Student{}
	var id int
	for std_row.Next(){
		err = std_row.Scan(&id, &std.Sid , &std.Sname, &std.Sphone, &std.Slvl, &std.Sen_date,  &std.Re_date, &std.Stype, &std.Entime, &std.Stime) 
		CheckErr(err)
		std.D1 = std.Sen_date.Format("2006-01-02")
		std.D2 = std.Re_date.Format("2006-01-02")
		stds = append(stds, std)
	}
	defer db.Close()
	return stds
}

// To Display on Enrollment Table 
func StdEn() []Student{
	// Start DB Connection
	db := DbConn()
	std_row, err := db.Query("SELECT std_id, std_name FROM students")
	CheckErr(err)
	std := Student{}
	stds := []Student{}
	for std_row.Next(){
		err = std_row.Scan(&std.Sid , &std.Sname) 
		CheckErr(err)
		stds = append(stds, std)
	}
	defer db.Close()
	return stds
}

func CheckField(std Student){
}

//Get Student id, name and phone to be Exported as Excel
func StdsExcel(std_type, session, level, time string)[]Student{
	// Start DB Connection
	db := DbConn()
	std_row, err := db.Query("SELECT std_id, std_name, std_phone FROM students WHERE (DATE(last_reg) >= DATE(NOW()) - INTERVAL 30 DAY) AND `en_type` =? AND `std_time` =? AND `std_lvl` =? AND `time` =?", std_type, session, level, time)
	CheckErr(err)
	std := Student{}
	stds := []Student{}
	for std_row.Next(){
		err = std_row.Scan(&std.Sid , &std.Sname, &std.Sphone) 
		CheckErr(err)
		stds = append(stds, std)
	}
	defer db.Close()
	return stds

}

func AllStdsExcel()[]Student{
	// Start DB Connection
	db := DbConn()
	std_row, err := db.Query("SELECT std_id, std_name, std_phone FROM students")
	CheckErr(err)
	std := Student{}
	stds := []Student{}
	for std_row.Next(){
		err = std_row.Scan(&std.Sid , &std.Sname, &std.Sphone) 
		CheckErr(err)
		stds = append(stds, std)
	}
	defer db.Close()
	return stds

}

// Student Enrollment
func EditStd(std Student){
	// Start DB Connection
	db := DbConn()
	// Get en Date using std_id

	var enroll_date time.Time
	date_row := db.QueryRow("SELECT en_date from students WHERE std_id=?",std.Sid)
	err = date_row.Scan(&enroll_date)
	CheckErr(err)


	row2, err := db.Query("UPDATE `edc_db`.`students` SET `std_name` =?, `std_lvl` = ?, `en_type` =?, `time`=?, `std_time` = ? WHERE `std_id` = ?", std.Sname, std.Slvl, std.Stype, std.Entime, std.Stime, std.Sid)
	CheckErr(err)
	row2.Close()
	
	defer db.Close()
	return
}
