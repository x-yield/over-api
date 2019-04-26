package models

type TankSession struct {
	tableName     struct{} `sql:"tank_sessions,alias:ts"`
	Id            int32
	Tank          string
	Conf          string
	Name          string
	Failures      []string `sql:",array"`
	Stage         string
	Status        string
	ExternalId    string
	OverloadId    int32
	ExternalJoint string
	OverloadJoint int32
	Author        string
}
