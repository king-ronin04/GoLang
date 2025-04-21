package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var tpl = template.Must(template.New("users").Parse(`
<!DOCTYPE html>
<html>
<head><title>Users</title></head>
<body>
<h1>User List</h1>
<ul>
{{range .}}
	<li>{{.ID}}: {{.Name}} - {{.Email}}</li>
{{end}}
</ul>
</body>
</html>
`))

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				log.Println(err)
				continue
			}
			users = append(users, u)
		}

		tpl.Execute(w, users)
	})

	fmt.Println("Server running at http://localhost:8080/users")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
