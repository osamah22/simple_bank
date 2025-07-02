package api

import (
	"github.com/gin-gonic/gin"
	"github.com/osamah22/simple_bank/db/db"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewStore(store *db.Store) *Server {
	return &Server{
		store: store,
	}
}
