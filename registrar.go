package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Placement Test
func PtGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "reg_pt.html", nil)
}

// Placement Test Post Request
func PtPostHandler(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	r.ParseForm()
	session, _ := store.Get(r, "session")
	var std = Student{}
	var std_id, receiptId int

	std.Sname = r.PostForm.Get("name")
	std.Sphone = r.PostForm.Get("phone")
	std.Sen_date, _ = time.Parse("2006-01-02", r.PostForm.Get("ptdate"))
	check := r.FormValue("checkbox")
	actDate := time.Now()

	//The Student id will be Generated using the actDate
	std.Sid, std_id, _ = std.GenerateId()
	AddStd(std)

	// Reciept id
	var laid int
	last_actid := db.QueryRow("SELECT IFNULL(MAX(`act_id`), 0) AS act_id FROM edc_db.actions;")
	err = last_actid.Scan(&laid)
	CheckErr(err)
	receiptId = laid + 1

	user_name, _ := session.Values["username"].(string)
	if check == "on" {
		fe := 0
		defer db.Close()
		Addaction(user_name, "Placement Test", std_id, receiptId, actDate.Format("2006-01-02"), std.Sen_date, fe)
		http.Redirect(w, r, "reg_dashboard", 302)
	} else {
		fe := GetFee("Placement Test")
		Addaction(user_name, "Placement Test", std_id, receiptId, actDate.Format("2006-01-02"), std.Sen_date, fe.Fee)

		//fmt.Println(user_name, "Placement Test", std_id, receiptId, actDate.Format("2006-01-02"), std.Sen_date,  fe.Fee)
		y := std.Sen_date.Year()
		fmt.Println(y)
		recp := Receipt{Std_name: std.Sname,
			Std_session: " ",
			Std_Time:    " ",
			Std_id:      std.Sid,
			Recp_id:     receiptId,
			Recp_type:   "Placement Test",
			Std_lvl:     " ",
			Std_fees:    fe.Fee,
			Date:        std.Sen_date.Format("2006-01-02"),
			Year:        y,
			User_name:   user_name,
		}

		var recps = []Receipt{recp, recp}
		defer db.Close()
		templates.ExecuteTemplate(w, "reg_receipt.html", recps)
	}
}

// Enrollment
func EnGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "reg_en.html", nil)
}

func EnPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	actDate := time.Now()
	r.ParseForm()
	var std = Student{}
	std.id, err = strconv.Atoi(r.PostForm.Get("id"))
	CheckErr(err)

	std.Slvl = r.PostForm.Get("level")
	std.Stype = r.PostForm.Get("type")
	std.Entime = r.PostForm.Get("time")   // 11:00 - 13:00
	std.Stime = r.PostForm.Get("session") // Regular / Midmonth
	std.Sen_date, err = time.Parse("2006-01-02", r.PostForm.Get("endate"))
	repeat := r.FormValue("checkbox")

	//Check if Empty
	checkLvl := IsEmpty(std.Slvl)
	checkType := IsEmpty(std.Stype)
	checkTime := IsEmpty(std.Entime)
	checkSession := IsEmpty(std.Stime)
	checkEndate := IsEmpty(r.PostForm.Get("endate"))

	if checkLvl || checkType || checkTime || checkSession || checkEndate {
		templates.ExecuteTemplate(w, "reg_msg.html", "Please Fill all Feilds")
		return
	}

	UpdateStd(std)
	//fmt.Println(std.Sid)

	user_name, _ := session.Values["username"].(string)

	// Start DB Connection
	db := DbConn()

	// Get id using std_id
	var std_id int
	std_id = std.id
	/*id_row := db.QueryRow("SELECT id from students WHERE id=?",std.id)
	err = id_row.Scan(&std_id)
	CheckErr(err)*/

	std.Sname = GetStdName(std_id)

	// Reciept id
	var laid int
	last_actid := db.QueryRow("SELECT IFNULL(MAX(`act_id`), 0) AS act_id FROM edc_db.actions;")
	err = last_actid.Scan(&laid)
	CheckErr(err)
	receiptId := laid + 1
	var fe Fees
	if std.Stype == "Communication" || std.Stype == "English Club" {
		fe = GetFee(std.Slvl)
	} else {
		fe = GetFee(std.Stype)
	}
	if repeat == "on" {
		fe.Fee = fe.Fee / 2
	}
	Addaction(user_name, "Enrollment", std_id, receiptId, actDate.Format("2006-01-02"), std.Sen_date, fe.Fee)

	y := std.Sen_date.Year()

	recp := Receipt{Std_name: std.Sname,
		Std_session: std.Stime,
		Std_Time:    std.Entime,
		Std_id:      std.id,
		Recp_id:     receiptId,
		Recp_type:   "Enrollment",
		Std_lvl:     std.Slvl,
		Std_fees:    fe.Fee,
		Date:        std.Sen_date.Format("2006-01-02"),
		Year:        y,
		User_name:   user_name,
	}

	var recps = []Receipt{recp, recp}
	defer db.Close()
	templates.ExecuteTemplate(w, "reg_receipt.html", recps)

}

// Certificate
func CertGetHandler(w http.ResponseWriter, r *http.Request) {
	err = templates.ExecuteTemplate(w, "reg_cert.html", nil)
	CheckErr(err)
}

func CertPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	actDate := time.Now()
	r.ParseForm()
	var std = Student{}
	std.id, err = strconv.Atoi(r.PostForm.Get("id"))
	CheckErr(err)
	std.Sen_date, err = time.Parse("2006-01-02", r.PostForm.Get("certdate"))
	statmentType := r.PostForm.Get("statment_type")
	std.Slvl = r.PostForm.Get("level")
	user_name, _ := session.Values["username"].(string)

	////Check if Empty
	checkType := IsEmpty(statmentType)
	checkCerdate := IsEmpty(r.PostForm.Get("certdate"))

	if checkType || checkCerdate {
		templates.ExecuteTemplate(w, "reg_msg.html", "Please Fill Both Date and Type Feilds")
		return
	}

	// Start DB Connection
	db := DbConn()
	var std_id int
	std_id = std.id
	/* Get id using std_id
	id_row := db.QueryRow("SELECT id from students WHERE std_id=?",std.Sid)
	err = id_row.Scan(&std_id)
	CheckErr(err) */

	std.Sname = GetStdName(std_id)

	// Reciept id
	var laid int
	last_actid := db.QueryRow("SELECT IFNULL(MAX(`act_id`), 0) AS act_id FROM edc_db.actions;")
	err = last_actid.Scan(&laid)
	CheckErr(err)
	receiptId := laid + 1

	fe := GetFee(statmentType)

	// Include (Student Name + type of statment , Date, Receipt no, Fees) on The Receipt

	Addaction(user_name, statmentType, std_id, receiptId, actDate.Format("2006-01-02"), std.Sen_date, fe.Fee)
	y := std.Sen_date.Year()

	recp := Receipt{Std_name: std.Sname,
		Std_session: "",
		Std_Time:    "",
		Std_id:      std.id,
		Recp_id:     receiptId,
		Recp_type:   statmentType,
		Std_lvl:     std.Slvl,
		Std_fees:    fe.Fee,
		Date:        std.Sen_date.Format("2006-01-02"),
		Year:        y,
		User_name:   user_name,
	}

	var recps = []Receipt{recp, recp}
	defer db.Close()
	templates.ExecuteTemplate(w, "reg_receipt.html", recps)

}

// Report
func ReportGetHandler(w http.ResponseWriter, r *http.Request) {
	rec := Record{}
	session, _ := store.Get(r, "session")
	usr, _ := session.Values["username"].(string)
	rec = RegReport(usr)

	err = templates.ExecuteTemplate(w, "reg_report.html", rec)
	CheckErr(err)

}
