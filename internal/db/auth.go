package db

func (d *db) GetHashedPassword(login string) (string, error) {
	// language=PostgreSQL
	q := `
SELECT password FROM auth WHERE login = $1
`
	return q, nil
}
