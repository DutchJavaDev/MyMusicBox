package db

import (
	"api/logging"
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const databaseFileName = "localDatabase.db"

func CreateDatabase() {
	db, err := sql.Open("sqlite", databaseFileName)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	initScriptPath := "db/init_script.sql"

	initQuery, err := os.ReadFile(initScriptPath)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(initQuery))

	if err != nil {
		panic(err)
	}
}

func createConnection() (*sql.DB, error) {
	return sql.Open("sqlite", databaseFileName)
}

func InsertMusic(music *Music) {
	con, _ := createConnection()
	defer con.Close()

	query := "insert into music (Name, Url, Path) values (?, ?, ?)"

	prep, _ := con.Prepare(query)
	defer prep.Close()

	result, err := prep.Exec(music.Name, music.Url, music.Path)

	if err != nil {
		logging.Log(err.Error())
	}

	row, _ := result.LastInsertId()

	logging.Log(row)

}

func GetMusic() *[]Music {
	con, _ := createConnection()
	defer con.Close()

	query := "select Id, Name, Url, Path from music"

	var m Music

	row, _ := con.Query(query)

	musics := make([]Music, 0)

	for row.Next() {
		scanError := row.Scan(&m.Id, &m.Name, &m.Url, &m.Path)

		if scanError != nil {
			logging.Log(scanError.Error())
			continue
		}

		musics = append(musics, m)
	}

	return &musics
}
