package models

type Aggregate struct {
	Id           int32
	JobId        int32 `sql:",notnull"`
	Label        string
	Q50          float32 `sql:",notnull"`
	Q75          float32 `sql:",notnull"`
	Q80          float32 `sql:",notnull"`
	Q85          float32 `sql:",notnull"`
	Q90          float32 `sql:",notnull"`
	Q95          float32 `sql:",notnull"`
	Q98          float32 `sql:",notnull"`
	Q99          float32 `sql:",notnull"`
	Q100         float32 `sql:",notnull"`
	Avg          float32 `sql:",notnull"`
	OkCount      int64   `sql:",notnull"`
	ErrCount     int64   `sql:",notnull"`
	ResponseCode string
	NetRecv      float32 `sql:",notnull"`
	NetSend      float32 `sql:",notnull"`
}
