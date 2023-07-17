package sqlstorage

//import (
//	"context"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
//	"testing"
//	"time"
//)
//
//func TestStorageCreate(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer db.Close()
//
//	event := models.Event{
//		Title:                "test title",
//		Date:                 time.Now(),
//		Duration:             time.Second,
//		Description:          "test description",
//		UserID:               4,
//		NotificationInterval: time.Second,
//	}
//
//	ctx := context.Background()
//
//	mock.ExpectQuery()
//}
