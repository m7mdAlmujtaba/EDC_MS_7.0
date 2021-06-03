package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"time"
) 

type Scount struct{
	Ts, Cs, Ps , U int
	List [10]int
}

// Users
func usersGetHandler(w http.ResponseWriter,r *http.Request){
	var users = []User{}
	users = UsersDisp()
	templates.ExecuteTemplate(w, "admin_users.html", users)
}

func usersPostHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var user User
	user.Uname = r.PostForm.Get("username")
	user.Upass = r.PostForm.Get("userpass")
	user.Utype = r.PostForm.Get("usertype")

	user_space := strings.Contains(user.Uname, " ")

	emptyuser := IsEmpty(user.Uname)
	emptypass := IsEmpty(user.Upass)
	if emptyuser || emptypass {
		templates.ExecuteTemplate(w, "msg.html", "Invalid Username or Password (Empty!)")
		return
	}
	if user_space{
		templates.ExecuteTemplate(w, "msg.html", "Username can not contain space")
		return
	}
	existuser := IsExist(user.Uname)
	if existuser{
		templates.ExecuteTemplate(w, "msg.html", "This Username is already Exist .. Try Another Username")
		return
	}

	AddUser(user) 
	http.Redirect(w, r, "/admin_users", 302)
}

func delUsersPostHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var user User
	user.ID, _ = strconv.Atoi(r.PostForm.Get("id"))
	DelUser(user.ID)
	http.Redirect(w, r, "/admin_users", 302)
}

func editpassHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var user User
	user.ID, _ = strconv.Atoi(r.PostForm.Get("id"))
	user.Upass = r.PostForm.Get("userpass")
	emptypass := IsEmpty(user.Upass)
	if emptypass {
		templates.ExecuteTemplate(w, "msg.html", "Invalid Username or Password (Empty!)")
		return
	}
	Editpass(user)
	http.Redirect(w, r, "/admin_users", 302)
}


// Fees
func feesGetHandler(w http.ResponseWriter,r *http.Request){
	var fees = []Fees{}
	fees = FeesDisp()
	templates.ExecuteTemplate(w, "admin_fees.html", fees)
}

func feesPostHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var fee Fees
	if r.PostForm.Get("editfees") == ""{
		templates.ExecuteTemplate(w, "msg.html", "Please Enter a Valid Value")
		return
	}
	fee.Id, err = strconv.Atoi(r.PostForm.Get("id"))
	CheckErr(err)
	fee.Fee, err = strconv.Atoi(r.PostForm.Get("editfees"))
	CheckErr(err)
	if fee.Fee < 0 {
		templates.ExecuteTemplate(w, "msg.html", "Please Enter a Valid Value (Positive)")
		return
	}
	EditFees(fee)
	http.Redirect(w, r, "/admin_fees", 302)
}


// Report
func actionsGetHandler(w http.ResponseWriter,r *http.Request){
	var acts = []Action{}
	acts = ActDisp()
	templates.ExecuteTemplate(w, "admin_actions.html", acts)
}

func reportsGetHandler(w http.ResponseWriter,r *http.Request){
	var recs = []Record{}
	recs = AdminReport()
	templates.ExecuteTemplate(w, "admin_report.html", recs)
}

// Stat

func recordsGetHandler(w http.ResponseWriter,r *http.Request){
	var counts Scount
	counts.Ts = GetTotalRows("SELECT count(*) FROM students")
	counts.Cs = GetTotalRows("SELECT count(*) FROM edc_db.students where en_date != '1000-01-01' and year(last_reg) = year(NOW()) and month(last_reg) = month(NOW()) ;")
	counts.Ps = GetTotalRows("SELECT count(*) FROM edc_db.students where en_date = '1000-01-01' and year(last_reg) = year(NOW()) and month(last_reg) = month(NOW()) ;")
	counts.U = GetTotalRows("SELECT count(*) FROM users")
	// No need
	counts.List = [10]int{12, 19, 3, 5, 2, 3, 19, 10 , 22, 1}
	templates.ExecuteTemplate(w, "admin_records.html", counts)
}

// Get The Total Number of Students in Each Level
type LvlsCount struct{
	Regular []int `json:"Regular"`
	Midmonth []int `json:"Midmonth"`
} 

func LvlCountHandler(w http.ResponseWriter,r *http.Request){
	db := DbConn()
	var n, n2 int
	l := [11]string{"Pre zero", "Pre 1", "Pre 2", "Level 1","Level 2", "Level 3", "Level 4", "Level 5", "Level 6", "Level 7", "Level 8"}
	var cnt []int
	var cnt2 []int
	for _, x := range l{
		regular, err :=  db.Query("SELECT count(id) FROM edc_db.students WHERE MONTH(last_reg) = MONTH(CURRENT_DATE()) AND YEAR(last_reg) = YEAR(CURRENT_DATE()) AND std_time = 'Regular' AND std_lvl=? ;", x)
			for regular.Next(){
				err = regular.Scan(&n)
				cnt = append(cnt, n)
				CheckErr(err)
			}
			mid_month, err :=  db.Query("SELECT count(id) FROM edc_db.students WHERE MONTH(last_reg) = MONTH(CURRENT_DATE()) AND YEAR(last_reg) = YEAR(CURRENT_DATE()) AND std_time = 'Midmonth' AND std_lvl=? ;", x)
			for mid_month.Next(){
				err = mid_month.Scan(&n2)
				cnt2 = append(cnt2, n2)
				CheckErr(err)
			}
	}
	levelscnt := LvlsCount{Regular: cnt, Midmonth: cnt2}
	fmt.Println(levelscnt)
	defer db.Close()
	json.NewEncoder(w).Encode(levelscnt)
}

// Number of Students in Each Level and Time
func recordsPostHandler(w http.ResponseWriter,r *http.Request){
	var recs = []Record{}
	d1 := r.URL.Query().Get("date1")
	d2 := r.URL.Query().Get("date2")

	fmt.Println(d1, d2)
	if d1 == "" || d2 == ""{
		d1 := time.Now().Format("2006-01-02")
		d2 := time.Now().Format("2006-01-02")
		date1, _ := time.Parse("2006-01-2", d1 )
		date2, _ := time.Parse("2006-01-2", d2)
		recs = Recordbetween(date1, date2)
	}else{
		date1, _ := time.Parse("2006-01-2", d1 )
		date2, _ := time.Parse("2006-01-2", d2)
		if date1.After(date2){
			recs = Recordbetween(date2, date1)
		}else{
			recs = Recordbetween(date1, date2)
		}
	}

	
	html :=`                    
	<table id="myTable" class="table table-dark">
	<thead>
	  <tr>
		<th scope="col">User</th>
		<th scope="col">Placement Test Students</th>
		<th scope="col">Enrolled Students</th>
		<th scope="col">Certificates/Statements</th>
		<th scope="col">Total Income</th>
	  </tr>
	</thead>
	<tbody>`

	for _,v := range recs{
		html += `<tr>
		<td>`+ v.Usr +`</td>
		<td>`+ strconv.Itoa(v.Mptcount) + `</td>
		<td>`+ strconv.Itoa(v.Mencount) + `</td>
		<td>`+ strconv.Itoa(v.Mcercount) +`</td>
		<td>`+ strconv.Itoa(v.Mt) +`</td>
	</tr>`
	}
			
	html += `
	<tr>
		<th>Total</th>
			<td id="total_pt"></td>
			<td id="total_en"></td>
			<td id="total_cer"></td>
			<td id="total_in"></td>
	</tr>
	</tbody>
	</table>
	`
fmt.Fprintf(w, html)
}

func Recordbetween(date1, date2 time.Time)[]Record{
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
		month_pt, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE (act_date between ? and ?) AND user_name = ? AND `type` ='Placement Test' ",date1 ,date2, usr)
		CheckErr(err)
		for month_pt.Next(){
			err = month_pt.Scan(&rec.Mptfee, &rec.Mptcount)
			CheckErr(err)
		}
		
		month_en, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE  (act_date between ? and ?) AND user_name = ? AND `type` ='Enrollment' ",date1 ,date2, usr)
		CheckErr(err)
		for month_en.Next(){
			err = month_en.Scan(&rec.Menfee, &rec.Mencount)
			CheckErr(err)
		}

		month_cer, err :=  db.Query("SELECT IFNULL(sum(fee), 0), count(act_id) FROM actions WHERE  (act_date between ? and ?) AND user_name = ? AND (`type` ='Certificate' OR `type` = 'Statement'  OR `type` = 'Freeze')",date1 ,date2, usr)
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

func editStd(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var std = Student{}
	std.id, err = strconv.Atoi(r.PostForm.Get("id"))
	CheckErr(err)

	std.Sname = r.PostForm.Get("name")
	std.Sphone = r.PostForm.Get("phone")
	std.Slvl = r.PostForm.Get("level")
	std.Stype = r.PostForm.Get("type")
	std.Entime = r.PostForm.Get("time") // 11:00 - 13:00
	std.Stime = r.PostForm.Get("session") // Regular / Midmonth

	//Check if Empty
	checkname := IsEmpty(std.Sname)
	checkphone := IsEmpty(std.Sphone)
	checkLvl := IsEmpty(std.Slvl)
	checkType := IsEmpty(std.Stype)
	checkTime := IsEmpty(std.Entime)
	checkSession := IsEmpty(std.Stime)

	if checkname || checkphone || checkLvl || checkType || checkTime || checkSession {
		templates.ExecuteTemplate(w, "msg.html", "Failed! Some fields are empty")
		return
	}
	EditStd(std)
	http.Redirect(w, r, "admin_disp", 302)
}
