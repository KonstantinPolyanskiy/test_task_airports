package common

type IndexEntry struct {
	Key          string
	NumericValue float64

	// IsNumeric Для кейса, когда колонка представляет из себя число
	IsNumeric bool

	Offset int64

	// Length Длинна строки в байтах
	Length int
}
