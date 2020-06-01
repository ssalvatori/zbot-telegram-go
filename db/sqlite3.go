package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	log "github.com/sirupsen/logrus"
)

type ZbotSqlite3Database struct {
	Db   *sql.DB
	File string
}

func (d *ZbotSqlite3Database) GetConnectionString() string {
	return fmt.Sprintf("DB: %s", d.File)
}

func (d *ZbotSqlite3Database) UserIgnoreList() ([]UserIgnore, error) {
	log.Debug("Getting ignore list")
	statement := "SELECT username, since, until FROM ignore_list"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query()
	if err2 != nil {
		panic(err2)
	}
	defer rows.Close()

	var users []UserIgnore
	var user UserIgnore
	for rows.Next() {
		err2 := rows.Scan(&user.Username, &user.Since, &user.Until)
		if err2 != nil {
			return nil, err2
		}
		users = append(users, user)
	}
	return users, nil
}

//Init start connection with sql database
func (d *ZbotSqlite3Database) Init() error {

	if _, err := os.Stat(d.File); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Sqlite file [%s] does not exist!", d.File))
	}

	log.Info("Connecting to database: " + d.File)
	db, err := sql.Open("sqlite3", "file:"+d.File+"?cache=shared&mode=rwc")
	if err != nil {
		log.Error(err)
		return err
	}
	if db == nil {
		log.Error(err)
		return errors.New("Error connecting")
	}
	d.Db = db

	return nil
}

func (d *ZbotSqlite3Database) Close() {
	log.Debug("Closing connection")
	d.Db.Close()
}

func (d *ZbotSqlite3Database) Statistics() (string, error) {
	statement := "select count(*) as total from definitions"
	var totalCount string
	err := d.Db.QueryRow(statement).Scan(&totalCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return totalCount, errors.New("No Rows found")
		} else {
			return totalCount, err
		}
	}

	return totalCount, err
}

func (d *ZbotSqlite3Database) Top() ([]DefinitionItem, error) {

	statement := "SELECT term FROM definitions ORDER BY hits DESC LIMIT 10"
	rows, err := d.Db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DefinitionItem
	for rows.Next() {
		var key string
		err2 := rows.Scan(&key)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, DefinitionItem{Term: key})
	}

	return items, nil
}
func (d *ZbotSqlite3Database) Rand() (DefinitionItem, error) {
	var def DefinitionItem

	statement := "SELECT term, meaning FROM definitions ORDER BY random() LIMIT 1"
	rows, err := d.Db.Query(statement)
	if err != nil {
		return def, err
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&def.Term, &def.Meaning)
		if err2 != nil {
			return def, err2
		}
	}

	return def, nil

}

func (d *ZbotSqlite3Database) Last() (DefinitionItem, error) {
	var def DefinitionItem
	statement := "SELECT term, meaning FROM definitions ORDER BY id DESC LIMIT 1"

	err := d.Db.QueryRow(statement).Scan(&def.Term, &def.Meaning)
	if err != nil {
		if err == sql.ErrNoRows {
			return def, errors.New("Nothing defined")
		} else {
			log.Fatal(err)
			return def, err
		}
	}

	return def, nil
}
func (d *ZbotSqlite3Database) Get(term string) (DefinitionItem, error) {
	var def DefinitionItem
	statement := "SELECT id, term, meaning, author, date FROM definitions WHERE term = ? COLLATE NOCASE LIMIT 1"
	err := d.Db.QueryRow(statement, term).Scan(&def.Id, &def.Term, &def.Meaning, &def.Author, &def.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return DefinitionItem{Term: "", Meaning: ""}, nil
		} else {
			log.Fatal(err)
			return def, err
		}
	}

	statement = "UPDATE definitions SET hits = hits + 1 WHERE id = ?"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return def, err
	}

	_, err = stmt.Exec(def.Id)
	if err != nil {
		return def, err
	}

	return def, nil
}

func (d *ZbotSqlite3Database) _set(term string, def DefinitionItem) (sql.Result, error) {
	statement := "INSERT INTO definitions (term, meaning, author, locked, active, date, hits, link) VALUES (?,?,?,?,?,?,?,?)"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	return stmt.Exec(term, def.Meaning, def.Author, 1, 1, def.Date, 0, 0)

}

func (d *ZbotSqlite3Database) Set(def DefinitionItem) (string, error) {
	count := 1
	term := def.Term
	log.Debug(def)
	for {
		_, err := d._set(term, def)
		if err != nil {
			log.Debug("SQL insert error: ", err.Error())
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				term = fmt.Sprintf("%s%d", def.Term, count)
				log.Debug(fmt.Sprintf("New Term: %s", term))
				count = count + 1
			} else {
				return "", err
			}
		} else {
			log.Debug("trying with: ", term)
			break
		}
	}
	return term, nil

}

func (d *ZbotSqlite3Database) Append(def DefinitionItem) error {
	statement := "UPDATE definitions SET meaning = meaning || ?, date = ?, author = ? WHERE term = ?"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(fmt.Sprintf(" %s", def.Meaning), def.Date, def.Author, def.Term)
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}

func (d *ZbotSqlite3Database) Find(criteria string) ([]DefinitionItem, error) {
	var items []DefinitionItem
	statement := "SELECT term FROM definitions WHERE meaning like ? ORDER BY random() COLLATE NOCASE LIMIT 20"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return items, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil {
		return items, err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return items, err2
		}
		items = append(items, DefinitionItem{Term: result})
	}
	return items, nil
}
func (d *ZbotSqlite3Database) Search(criteria string) ([]DefinitionItem, error) {
	statement := "SELECT term FROM definitions WHERE term like ? ORDER BY random() COLLATE NOCASE LIMIT 10"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil {
		panic(err2)
	}
	defer rows.Close()

	var items []DefinitionItem
	var result string
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, DefinitionItem{Term: result})
	}
	return items, nil
}

func (d *ZbotSqlite3Database) UserLevel(username string) (string, error) {
	var level string
	statement := "SELECT level FROM users WHERE username = ? COLLATE NOCASE LIMIT 1"
	err := d.Db.QueryRow(statement, username).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			return "0", nil
		} else {
			return level, err
		}
	}

	return level, nil
}
func (d *ZbotSqlite3Database) UserIgnoreInsert(username string) error {
	statement := "INSERT INTO ignore_list (username, since, until) VALUES (?,?,?)"
	stmt, err := d.Db.Prepare(statement)

	if err != nil {
		return err
	}

	since := time.Now().Unix()
	tenMinutes := 10 * time.Minute
	until := time.Now().Add(tenMinutes).Unix()

	_, err = stmt.Exec(username, since, until)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//UserCheckIgnore check if user there is any row in ignore_list table for some username
// and an until greater than the current time
func (d *ZbotSqlite3Database) UserCheckIgnore(username string) bool {
	ignored := false

	now := time.Now().Unix()

	var level string
	statement := "SELECT count(*) as total FROM ignore_list WHERE username = ? AND until >= ?"
	err := d.Db.QueryRow(statement, username, now).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			ignored = false
		} else {
			log.Error(err)
			return ignored
		}
	}
	levelInt, _ := strconv.Atoi(level)
	log.Debug("Ingored ", levelInt)
	if levelInt > 0 {
		ignored = true
	}

	return ignored
}
func (d *ZbotSqlite3Database) UserCleanupIgnorelist() error {
	for {
		log.Debug("Cleaning ignore list")
		now := time.Now().Unix()
		statement := "DELETE FROM ignore_list WHERE until <= ?"
		stmt, err := d.Db.Prepare(statement)
		_, err = stmt.Query(now)
		if err != nil {
			if err == sql.ErrNoRows {

			} else {
				log.Fatal(err)
				return err
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func (d *ZbotSqlite3Database) Lock(item DefinitionItem) error {

	statement := "UPDATE definitions SET locked = 1, locked_by = ? WHERE term = ?"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.Author, item.Term)
	if err != nil {
		return err
	}

	return nil

}

func (d *ZbotSqlite3Database) Forget(item DefinitionItem) error {
	statement := "DELETE definitions WHERE term = ? LIMIT 1"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.Term)
	if err != nil {
		return err
	}

	return nil
}