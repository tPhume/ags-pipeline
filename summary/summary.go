package summary

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Data map[string]float64

type Summary struct {
	Id           string
	UserId       string
	ControllerId string
	Date         string
	Data         Data
}

type ReadStorage interface {
	Read(ctx context.Context, summary map[string]*Summary) error
}

type WriteStorage interface {
	Write(ctx context.Context, summary map[string]*Summary) error
}

type Storage struct {
	Reader ReadStorage
	Writer WriteStorage
}

func (s *Storage) Handle(ctx *gin.Context) {
	panic("implement me")
}