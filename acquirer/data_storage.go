package acquirer

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type IDataStorage interface {
	Get(string) (*DBMetaInfo, bool)
	Put(*DBMetaInfo) bool
}

var DataStorage = newSqliteStorage()

type SqliteStorage struct {
	db *sql.DB
}

func newSqliteStorage() *SqliteStorage {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	storage := &SqliteStorage{db: db}
	storage.init()

	return storage
}

func (s *SqliteStorage) Close() {
	if err := s.db.Close(); err != nil {
		fmt.Println(err.Error())
	}
}

func (s *SqliteStorage) init() {
	const SQL = "CREATE TABLE IF NOT EXISTS metainfo (id INTEGER PRIMARY KEY AUTOINCREMENT,infohash TEXT, metainfo TEXT)"
	_, err := s.db.Exec(SQL)
	if err != nil {
		panic(err.Error())
	}
}

func (s *SqliteStorage) Get(infoHash string) *DBMetaInfo {
	const SQL = "SELECT id, infohash, metainfo FROM metainfo WHERE infohash = ?"
	info := DBMetaInfo{}

	if err := s.db.QueryRow(SQL, infoHash).Scan(&info.Id, &info.InfoHash, &info.Metadata); err != nil {
		fmt.Println(err)
		return nil
	}

	return &info
}

func (s *SqliteStorage) Put(info *DBMetaInfo) {
	const SQl = "INSERT INTO metainfo(infohash, metainfo) values(?,?)"

	stmt, err := s.db.Prepare(SQl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.InfoHash, info.Metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
}
