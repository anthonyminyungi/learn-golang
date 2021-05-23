package mydict

import "errors"

type Dictionary map[string]string

var (
	errNotFound     = errors.New("not found")
	errAlreadyExist = errors.New("the word you're going to add is already exist")
	errCantUpdate   = errors.New("you can't update the word because it's not exist")
	errCantDelete   = errors.New("you can't delete the word because it's not exist")
)

func (d Dictionary) Search(word string) (string, error) {
	value, exist := d[word]
	if exist {
		return value, nil
	}

	return "", errNotFound
}

func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errAlreadyExist
	}

	return nil
}

func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		return errCantUpdate
	case nil:
		d[word] = def
	}

	return nil
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		return errCantDelete
	case nil:
		delete(d, word)
	}
	return nil
}
