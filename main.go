package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var err error
var templates *template.Template
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func DbConn() (db *sql.DB) {
	db, err = sql.Open("mysql", "root:edc_pass_321@tcp(127.0.0.1:3306)/edc_db?parseTime=true")
	CheckErr(err)
	return db
}

// No cache handler
var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	// Parse Templates Folder
	templates = template.Must(template.ParseGlob("templates/*.html"))

	Route()
	fmt.Println("Server running on port :8000")
	http.ListenAndServe(":8000", nil)
}

func Route() {
	r := mux.NewRouter()

	// General
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")

	r.HandleFunc("/admin_dashboard", AuthanticatedAdmin(adminHomeHandler)).Methods("GET")
	r.HandleFunc("/reg_dashboard", AuthanticatedRegistrar(registrarHomeHandler)).Methods("GET")
	r.HandleFunc("/sa_dashboard", AuthanticatedAffair(affairHomeHandler)).Methods("GET")
	r.HandleFunc("/se_dashboard", AuthanticatedSecurity(securityHomeHandler)).Methods("GET")

	// Display Students
	r.HandleFunc("/ajax", AuthanticatedUser(AjaxGetHandler)).Methods("GET")
	r.HandleFunc("/disp", AuthanticatedUser(DispGetHandler)).Methods("GET")

	// Admin
	r.HandleFunc("/admin_records", AuthanticatedAdmin(recordsGetHandler)).Methods("GET")
	r.HandleFunc("/reportstbl", AuthanticatedAdmin(recordsPostHandler)).Methods("GET")
	r.HandleFunc("/lvlcount", AuthanticatedAdmin(LvlCountHandler)).Methods("GET")

	r.HandleFunc("/admin_fees", AuthanticatedAdmin(feesGetHandler)).Methods("GET")
	r.HandleFunc("/admin_fees", AuthanticatedAdmin(feesPostHandler)).Methods("POST")

	r.HandleFunc("/admin_disp", AuthanticatedAdmin(adminDispGetHandler)).Methods("GET")
	r.HandleFunc("/get_stds", AuthanticatedAdmin(getStds)).Methods("GET")
	r.HandleFunc("/edit_std", AuthanticatedAdmin(editStd)).Methods("POST")

	r.HandleFunc("/admin_users", AuthanticatedAdmin(usersGetHandler)).Methods("GET")
	r.HandleFunc("/admin_users", AuthanticatedAdmin(usersPostHandler)).Methods("POST")
	r.HandleFunc("/remove", AuthanticatedAdmin(delUsersPostHandler)).Methods("POST")
	r.HandleFunc("/edit_pass", AuthanticatedAdmin(editpassHandler)).Methods("POST")

	r.HandleFunc("/admin_report", AuthanticatedAdmin(reportsGetHandler)).Methods("GET")
	r.HandleFunc("/admin_acts", AuthanticatedAdmin(actionsGetHandler)).Methods("GET")
	r.HandleFunc("/ac", AuthanticatedAdmin(ActionsTableGetHandler)).Methods("GET")

	// Registrar
	r.HandleFunc("/reg_pt", AuthanticatedRegistrar(PtGetHandler)).Methods("GET")
	r.HandleFunc("/reg_pt", AuthanticatedRegistrar(PtPostHandler)).Methods("POST")

	r.HandleFunc("/en", AuthanticatedRegistrar(EnTableGetHandler)).Methods("GET")
	r.HandleFunc("/reg_en", AuthanticatedRegistrar(EnGetHandler)).Methods("GET")
	r.HandleFunc("/reg_en", AuthanticatedRegistrar(EnPostHandler)).Methods("POST")

	r.HandleFunc("/cert", AuthanticatedRegistrar(CertTableGetHandler)).Methods("GET")
	r.HandleFunc("/reg_cert", AuthanticatedRegistrar(CertGetHandler)).Methods("GET")
	r.HandleFunc("/reg_cert", AuthanticatedRegistrar(CertPostHandler)).Methods("POST")

	r.HandleFunc("/reg_report", AuthanticatedRegistrar(ReportGetHandler)).Methods("GET")

	// Student's Affair
	r.HandleFunc("/lvl", AuthanticatedAffair(Lvlhandler)).Methods("GET")
	r.HandleFunc("/lvls", AuthanticatedAffair(Lvlshandler)).Methods("GET")

	r.HandleFunc("/excel", AuthanticatedAffair(ExcelTableGetHandler)).Methods("GET")
	r.HandleFunc("/sa_attend", AuthanticatedAffair(SaAttendGetHandler)).Methods("GET")
	r.HandleFunc("/sa_attend", AuthanticatedAffair(SaAttendPostHandler)).Methods("POST")

	r.HandleFunc("/sa_stat", AuthanticatedAffair(SaReportGetHandler)).Methods("GET")

	/*r.HandleFunc("/std_img_page", AuthanticatedAffair(stdsImgPagetHandler)).Methods("GET")
	r.HandleFunc("/std_img", AuthanticatedAffair(getstdImgPage)).Methods("GET")*/

	r.HandleFunc("/profile", AuthanticatedSecurity(profilesPageHandler)).Methods("GET")
	r.HandleFunc("/stds_profile", AuthanticatedSecurity(profilesGetHandler)).Methods("GET")
	r.HandleFunc("/stds_profile_upload", AuthanticatedSecurity(uploadProfilePicture)).Methods("POST")

	r.HandleFunc("/cards", AuthanticatedSecurity(cardsPageHandler)).Methods("GET")
	r.HandleFunc("/stds_cards_page", AuthanticatedSecurity(stdsCardsGetHandler)).Methods("GET")
	r.HandleFunc("/cardshtml", AuthanticatedSecurity(cardsHtmlPageHandler)).Methods("GET")

	// Serve Static Files
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	fu := http.FileServer(http.Dir("./uploads"))
	r.PathPrefix("/uploads").Handler(NoCache(http.StripPrefix("/uploads", fu)))

	http.Handle("/", r)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	user_pass := r.PostForm.Get("pass")
	Users := QueryUser(username)
	//userid = Users.
	//usr = Users.Uname
	if user_pass == Users.Upass {
		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		//session.Values["id"] = userid
		session.Values["authenticated"] = true
		session.Values["type"] = Users.Utype
		session.Values["id"] = Users.ID
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch Users.Utype {
		case "Admin":
			// fmt.Println(users.uid, users.uname, users.upass, users.utype)
			http.Redirect(w, r, "admin_dashboard", 302)
		case "Registrar":
			// fmt.Println(users.uid, users.uname, users.upass, users.utype)
			http.Redirect(w, r, "reg_dashboard", 302)
		case "Affair":
			// fmt.Println(users.uid, users.uname, users.upass, users.utype)
			http.Redirect(w, r, "sa_dashboard", 302)
		case "Security":
			// fmt.Println(users.uid, users.uname, users.upass, users.utype)
			http.Redirect(w, r, "se_dashboard", 302)
		}

	} else {
		templates.ExecuteTemplate(w, "login.html", "Invalid Username or Password")
		//http.Redirect(w, r, "/login", 302)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		//http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 302)
		return
	}
	switch session.Values["type"] {
	case "Admin":
		http.Redirect(w, r, "/admin_dashboard", 302)
		return
	case "Registrar":
		http.Redirect(w, r, "/reg_dashboard", 302)
		return
	case "Affair":
		http.Redirect(w, r, "/sa_dashboard", 302)
		return
	case "Security":
		http.Redirect(w, r, "/se_dashboard", 302)
		return
	default:
		http.Redirect(w, r, "/login", 302)
		return
	}

}

func AuthanticatedUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			//templates.ExecuteTemplate(w, "login.html", untyped)
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func AuthanticatedAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			//templates.ExecuteTemplate(w, "login.html", untyped)
			http.Redirect(w, r, "/login", 302)
			return
		}
		switch session.Values["type"] {
		case "Admin":
			handler.ServeHTTP(w, r)
			return
		case "Registrar":
			http.Redirect(w, r, "/reg_dashboard", 302)
			return
		case "Affair":
			http.Redirect(w, r, "/sa_dashboard", 302)
			return
		case "Security":
			http.Redirect(w, r, "/se_dashboard", 302)
			return
		}

	}
}

func AuthanticatedRegistrar(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/login", 302)
			return
		}
		switch session.Values["type"] {
		case "Admin":
			http.Redirect(w, r, "/admin_dashboard", 302)
			return
		case "Registrar":
			handler.ServeHTTP(w, r)
			return
		case "Affair":
			http.Redirect(w, r, "/sa_dashboard", 302)
			return
		case "Security":
			http.Redirect(w, r, "/se_dashboard", 302)
			return
		}

	}
}

func AuthanticatedAffair(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/login", 302)
			return
		}
		switch session.Values["type"] {
		case "Admin":
			http.Redirect(w, r, "/admin_dashboard", 302)
			return
		case "Registrar":
			http.Redirect(w, r, "/reg_dashboard", 302)
			return
		case "Affair":
			handler.ServeHTTP(w, r)
			return
		case "Security":
			http.Redirect(w, r, "/se_dashboard", 302)
			return
		}
	}
}

func AuthanticatedSecurity(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/login", 302)
			return
		}
		switch session.Values["type"] {
		case "Admin":
			http.Redirect(w, r, "/admin_dashboard", 302)
			return
		case "Registrar":
			http.Redirect(w, r, "/reg_dashboard", 302)
			return
		case "Affair":
			http.Redirect(w, r, "/sa_dashboard", 302)
			return
		case "Security":
			handler.ServeHTTP(w, r)
			return
		}
	}
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	User := session.Values["username"]
	templates.ExecuteTemplate(w, "admin_dashboard.html", User)
}

func registrarHomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	User := session.Values["username"]
	templates.ExecuteTemplate(w, "reg_dashboard.html", User)
}

func affairHomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	User := session.Values["username"]
	templates.ExecuteTemplate(w, "sa_dashboard.html", User)

}

func securityHomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	User := session.Values["username"]
	templates.ExecuteTemplate(w, "se_dashboard.html", User)

}

func CheckErr(err error) {
	file, e := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalln("Faild")
	}
	log.SetOutput(file)
	if err != nil {
		panic(err)
		log.Println(err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = nil
	session.Values["authenticated"] = false
	session.Values["type"] = nil
	session.Values["id"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)

}
