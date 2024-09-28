package db

type DB interface {
	GetHashedPassword(login string) (string, error)
}
