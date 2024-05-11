package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8088

type application struct {
	DSN    string
	Domain string
	DB repository.DatabaseRepo
	auth Auth
	JWTSecret string
	JWTIssuer string
	JWTAudience string
	CookieDomain string
}

/**
type Movie struct {
	ID int `json:"id"`
	Title string `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	RunTime int `json:"runtime"`
	MPAARating string `json:"mpaa_rating"`
	Description string `json:"description"`
	Image string `json:"image"`
	CreatedAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}
*/
func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgrase connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath: "/",
		CookieName: "__Host-refresh_token",
		CookieDomain: app.CookieDomain,
	}

	log.Println("Starting application on port", port)

	http.HandleFunc("/", app.Home)
	address := fmt.Sprintf(":%d", port)

	 err = http.ListenAndServe(address, app.routes())
	// err := http.ListenAndServe(address, nil)

	if err != nil {
	  	log.Fatalf("failed to start server: %s", err)
	  }
}

func (app *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}