package postgre

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

//go:embed migration
var migration embed.FS

type (
	postgre struct {
		db *sql.DB
	}
	IPostgre interface {
		Migrate() error
		Begin() (ITransaction, error)
		Close() error
	}
)

func NewPostgre(dsn string) IPostgre {
	var (
		err error
		pg  postgre
	)
	for i := 0; i < 10; i++ {
		pg.db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Panic(err)
		}
		if err = pg.db.Ping(); err != nil {
			log.Println("\nReconnecting to your database host....")
		} else {
			log.Println("\nconnected to database")
			break
		}
		time.Sleep(4 * time.Second)
	}

	if err != nil {
		log.Panic("cannot connect to database")
	}
	if err = executeFiles(pg.db, "migration/function"); err != nil {
		log.Panicf("cannot create db function error : %v", err)
	}
	return &pg
}

func executeFiles(db *sql.DB, location string) (err error) {
	files, err := migration.ReadDir(location)
	if err != nil {
		return
	}
	for _, v := range files {
		if err = executeFile(db, location+"/"+v.Name()); err != nil {
			return
		}
	}
	return
}

func executeFile(db *sql.DB, fileName string) (err error) {
	if sql, err := migration.ReadFile(fileName); err == nil {
		if _, err := db.Exec(string(sql)); err != nil {
			return fmt.Errorf("error executing file %s : %s", fileName, err.Error())
		}
	}
	return
}

func (r *postgre) Migrate() (err error) {
	return executeFiles(r.db, "migration/table")
}

func (r *postgre) Seed(fileName string) (err error) {
	return executeFile(r.db, "migration/seed/"+fileName)
}

func (r *postgre) Close() (err error) {
	r.db.Close()
	return
}

func (r *postgre) Begin() (out ITransaction, err error) {
	tx, err := r.db.Begin()
	out = &transaction{tx: tx}
	return
}
