package models

type Joint struct {
	Id   int32
	Jobs []int32 `sql:",array"`
	Name string
}
