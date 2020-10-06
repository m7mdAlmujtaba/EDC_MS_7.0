package main

import (
	"time"
)

// Action Struct ...
type Action struct{
	Aid, Areceipt, Fee int
	Auser, Atype,Asname, DAdate, DDate string
	Adate, Date time.Time
}

// Receipt
type Receipt struct{
	Std_name, Std_session, Std_Time, Recp_type,Std_lvl, Date string
	Std_id, Recp_id , Std_fees, Year int
}

// Records
type Record struct{
	Usr, Repdate, Reptime string
	Dptfee, Dptcount, Denfee, Dencount, Dcerfee, Dcercount, Dt int
	Wptfee, Wptcount, Wenfee, Wencount, Wcerfee, Wcercount, Wt int
	Mptfee, Mptcount, Menfee, Mencount, Mcerfee, Mcercount, Mt int
}

// Record Method to Calculate
func (r *Record)TotalFees(){
	r.Dt = r.Dptfee + r.Denfee + r.Dcerfee
	r.Wt = r.Wptfee + r.Wenfee + r.Wcerfee
	r.Mt = r.Mptfee + r.Menfee + r.Mcerfee
}


// Record Action (Student Enrollment / Placement Test / Certificate or Statment) // REGISTRAR
func Addaction(user_name string, atype string, std_id int, receipt_id int, date string, date_c time.Time,fee int) {
	// Start DB Connection
	db := DbConn()
	act_date, _ := time.Parse("2006-01-02", date)
	row, err := db.Query("INSERT INTO `edc_db`.`actions` (`user_name`, `type`, `student_id`, `receipt_id`, `act_date`, `date_col`, `fee`) VALUES(?, ?, ?, ?, ?, ?, ?);", user_name, atype, std_id, receipt_id, act_date, date_c, fee)
	CheckErr(err)
	row.Close()
	defer db.Close()
	return
}

// Display Actions // ADMIN
func ActDisp()[]Action{
	// Start DB Connection
	db := DbConn()
	act_row, err := db.Query("SELECT * FROM actions ORDER BY act_id DESC")
	CheckErr(err)
	act := Action{}
	acts := []Action{}
	var id int
	for act_row.Next(){
		err = act_row.Scan(&act.Aid, &act.Auser, &id, &act.Atype, &act.Areceipt, &act.Adate, &act.Date, &act.Fee) 
		CheckErr(err)
		act.DAdate = act.Adate.Format("2006-01-02")
		act.DDate = act.Date.Format("2006-01-02")
		act.Asname = GetStdName(id)
		acts = append(acts, act)
	}
	defer db.Close()
	return acts
}

func RegReport(usr string) Record{
	// Start DB Connection
	db := DbConn()
	rec := Record{}
	rec.Usr = usr
	rec.Repdate = time.Now().Format("2006-01-02")
	rec.Reptime = time.Now().Format("15:04:05")
	day_pt, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name =? AND `type` ='Placement Test' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_pt.Next(){
			err = day_pt.Scan(&rec.Dptfee, &rec.Dptcount)
			CheckErr(err)
		}
		
		day_en, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name =? AND `type` ='Enrollment' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_en.Next(){
			err = day_en.Scan(&rec.Denfee, &rec.Dencount)
			CheckErr(err)
		}

		day_cer, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statement' OR `type` = 'Freeze') ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_cer.Next(){
			err = day_cer.Scan(&rec.Dcerfee, &rec.Dcercount)
			CheckErr(err)
		}
	defer db.Close()
	return rec
}

// Get Number Of Students And Fees
func AdminReport() []Record{
	// Start DB Connection
	db := DbConn()
	// users list
	usrs := ActionsUsernames()
	rec := Record{}
	recs := []Record{}
	// struct to save record
	for _, usr := range usrs{ 
		rec.Usr = usr
		// Day Records
		day_pt, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name =? AND `type` ='Placement Test' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_pt.Next(){
			err = day_pt.Scan(&rec.Dptfee, &rec.Dptcount)
			CheckErr(err)
		}
		
		day_en, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name =? AND `type` ='Enrollment' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_en.Next(){
			err = day_en.Scan(&rec.Denfee, &rec.Dencount)
			CheckErr(err)
		}

		day_cer, err := db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) = DATE(NOW()) AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statement'  OR `type` = 'Freeze') ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for day_cer.Next(){
			err = day_cer.Scan(&rec.Dcerfee, &rec.Dcercount)
			CheckErr(err)
		}

		// Week Records
		week_pt, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= (DATE(NOW()) - INTERVAL 7 DAY) AND user_name =? AND `type` ='Placement Test' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for week_pt.Next(){
			err = week_pt.Scan(&rec.Wptfee, &rec.Wptcount)
			CheckErr(err)
		}

		week_en, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= DATE(NOW()) - INTERVAL 7 DAY AND user_name = ? AND `type` ='Enrollment' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for week_en.Next(){
			err = week_en.Scan(&rec.Wenfee, &rec.Wencount)
			CheckErr(err)
		}

		week_cer, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= DATE(NOW()) - INTERVAL 7 DAY AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statement'  OR `type` = 'Freeze') ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for week_cer.Next(){
			err = week_cer.Scan(&rec.Wcerfee, &rec.Wcercount)
			CheckErr(err)
		}

		// Month Records
		month_pt, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= DATE(NOW()) - INTERVAL 30 DAY AND user_name = ? AND `type` ='Placement Test' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for month_pt.Next(){
			err = month_pt.Scan(&rec.Mptfee, &rec.Mptcount)
			CheckErr(err)
		}
		
		month_en, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= DATE(NOW()) - INTERVAL 30 DAY AND user_name = ? AND `type` ='Enrollment' ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for month_en.Next(){
			err = month_en.Scan(&rec.Menfee, &rec.Mencount)
			CheckErr(err)
		}

		month_cer, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE DATE(act_date) >= DATE(NOW()) - INTERVAL 30 DAY AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statement'  OR `type` = 'Freeze') ORDER BY `act_id` DESC",usr)
		CheckErr(err)
		for month_cer.Next(){
			err = month_cer.Scan(&rec.Mcerfee, &rec.Mcercount)
			CheckErr(err)
		}

		rec.TotalFees()
		recs = append(recs, rec)
	}
	defer db.Close()
	return recs

}

// Each Month Record
func MonthRecord(date time.Time)[]Record{
	// Start DB Connection
	db := DbConn()
	// users list
	usrs := ActionsUsernames()
	rec := Record{}
	recs := []Record{}
	// struct to save record
	for _, usr := range usrs{ 
		rec.Usr = usr
		// Month Records
		month_pt, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE YEAR(act_date) = ? AND MONTH(act_date) = ? AND user_name = ? AND `type` ='Placement Test' ",date.Year() ,date.Month(), usr)
		CheckErr(err)
		for month_pt.Next(){
			err = month_pt.Scan(&rec.Mptfee, &rec.Mptcount)
			CheckErr(err)
		}
		
		month_en, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE YEAR(act_date) = ? AND MONTH(act_date) = ? AND user_name = ? AND `type` ='Enrollment' ",date.Year() ,date.Month(), usr)
		CheckErr(err)
		for month_en.Next(){
			err = month_en.Scan(&rec.Menfee, &rec.Mencount)
			CheckErr(err)
		}

		month_cer, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE YEAR(act_date) = ? AND MONTH(act_date) = ? AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statment'  OR `type` = 'Freeze')",date.Year() ,date.Month(), usr)
		CheckErr(err)
		for month_cer.Next(){
			err = month_cer.Scan(&rec.Mcerfee, &rec.Mcercount)
			CheckErr(err)
		}

		rec.TotalFees()
		recs = append(recs, rec)
	}
	defer db.Close()
	return recs
}

type Counts struct{
	/*t9, t11, t13, t15, t17 string*/
	Level string
	Times []int
	STimes []int
}

func SaStat() []Counts{
	// Start DB Connection
	db := DbConn()
	/*s := [2]string{"Regular", "Midmonth"}*/

	l := [11]string{"Pre zero", "Pre 1", "Pre 2", "Level 1","Level 2", "Level 3", "Level 4", "Level 5", "Level 6", "Level 7", "Level 8"}
	t := [5]string{"09:00 - 11:00", "11:00 - 13:00", "13:00 - 15:00", "15:00 - 17:00", "17:00 - 19:00"}

	var count,count2 int
	
	cnts := []Counts{}
	for _, x := range l{
		cnt := Counts{}
		cnt.Level = x
		for _, y := range t{
			regular, err :=  db.Query("SELECT count(id) FROM edc_db.students WHERE MONTH(last_reg) = MONTH(CURRENT_DATE()) AND YEAR(last_reg) = YEAR(CURRENT_DATE()) AND std_time = 'Regular' AND std_lvl=?  AND time=?;", x, y)
			for regular.Next(){
				err = regular.Scan(&count)
				cnt.Times = append(cnt.Times, count)
				CheckErr(err)
			}
			mid_month, err :=  db.Query("SELECT count(id) FROM edc_db.students WHERE MONTH(last_reg) = MONTH(CURRENT_DATE()) AND YEAR(last_reg) = YEAR(CURRENT_DATE()) AND std_time = 'Midmonth' AND std_lvl=?  AND time=?;", x, y)
			for mid_month.Next(){
				err = mid_month.Scan(&count2)
				cnt.STimes = append(cnt.STimes, count2)
				CheckErr(err)
			}
		}
		cnts = append(cnts, cnt)
	}
	defer db.Close()
	return cnts
}