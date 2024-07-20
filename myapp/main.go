package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    _ "github.com/lib/pq"
    "html/template"
)

var db *sql.DB

func init() {
    var err error
    connStr := "user=postgres password=new_password dbname=candidate_db sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./static/index.html")
})

    http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/ask.html")
    })
    http.HandleFunc("/role", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            role := r.FormValue("role")
            if role == "candidate" {
                http.Redirect(w, r, "/register_candidate", http.StatusSeeOther)
            } else if role == "recruiter" {
                http.Redirect(w, r, "/register_recruiter", http.StatusSeeOther)
            } else {
                http.Error(w, "Invalid role", http.StatusBadRequest)
            }
        } else {
            http.ServeFile(w, r, "./static/ask.html")
        }
    })

    http.HandleFunc("/register_candidate", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/register_candidate.html")
    })

    http.HandleFunc("/register_recruiter", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/register_recruiter.html")
    })

    http.HandleFunc("/submit_candidate_registration", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            name := r.FormValue("name")
            email := r.FormValue("email")
            password := r.FormValue("password")
            dob := r.FormValue("dob")
            qualification := r.FormValue("qualification")
            github := r.FormValue("github")
            mobile := r.FormValue("mobile")

            _, err := db.Exec(`
                INSERT INTO candidate11 (name, email, password, dob, qualification, github, mobile)
                VALUES ($1, $2, $3, $4, $5, $6, $7)`,
                name, email, password, dob, qualification, github, mobile)

            if err != nil {
                log.Println("Error inserting data:", err)
                http.Error(w, "Unable to save data", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/?name="+name, http.StatusSeeOther)
        } else {
            http.ServeFile(w, r, "./static/register_candidate.html")
        }
    })

    http.HandleFunc("/submit_recruiter_registration", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            name := r.FormValue("name")
            email := r.FormValue("email")
            password := r.FormValue("password")
            post := r.FormValue("post")
            company := r.FormValue("company")
            experience := r.FormValue("experience")
            branch := r.FormValue("branch")
            mobile := r.FormValue("mobile")

            _, err := db.Exec(`
                INSERT INTO recruiter (name, email, password, post, company, experience, branch, mobile)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
                name, email, password, post, company, experience, branch, mobile)

            if err != nil {
                log.Println("Error inserting data:", err)
                http.Error(w, "Unable to save data", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/?name="+name, http.StatusSeeOther)
        } else {
            http.ServeFile(w, r, "./static/register_recruiter.html")
        }
    })

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            password := r.FormValue("password")
    
            var dbPassword, dbName, dbEmail, dbRole string
            err := db.QueryRow(`
                SELECT name, email, password, 'candidate' AS role FROM candidate11 WHERE email=$1
                UNION
                SELECT name, email, password, 'recruiter' AS role FROM recruiter WHERE email=$1
                UNION
                SELECT name, email, password, 'admin' AS role FROM admin WHERE email=$1
            `, username).Scan(&dbName, &dbEmail, &dbPassword, &dbRole)
            
            if err != nil {
                if err == sql.ErrNoRows {
                    http.Error(w, "Invalid username or password", http.StatusUnauthorized)
                } else {
                    log.Println("Error querying database:", err)
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                }
                return
            }
    
            if dbPassword == password {
                http.SetCookie(w, &http.Cookie{
                    Name:  "username",
                    Value: dbEmail,
                    Path:  "/",
                })
                
                // Redirect to the appropriate dashboard based on the role
                switch dbRole {
                case "candidate":
                    http.Redirect(w, r, "/candidate_dashboard", http.StatusSeeOther)
                case "recruiter":
                    http.Redirect(w, r, "/recruiter_dashboard", http.StatusSeeOther)
                case "admin":
                    http.Redirect(w, r, "/admin_dashboard", http.StatusSeeOther)
                default:
                    http.Error(w, "Invalid role", http.StatusUnauthorized)
                }
                return
            } else {
                http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            }
        }
    
        name := r.URL.Query().Get("name")
        if name != "" {
            fmt.Fprintf(w, "User %s successfully registered! Please log in.\n", name)
        }
        http.ServeFile(w, r, "./static/login.html")
    })
    
    http.HandleFunc("/candidate_dashboard", func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("username")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            log.Println("Error retrieving cookie:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        username := cookie.Value
        var name, email, dob, qualification, github, mobile string

        err = db.QueryRow("SELECT name, email, dob, qualification, github, mobile FROM candidate11 WHERE email=$1", username).Scan(&name, &email, &dob, &qualification, &github, &mobile)
        if err != nil {
            log.Println("Error querying candidate details:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        data := struct {
            Name         string
            Email        string
            Dob          string
            Qualification string
            GitHub       string
            Mobile       string
        }{
            Name:         name,
            Email:        email,
            Dob:          dob,
            Qualification: qualification,
            GitHub:       github,
            Mobile:       mobile,
        }

        tmpl, err := template.ParseFiles("./static/candidate_dashboard.html")
        if err != nil {
            log.Println("Error parsing template:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        err = tmpl.Execute(w, data)
        if err != nil {
            log.Println("Error executing template:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    })


    
    http.HandleFunc("/recruiter_dashboard", func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("username")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            log.Println("Error retrieving cookie:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    
        username := cookie.Value
        var name, email, post, company, experience, branch, mobile string
    
        err = db.QueryRow("SELECT name, email, post, company, experience, branch, mobile FROM recruiter WHERE email=$1", username).Scan(&name, &email, &post, &company, &experience, &branch, &mobile)
        if err != nil {
            log.Println("Error querying recruiter details:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    
        additionalDetails := fmt.Sprintf("Post: %s<br>Company: %s<br>Experience: %s<br>Branch: %s<br>Mobile: %s", post, company, experience, branch, mobile)
    
        fmt.Fprintf(w, "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width, initial-scale=1.0'><title>Recruiter Dashboard</title><link rel='stylesheet' href='/static/css/style.css'></head><body><h1>Recruiter Dashboard</h1><p>Name: %s<br>Email: %s<br>%s</p><p><a href='/'>Logout</a></p></body></html>", name, email, additionalDetails)
    })
    
    http.HandleFunc("/admin_dashboard", func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("username")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            log.Println("Error retrieving cookie:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    
        username := cookie.Value
        var name, email, role, additionalDetails string
    
        err = db.QueryRow(`
            SELECT name, email, 'admin' AS role FROM admin WHERE email=$1
        `, username).Scan(&name, &email, &role)
        
        if err != nil {
            log.Println("Error querying database:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    
       
        
    
        fmt.Fprintf(w, "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width, initial-scale=1.0'><title>Admin Dashboard</title><link rel='stylesheet' href='/static/css/style.css'></head><body><h1>Admin Dashboard</h1><p>Name: %s<br>Email: %s<br>Role: %s<br>%s</p><p><a href='/'>Logout</a></p></body></html>", name, email, role, additionalDetails)
    })
    

    
    

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server started at http://localhost:%s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
