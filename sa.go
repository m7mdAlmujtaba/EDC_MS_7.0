package main

import (
	"net/http"
)


func SaReportGetHandler(w http.ResponseWriter, r *http.Request){
	// Here Goes The Session Check
	var cnts = []Counts{}
	cnts = SaStat()
	err = templates.ExecuteTemplate(w, "sa_stat.html", cnts)
	CheckErr(err)
}

func SaAttendGetHandler(w http.ResponseWriter, r *http.Request){
	// Here Goes The Session Check
	var stds = []Student{}
	stds = AllStdsExcel()
	err = templates.ExecuteTemplate(w, "sa_excel.html", stds)
	CheckErr(err)
}

func SaAttendPostHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	std_type :=  r.PostForm.Get("type")
	sess :=  r.PostForm.Get("session")
	lvl := r.PostForm.Get("level")
	tim := r.PostForm.Get("time")
	var stds = []Student{}
	stds = StdsExcel(std_type, sess, lvl, tim)
	err = templates.ExecuteTemplate(w, "sa_excel.html", stds)
	CheckErr(err)
}