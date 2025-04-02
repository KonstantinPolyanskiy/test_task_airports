package common

import "os"

type Index struct {
	Entries []IndexEntry

	File *os.File

	Column    int
	IsNumeric bool
}
