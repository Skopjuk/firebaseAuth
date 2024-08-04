package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/http"
	"tryFirebase/midleware"
)

func main() {
	opt := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	premiumHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Context().Value("verifiedToken")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tmpl, err := template.ParseFiles("static/premium.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error parsing templates: %v\n", err)
			return
		}

		data := map[string]interface{}{
			"Token": token,
		}

		tmpl.Execute(w, data)
	})
	
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.Handle("/premium", midleware.AuthMiddleWare(client)(premiumHandler))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error creating server: %v\n", err)
	}
}
