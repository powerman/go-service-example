package dal

import (
	"time"
)

const (
	sqlContactAdd = `
 INSERT INTO Contact (name)
 VALUES (:name)
	`
	sqlContactLst = `
 SELECT id, name, ctime
   FROM Contact
  WHERE id > :since_id
  ORDER BY id ASC
  LIMIT :limit
	`
)

type (
	argContactAdd struct {
		Name string
	}

	argContactLst struct {
		SinceID int
		Limit   int
	}
	rowContactLst struct {
		ID    int
		Name  string
		Ctime time.Time
	}
)
