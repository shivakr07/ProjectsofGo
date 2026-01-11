// DB Setup ----------
package storage

// we will use interfaces here
// we can make it like pluging as we did for payment methods
// so we can switch to any DB with minimal changes

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
}
