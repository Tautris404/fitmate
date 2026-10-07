package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}
