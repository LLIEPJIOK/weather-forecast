package repository

import "fmt"

type ErrNotFound struct {
	id int
}

func NewErrNotFound(id int) error {
	return ErrNotFound{
		id: id,
	}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("no record with id=%d", e.id)
}
