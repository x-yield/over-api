package models

type Collection struct {
	Id      int32
	Env     string
	Project string
	Service string
	Ref     string
	Type    string
	Name    string
	Author  string
}
