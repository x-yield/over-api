package models

type Job struct {
	Id                 int32
	TestStart          float64
	TestStop           float64
	Config             string
	Author             string
	RegressionId       string
	CollectionIds      []int32 `sql:",array"`
	Description        string
	Tank               string
	Target             string
	EnvironmentDetails string
	Status             string
	AutostopTime       float64
	AutostopMessage    string
	Imbalance          int32
}
