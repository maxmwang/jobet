// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"
)

type Company struct {
	Name      string
	Alias     string
	Site      string
	Priority  int64
	CreatedAt time.Time
}
