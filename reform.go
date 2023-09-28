package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Store struct {
	db *reform.DB
}

// Test project. Crud works.
// No field validation!
func main() {
	//Use in prod with:
	//db, err := GetDB(os.Getenv("PG_URI"))
	db, err := GetDB("postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("DB connection success")
	e := echo.New()
	logger := log.New(os.Stderr, "SQL: ", log.Flags())

	reformDb := reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))
	store := Store{db: reformDb}

	e.GET("/users", store.getUsers)
	e.POST("/users", store.saveUser)
	e.GET("/users/:id", store.getUser)
	e.PUT("/users", store.updateUser)
	e.DELETE("/users/:id", store.deleteUser)
	e.Logger.Fatal(e.Start(":1323"))
}

// e.PUT("/users/", store.updateUser)
func (s *Store) updateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "Bad body request")
	}
	if err := s.db.Update(u); err != nil {
		return c.String(http.StatusBadRequest, "No such user")
	}
	return c.JSON(http.StatusOK, "Updated")
}

// e.POST("/users", store.saveUser)
func (s *Store) saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := s.db.Save(u); err != nil {
		return c.String(http.StatusBadRequest, "User not saved")
	}
	return c.JSON(http.StatusOK, "Saved")
}

// e.DELETE("/users/:id", store.deleteUser)
func (s *Store) deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	res, err := s.db.DeleteFrom(UserTable, fmt.Sprintf("where id = %d", id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Удалено строк "+strconv.Itoa(int(res)))
}

// e.GET("/users/:id", store.getUser)
func (s *Store) getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	res, err := s.db.FindAllFrom(UserTable, "id", id)
	if res == nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, res)
}

// e.GET("/", store.getUsers)
func (s *Store) getUsers(c echo.Context) error {
	res, err := s.db.SelectAllFrom(UserTable, "")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func GetDB(uri string) (*sql.DB, error) {
	DB, err := PgxCreateDB(uri)
	if err != nil {
		return nil, err
	}
	DB.SetMaxIdleConns(2)
	DB.SetMaxOpenConns(4)
	DB.SetConnMaxLifetime(time.Duration(30) * time.Minute)
	return DB, nil
}

func PgxCreateDB(uri string) (*sql.DB, error) {
	connConfig, _ := pgx.ParseConfig(uri)
	afterConnect := stdlib.OptionAfterConnect(func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `
			 CREATE TABLE IF NOT EXISTS users(
			 	id SERIAL,
				name varchar NOT NULL,
				age int
			 );
		`)
		if err != nil {
			return err
		}
		return nil
	})
	pgxdb := stdlib.OpenDB(*connConfig, afterConnect)
	return pgxdb, nil
}
