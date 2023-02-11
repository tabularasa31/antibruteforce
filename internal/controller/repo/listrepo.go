package repo

import (
	"github.com/boltdb/bolt"
)

// ListRepo -.
type ListRepo struct {
	db *bolt.DB
}

const addressBucket = "blackandwhite"

func NewListRepo(db *bolt.DB) *ListRepo {
	return &ListRepo{db: db}
}

// Save Сохранить адрес в lists -.
func (l *ListRepo) Save(subnet string, list string) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(addressBucket)).Put([]byte(subnet), []byte(list))
		return err
	})

}

// Delete Удалить адрес из lists -.
func (l *ListRepo) Delete(subnet string) error {
	return l.db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(addressBucket)).Delete([]byte(subnet))
		return err
	})
}
