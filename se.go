package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Get Students Profiles
func profilesPageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "stds_profile.html", nil)
}

// Generate Pagination for Display Students For Admin
func profilesGetHandler(w http.ResponseWriter, r *http.Request) {

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	if len(pg) < 1 {
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1 {
		page = 1
	}

	if search != "" {
		query += "WHERE id LIKE '%" + search + "%' or std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE id LIKE '%" + search + "%' or std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page - 1) * limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)

	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))

	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Image</th>
		<th scope="col">Student Name</th>
		<th scope="col">Level</th>
		<th scope="col">Edit</th>
	  </tr>
	</thead>
	<tbody>`

	if count > 0 {
		for _, i := range stds {
			if i.Pic_path == "" {
				i.Pic_path = "../uploads/avatar.jpg"
			}
			html += `
			<tr>
				
				<td>` + strconv.Itoa(i.id) + `</td>
				<td> <img src="` + i.Pic_path + `" style='width:50px; height:auto;' class="profile-pic"> </td>
				<td>` + i.Sname + `</td>
				<td>` + i.Slvl + `</td>
				<td> <button class="popup btn btn--radius-2 btn--red m-r-55" data-id="` + strconv.Itoa(i.id) + `" data-pic= " ` + i.Pic_path + `" data-name= " ` + i.Sname + `"  data-phone= " ` + i.Sphone + `"  data-toggle="modal" data-target="#modal"> <i class="fa fa-pencil-square-o" aria-hidden="true"></i> Add/Edit Image </button> </td>
				</td>

			</tr>`
		}
	} else {
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
			for cnt := 1; cnt <= 5; cnt++ {
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		} else {
			end_limit := total_pages - 5
			if page > end_limit {
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++ {
					pages_array = append(pages_array, cnt)
				}

			} else {
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page - 1; cnt <= page+1; cnt++ {
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	} else {
		for cnt := 1; cnt <= total_pages; cnt++ {
			pages_array = append(pages_array, cnt)
		}
	}
	// fmt.Println(pages_array)

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++ {
		if page == pages_array[cnt] {
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">` + strconv.Itoa(x) + `
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			`
			previous_id := x - 1
			if previous_id > 0 {
				previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(previous_id) + `">Previous</a></li>`
			} else {
				previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1
			if next_id >= total_pages {
				next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			} else {
				next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(next_id) + `">Next</a></li>`
			}

		} else {
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			} else {
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(pages_array[cnt].(int)) + `">` + strconv.Itoa(pages_array[cnt].(int)) + `</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w, html)
	//templates.ExecuteTemplate(w, "ajax.html", stds)
}

type StudentData struct {
	Quote string
}

func uploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	img_data := r.FormValue("image_data")
	id := r.FormValue("std_id")

	base64ToImage(img_data, id)

	http.Redirect(w, r, "profile", 302)
}

func base64ToImage(imgdata, id string) {
	db := DbConn()
	var f *os.File
	image := imgdata
	coI := strings.Index(string(image), ",")
	rawImage := string(image)[coI+1:]
	unbased, _ := base64.StdEncoding.DecodeString(rawImage)
	res := bytes.NewReader(unbased)
	var imgpath string
	switch strings.TrimSuffix(image[5:coI], ";base64") {
	case "image/png":
		pngI, err := png.Decode(res)
		CheckErr(err)
		f, _ = os.OpenFile("uploads/"+id+".png", os.O_WRONLY|os.O_CREATE, 0777)
		png.Encode(f, pngI)
		imgpath = "uploads/" + id + ".png"
		fmt.Println("generated", imgpath)
	case "image/jpeg":
		jpgI, err := jpeg.Decode(res)
		CheckErr(err)
		f, _ = os.OpenFile("uploads/"+id+".jpg", os.O_WRONLY|os.O_CREATE, 0777)
		jpeg.Encode(f, jpgI, &jpeg.Options{Quality: 75})
		imgpath = "uploads/" + id + ".jpg"
		fmt.Println("generated", imgpath)
	}

	row, err := db.Query("UPDATE `edc_db`.`students` SET `std_pic` = ? WHERE `id` = ?", imgpath, id)
	fmt.Println("image path is added to db")
	CheckErr(err)
	row.Close()
	defer db.Close()
}

func cardsPageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "se_cards.html", nil)
}

// Generate Pagination for Display Students For Admin
func stdsCardsGetHandler(w http.ResponseWriter, r *http.Request) {

	query := "SELECT * FROM students "
	row_count := "SELECT count(*) FROM students "
	pg := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	if len(pg) < 1 {
		pg = "1"
	}
	page, err := strconv.Atoi(pg)
	CheckErr(err)
	if err != nil || page < 1 {
		page = 1
	}

	if search != "" {
		query += "WHERE id LIKE '%" + search + "%' or std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "

		row_count += "WHERE id LIKE '%" + search + "%' or std_id LIKE '%" + search + "%' or std_name LIKE '%" + search + "%' "
	}

	start := (page - 1) * limit

	query += fmt.Sprintf("ORDER BY id DESC LIMIT %d, %d", start, limit)

	stds := GetStudents(query)
	count := GetTotalRows(row_count)

	total_pages := int(math.Ceil(float64(count) / float64(limit)))

	html := `<table id='myTable' class='table table-dark'>
	<thead>
	  <tr>
		<th scope="col">id</th>
		<th scope="col">Student Image</th>
		<th scope="col">Student Name</th>
		<th scope="col">Level</th>
		<th scope="col">Add to Cards</th>
	  </tr>
	</thead>
	<tbody>`

	if count > 0 {
		for _, i := range stds {
			if i.Pic_path == "" {
				i.Pic_path = "../uploads/avatar.jpg"
			}
			html += `
			<tr>
				
				<td>` + strconv.Itoa(i.id) + `</td>
				<td> <img src="` + i.Pic_path + `" style='width:50px; height:auto;' class="profile-pic"> </td>
				<td>` + i.Sname + `</td>
				<td>` + i.Slvl + `</td>
				<td> <button class="addstd-btn btn btn--radius-2 btn--red m-r-55" data-json='{"id": "` + strconv.Itoa(i.id) + `", "name": "` + i.Sname + `", "level": "` + i.Slvl + `", "time":"` + i.Entime + `", "date":"` + i.D1 + `", "imgpath":"` + i.Pic_path + `"}' onclick="addStdtoCards(this)"> <i class="fa fa-plus" aria-hidden="true"></i> Add to Cards </button> </td>
				</td>

			</tr>`
		}
	} else {
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
			for cnt := 1; cnt <= 5; cnt++ {
				pages_array = append(pages_array, cnt)
			}
			pages_array = append(pages_array, "...")
			pages_array = append(pages_array, total_pages)
		} else {
			end_limit := total_pages - 5
			if page > end_limit {
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := end_limit; cnt <= total_pages; cnt++ {
					pages_array = append(pages_array, cnt)
				}

			} else {
				pages_array = append(pages_array, 1)
				pages_array = append(pages_array, "...")
				for cnt := page - 1; cnt <= page+1; cnt++ {
					pages_array = append(pages_array, cnt)
				}
				pages_array = append(pages_array, "...")
				pages_array = append(pages_array, total_pages)
			}

		}
	} else {
		for cnt := 1; cnt <= total_pages; cnt++ {
			pages_array = append(pages_array, cnt)
		}
	}
	// fmt.Println(pages_array)

	// Links
	for cnt := 0; cnt < len(pages_array); cnt++ {
		if page == pages_array[cnt] {
			x := pages_array[cnt].(int)
			page_link += `
			<li class="page-item active">
				<a class="page-link" href="#">` + strconv.Itoa(x) + `
					<span class="sr-only">(Current)</span>
				</a>
			</li>
			`
			previous_id := x - 1
			if previous_id > 0 {
				previous_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(previous_id) + `">Previous</a></li>`
			} else {
				previous_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Previous</a>
			  </li>
			  `
			}
			next_id := x + 1
			if next_id >= total_pages {
				next_link = `
			  <li class="page-item disabled">
				<a class="page-link" href="#">Next</a>
			  </li>
				`
			} else {
				next_link = `<li class="page-item"><a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(next_id) + `">Next</a></li>`
			}

		} else {
			if pages_array[cnt] == "..." {
				page_link += `<li class="page-item disabled">
				<a class="page-link" href="#">...</a>
			  </li>
			  `
			} else {
				page_link += `
					<li class="page-item">
						<a class="page-link" href="javascript:void(0)" data-page_number="` + strconv.Itoa(pages_array[cnt].(int)) + `">` + strconv.Itoa(pages_array[cnt].(int)) + `</a>
					</li>
				`
			}
		}
	}

	html += previous_link + page_link + next_link
	fmt.Fprintf(w, html)
	//templates.ExecuteTemplate(w, "ajax.html", stds)
}

func cardsHtmlPageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "cards_page.html", nil)
}
