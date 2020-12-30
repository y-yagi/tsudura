package db

import (
	"github.com/dgraph-io/badger"
)

type Client struct {
	db *badger.DB
}

func Init(path string) (*Client, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Get(key []byte) ([]byte, error) {
	var val []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		if val, err = item.ValueCopy(nil); err != nil {
			return err
		}

		return nil
	})

	return val, err
}

func (c *Client) Set(key []byte, value []byte) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	return err
}

func (c *Client) Delete(key []byte) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})

	return err
}

func (c *Client) Term() {
	c.db.Close()
}
