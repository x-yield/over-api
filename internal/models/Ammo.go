package models

type Ammo struct {
	Id       int32
	Url      string
	Bucket   string
	Key      string
	LastUsed float64
	Type     string
	Author   string
}
