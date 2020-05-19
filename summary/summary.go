package summary

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Data map[string]float64

type Summary struct {
	Id           string
	UserId       string
	ControllerId string
	Data         Data
}

type ReadStorage interface {
	ReadMean(ctx context.Context, summary map[string]*Summary) error

	ReadMedian(ctx context.Context, summary map[string]*Summary) error
}

type WriteStorage interface {
	WriteMean(ctx context.Context, summary map[string]*Summary) error

	WriteMedian(ctx context.Context, summary map[string]*Summary) error
}

type Storage struct {
	Reader ReadStorage
	Writer WriteStorage
}

func (s *Storage) HandleMean(ctx *gin.Context) {
	summary := make(map[string]*Summary)

	// Read data from storage
	if err := s.Reader.ReadMean(ctx, summary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"at": "reader", "err": err, "details": err.Error()})
		return
	}

	// Write data to storage
	if err := s.Writer.WriteMean(ctx, summary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"at": "writer", "err": err, "details": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (s *Storage) HandleMedian(ctx *gin.Context) {
	summary := make(map[string]*Summary)

	// Read data from storage
	if err := s.Reader.ReadMedian(ctx, summary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"at": "reader", "err": err, "details": err.Error()})
		return
	}

	// Write data to storage
	if err := s.Writer.WriteMedian(ctx, summary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"at": "writer", "err": err, "details": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}
