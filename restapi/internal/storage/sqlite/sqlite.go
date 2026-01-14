package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	//since we are not using directly like obj.something so we are using indirectly so we used _
	"github.com/shivakr07/students-api/internal/config"
	"github.com/shivakr07/students-api/internal/types"
)

//here we implement the interfaces created in the storage.go

type Sqlite struct {
	Db *sql.DB
}

// since we don't have constructor concept but we replicate similar using New [as convention]
func New(cfg *config.Config) (*Sqlite, error) {
	//db connection
	//we need to pass the driver inside the open method and storage path
	//open method returns two thing instance of the db and error
	// need to install this driver : browse go sqlite driver [mattnn git]
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		//we return sqlite and error
		//since till here we are getting error so instead of sqlite we are returning the nil
		return nil, err
	}

	//create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	//if everthing okay then return sqlite
	return &Sqlite{
		Db: db,
	}, nil

}

// exec return two things res [result of the query] and error
// since we are not using that so we kept _ but then we need to remove the : from := but if you res, err := then we need to use the :

// 	//why nil in 	}, nil while returning db? bcause we need to return error since we have no error now so we pass the nil

// implementing func to implement interface
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	//to create the records in the db
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	//we need to close this statement also after function execution
	defer stmt.Close()

	// we put ? ? ? [placeholders] to avoid the SQL injection as we don't pass the data direct which we are receiving
	//these values we are reveiving the func
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	//in result we have query result
	// we get methods from Exec
	// LastInsertId() (int64, error) and RowsAffected() (int64, error)
	// [check by clicking ctrl + click to see the def]
	//why we are returning 0 [because in return type it should be int64]so 0 is zeroed value / empty value for int type
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
	//since here we don't have error

	//first we prepare the statement/query and then we bind the data
}

//how pluging helps
//if in future you want to use postgres
//except the query part like sql or postgre-sql
// you just need to import postgres and create struct of that
// and just implement that interface type Storage interface {CreateStudent ()}
// means you just need to implement that CreateStudent method over postgres struct
// and in the main instead of sqlite.New you will do postgres.New
// that's it ...
// or you can pass fake db also for testing
//POWER OF DEPENDENCY INJECTION

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
		//empty struct
	}
	defer stmt.Close()

	//whatever data we are getting from the db that needs to be deserialized so
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		//sometimes we get error like user not found
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("n o student found with id %s", fmt.Sprint(id))
		}
		//else this will be error mostly
		return types.Student{}, fmt.Errorf("qeury error : %w", err)
	}

	//if everything is okay
	return student, nil

	//this method adds the db data into the struct var student
	//ordering is very important
}
