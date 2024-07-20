package middleware

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(db *mongo.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			sessionToken := cookie.Value
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Check if the user exists with the given sessionToken (User ID)
			userID, err := primitive.ObjectIDFromHex(sessionToken)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			filter := bson.M{"_id": userID}
			user := struct {
				ID primitive.ObjectID `bson:"_id"`
			}{}
			err = db.Collection("users").FindOne(ctx, filter).Decode(&user)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the session_token cookie is present
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			// If the cookie is not present, redirect to the login page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Proceed to the requested page if the session token is present
		next.ServeHTTP(w, r)
	})
}
