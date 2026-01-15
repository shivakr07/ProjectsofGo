// DB Setup ----------
package storage

import "github.com/shivakr07/students-api/internal/types"

// we will use interfaces here
// we can make it like pluging as we did for payment methods
// so we can switch to any DB with minimal changes

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
