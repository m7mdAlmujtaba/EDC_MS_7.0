package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"math"
)
var limit int = 100

func GetTotalRows(cstr string) int{
	var rows int

	db := DbConn()
	row_count, err := db.Query(cstr)
	CheckErr(err)
	for row_count.Next(){
		err = row_count.Scan(&rows)
		CheckErr(err)
	}
	
	defer db.Close()
	return rows
}

func GetStudents(qstr string) []Student{
	db := DbConn()
	//std_row, err := db.Query("SELECT * FROM students LIMIT ?,?", i, j)
	std_row, err := db.Query(qstr)
	CheckErr(err)
	std := Student{}
	stds := []Student{}
	var id int
	for std_row.Next(){
		err = std_row.Scan(&id, &std.Sid , &std.Sname, &std.Sphone, &std.Slvl, &std.Sen_date,  &std.Re_date, &std.Stype, &std.Entime, &std.Stime) 
		CheckErr(err)
		d1 := std.Sen_date.Format("2006-01-02")
		d2 := std.Re_date.Format("2006-01-02")
		if d1 == "1000-01-01"{
			std.D1 = ""
		}else{
			std.D1 = d1
		}
		if d2 == "1000-01-01"{
			std.D2 = ""
		}else{
			std.D2 = d2
		}
		stds = append(stds, std)
	}
	defer db.Close()
	return stds
}

func GetActions(qstr string) []Action{
	db := DbConn()
	//std_row, err := db.Query("SELECT * FROM students LIMIT ?,?", i, j)
	act_row, err := db.Query(qstr)
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

func DispGetHandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "disp.html", nil)
}

func adminDispGetHandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "admin_disp.html", nil)
}

// Generate Pagination for Display Students
func AjaxGetHandler(w http.ResponseWriter, r *http.Request){

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")


	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}


	if search != ""{
		query += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page-1)*limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)



	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))

	/*
	fmt.Println(query)
	fmt.Println("Page Number: " , page )
	fmt.Println("Number of Rows: ", count)
	fmt.Println("Total Number of Pages: ", total_pages)

	for _,std := range stds{
		fmt.Println(std)
	}
	*/
	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Name</th>
		<th scope="col">Phone</th>
		<th scope="col">Level</th>
		<th scope="col">First Enrollment</th>
		<th scope="col">Last Registration</th>
		<th scope="col">Type</th>
		<th scope="col">Time</th>
		<th scope="col">Session</th>
	  </tr>
	</thead>
	<tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Sid) +`</td>
				<td>` + i.Sname +`</td>
				<td>` + i.Sphone +`</td>
				<td>` + i.Slvl +`</td>
				<td>` + i.D1 +`</td>
				<td>` + i.D2 +`</td>
				<td>` + i.Stype +`</td>
				<td>` + i.Entime +`</td>
				<td>` + i.Stime +`</td>

			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	
	// fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
	//templates.ExecuteTemplate(w, "ajax.html", stds)
}

// Generate Pagination for Display Students For Admin
func getStds(w http.ResponseWriter, r *http.Request){

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")


	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}


	if search != ""{
		query += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page-1)*limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)



	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))
/*
	fmt.Println(query)
	fmt.Println("Page Number: " , page )
	fmt.Println("Number of Rows: ", count)
	fmt.Println("Total Number of Pages: ", total_pages)

	for _,std := range stds{
		fmt.Println(std)
	}
*/
	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Name</th>
		<th scope="col">Phone</th>
		<th scope="col">Level</th>
		<th scope="col">Type</th>
		<th scope="col">Time</th>
		<th scope="col">Session</th>
		<th scope="col">Edit</th>
	  </tr>
	</thead>
	<tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Sid) +`</td>
				<td>` + i.Sname +`</td>
				<td>` + i.Sphone +`</td>
				<td>` + i.Slvl +`</td>
				<td>` + i.Stype +`</td>
				<td>` + i.Entime +`</td>
				<td>` + i.Stime +`</td>
				<td> <button class="popup btn btn--radius-2 btn--red m-r-55" data-id=`+ strconv.Itoa(i.Sid) +` data-name= " ` + i.Sname +`"  data-phone= " `+ i.Sphone +`"  data-toggle="modal" data-target="#modalForm"> <i class="fa fa-pencil-square-o" aria-hidden="true"></i> Edit </button> </td>

			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	// fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
	//templates.ExecuteTemplate(w, "ajax.html", stds)
}
// Generate Pagination for Students Enrollment
func EnTableGetHandler(w http.ResponseWriter, r *http.Request){

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")


	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}


	if search != ""{
		query += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page-1)*limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)



	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))

	/*
	fmt.Println(query)
	fmt.Println("Page Number: " , page )
	fmt.Println("Number of Rows: ", count)
	fmt.Println("Total Number of Pages: ", total_pages)

	for _,std := range stds{
		fmt.Println(std)
	}
*/
	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Name</th>
		<th scope="col">Level</th>
		<th scope="col">Form</th>
	  </tr>
	</thead>
	<tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Sid) +`</td>
				<td>` + i.Sname +`</td>
				<td>` + i.Slvl +`</td>
				<td> <button class="popup btn btn--radius-2 btn--red m-r-55" data-id="`+ strconv.Itoa(i.Sid) +`" data-toggle="modal" data-target="#modalForm"><i class="fa fa-user-plus" aria-hidden="true"></i>  Enroll </button> </td>

			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	//fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
}

// Generate Pagination for Students Certificate, Statements and Freeze 
func CertTableGetHandler(w http.ResponseWriter, r *http.Request){
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")


	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}


	if search != ""{
		query += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page-1)*limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)



	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))
/*
	fmt.Println(query)
	fmt.Println("Page Number: " , page )
	fmt.Println("Number of Rows: ", count)
	fmt.Println("Total Number of Pages: ", total_pages)

	for _,std := range stds{
		fmt.Println(std)
	}
*/
	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Name</th>
		<th scope="col">Level</th>
		<th scope="col">Submit</th>
	  </tr>
	</thead>
	<tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Sid) +`</td>
				<td>` + i.Sname +`</td>
				<td>` + i.Slvl +`</td>
				<td> <button class="popup btn btn--radius-2 btn--red m-r-55" data-id="`+ strconv.Itoa(i.Sid) +`" data-toggle="modal" data-target="#modalForm"> Statement </button> </td>

			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	//fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
}

func ActionsTableGetHandler(w http.ResponseWriter, r *http.Request){
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")


	query := "SELECT * FROM actions "
	row_count := "SELECT count(*) FROM actions "
	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}


	if search != ""{
		query += "WHERE user_name LIKE '%" + search + "%' or receipt_id LIKE '%" + search + "%' "

		row_count += "WHERE user_name LIKE '%" + search + "%' or receipt_id LIKE '%" + search + "%' "
	}

	start := (page-1)*limit

	query += fmt.Sprintf("ORDER BY act_id DESC LIMIT %d, %d", start, limit)



	stds := GetActions(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))
/*
	for _,std := range stds{
		fmt.Println(std)
	}
*/
	html := `
	<table id="myTable" class="table table-dark">
	  <thead>
		<tr>
		  <th scope="col">Receipt id</th>
		  <th scope="col">Username</th>
		  <th scope="col">Type</th>
		  <th scope="col">Student Name</th>
		  <th scope="col">Fees</th>
		  <th scope="col">Action Date</th>
		  <th scope="col">Time</th>
		</tr>
	  </thead>
	  <tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Aid) +`</td>
				<td>` + i.Auser +`</td>
				<td>` + i.Atype +`</td>
				<td>` + i.Asname +`</td>
				<td>` + strconv.Itoa(i.Fee) +`</td>
				<td>` + i.DAdate +`</td>
				<td>` + i.DDate +`</td>
			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	//fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
}

// Generate Pagination for Students Affair to Export to Excel
func ExcelTableGetHandler(w http.ResponseWriter, r *http.Request){
	pg := r.URL.Query().Get("page")
	level := r.URL.Query().Get("level")
	d1 := r.URL.Query().Get("from")
	d2 := r.URL.Query().Get("to")
	s_type := r.URL.Query().Get("stype")
	s_session := r.URL.Query().Get("ssession")
	stime := r.URL.Query().Get("time")

	fmt.Println(d1)

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "

	// Handle page
	if len(pg) < 1{
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1{
		page = 1
	}

	// Handle level
	if level == "non"{
		query += "where en_date = '1000-01-01' "
		row_count += "where en_date = '1000-01-01' "
	}else{
		switch level {
		case "p0":
			query += "where std_lvl = 'Pre 0' "
			row_count += "where std_lvl = 'Pre 0' "
		case "p1":
			query += "where std_lvl = 'Pre 1' "
			row_count += "where std_lvl = 'Pre 1' "
		case "p2":
			query += "where std_lvl = 'Pre 2' "
			row_count += "where std_lvl = 'Pre 2' "
		case "l1":
			query += "where std_lvl = 'Level 1' "
			row_count += "where std_lvl = 'Level 1' "
		case "l2":
			query += "where std_lvl = 'Level 2' "
			row_count += "where std_lvl = 'Level 2' "
		case "l3":
			query += "where std_lvl = 'Level 3' "
			row_count += "where std_lvl = 'Level 3' "
		case "l4":
			query += "where std_lvl = 'Level 4' "
			row_count += "where std_lvl = 'Level 4' "
		case "l5":
			query += "where std_lvl = 'Level 5' "
			row_count += "where std_lvl = 'Level 5' "
		case "l6":
			query += "where std_lvl = 'Level 6' "
			row_count += "where std_lvl = 'Level 6' "
		case "l7":
			query += "where std_lvl = 'Level 7' "
			row_count += "where std_lvl = 'Level 7' "
		case "l8":
			query += "where std_lvl = 'Level 8' "
			row_count += "where std_lvl = 'Level 8' "
		}
		
		// Handle type, session and time
		query += "and en_type = '" + s_type + "' and std_time = '" + s_session + "' and time = '" + stime + "' "
		row_count += "and en_type = '" + s_type + "' and std_time = '" + s_session + "' and time = '" + stime + "' "
	}
	// Handle date
	if d1 == ""{
		d1 = time.Now().Format("2006-01-02")
	}
	if d2 == ""{
		d2 = time.Now().Format("2006-01-02")
	}
	date1, _ := time.Parse("2006-01-2", d1)
	date2, _ := time.Parse("2006-01-2", d2)

	stds := []Student{}
	var count int
	if date1.After(date2){
		query += fmt.Sprintf("and date(last_reg) between '%04d-%02d-%02d' and '%04d-%02d-%02d' ", date2.Year(),date2.Month(), date2.Day(), date1.Year(),date1.Month(), date1.Day())
		row_count += fmt.Sprintf("and date(last_reg) between '%04d-%02d-%02d' and '%04d-%02d-%02d' ",date2.Year(),date2.Month(), date2.Day(), date1.Year(),date1.Month(), date1.Day())
	}else{
		query += fmt.Sprintf("and date(last_reg) between '%04d-%02d-%02d' and '%04d-%02d-%02d' ", date1.Year(),date1.Month(), date1.Day(),  date2.Year(),date2.Month(), date2.Day())
		row_count += fmt.Sprintf("and date(last_reg) between '%04d-%02d-%02d' and '%04d-%02d-%02d' ", date1.Year(),date1.Month(), date1.Day(),  date2.Year(),date2.Month(), date2.Day())
	}
	

	start := (page-1)*limit
	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)

	fmt.Println(query)

	stds = GetStudents(query)
	count = GetTotalRows(row_count)

	

	total_pages := int(math.Ceil(float64(count) / float64(limit)))

	fmt.Println(query)
	fmt.Println("Page Number: " , page )
	fmt.Println("Number of Rows: ", count)
	fmt.Println("Total Number of Pages: ", total_pages)

	for _,std := range stds{
		fmt.Println(std)
	}

	html := `<table id="myTable" class="order-table table table-dark">
	<thead>
	  <tr>
		  <th scope="col">id</th>
		  <th scope="col">Student Name</th>
		  <th scope="col">Phone</th>
		  <th scope="col"> Select All <input type="checkbox" class="" onclick="checkAll(this)"> </th>
	 
		</tr>
	</thead>
	<tbody>`
	if count > 0 {
		for _,i := range stds{
			html += `
			<tr>
				<td>` + strconv.Itoa(i.Sid) +`</td>
				<td>` + i.Sname +`</td>
				<td>` + i.Sphone +`</td>
				<td><input type="checkbox"/></td>
			</tr>`
		}	
	}else{
		html += `
			<tr>
				<td>No Data Found ..</td>
			</tr>
		`
	}
	html += `
	</tbody>
	</table> <br>
	<div align="center">
	<ul class="pagination">
	`
	
	previous_link := ``
	next_link := ``
	page_link := ``
	var pages_array []interface{}

	if total_pages > 4 {
		if page < 5 {
			for cnt := 1; cnt <= 5; cnt++{
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		}else{
			end_limit := total_pages - 5;
			if page > end_limit{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++{
					pages_array = append(pages_array, cnt)
				}

			}else{
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page-1; cnt <= page+1; cnt++{
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	}else{
		for cnt := 1; cnt <= total_pages; cnt++{
			pages_array = append(pages_array, cnt)
		}
	}
	fmt.Println(pages_array)
	

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++{
		if page == pages_array[cnt]{
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">`+ strconv.Itoa(x) +`
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			` 
			previous_id :=  x - 1
			if previous_id > 0{
			  previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(previous_id) + `">Previous</a></li>`
			}else{
			  previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1;
			if next_id >= total_pages{
			  next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			}else{
			  next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(next_id) + `">Next</a></li>`
			}

		}else{
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			}else{
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="`+ strconv.Itoa(pages_array[cnt].(int)) +`">`+ strconv.Itoa(pages_array[cnt].(int)) +`</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w,html)
}

type Level struct{
	Ln, Lp0, Lp1, Lp2, Ll1, Ll2, Ll3, Ll4, Ll5, Ll6, Ll7, Ll8 int
}

func Lvlshandler(w http.ResponseWriter, r *http.Request){
	var lvl Level
	query := "SELECT count(*) FROM edc_db.students where year(last_reg) = year(NOW()) and month(last_reg) = month(NOW()) "
	lvl.Ln = GetTotalRows(query +"and en_date = '1000-01-01'")
	lvl.Lp0 = GetTotalRows(query +"and std_lvl = 'Pre 0'")
	lvl.Lp1 = GetTotalRows(query +"and std_lvl = 'Pre 1'")
	lvl.Lp2 = GetTotalRows(query +"and std_lvl = 'Pre 2'")
	lvl.Ll1 = GetTotalRows(query +"and std_lvl = 'Level 1'")
	lvl.Ll2 = GetTotalRows(query +"and std_lvl = 'Level 2'")
	lvl.Ll3 = GetTotalRows(query +"and std_lvl = 'Level 3'")
	lvl.Ll4 = GetTotalRows(query +"and std_lvl = 'Level 4'")
	lvl.Ll5 = GetTotalRows(query +"and std_lvl = 'Level 5'")
	lvl.Ll6 = GetTotalRows(query +"and std_lvl = 'Level 6'")
	lvl.Ll7 = GetTotalRows(query +"and std_lvl = 'Level 7'")
	lvl.Ll8 = GetTotalRows(query +"and std_lvl = 'Level 8'")
	templates.ExecuteTemplate(w, "lvls.html", lvl)
}

func Lvlhandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "sa_excel.html", nil)
}