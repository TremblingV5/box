package main

import (
	"context"
	"errors"
	"github.com/TremblingV5/box/dbtx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

// TestEntity represents a test database entity
type TestEntity struct {
	Id   int64
	Name string
}

// serviceMethod is a method for creating a new entity in the database.
// but, it will return an error
// so, the creation will be rollback
func serviceMethod(ctx context.Context) (err error) {
	ctx, persist := dbtx.WithTXPersist(ctx)
	defer func() {
		persist(err)
	}()

	entity := &TestEntity{
		Id:   time.Time{}.Unix(),
		Name: "test",
	}
	err = DB.Create(entity).Error
	if err != nil {
		return err
	}

	return errors.New("some error")
}

// init initializes the database connection
func init() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/test_entity?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	DB = db
}

// main executes the service method to demonstrate transaction rollback
func main() {
	err := serviceMethod(context.Background())
	log.Println(err)
}
