package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Agent struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	UUID          string        `bson:"uuid" mgo:"index:1"`
	RelateId      int32         `bson:"relate_id" mgo:"index:-1"`
	HostName      string        `bson:"host_name"`
	HostIP        []string      `bson:"host_ip"`
	OsType        string        `bson:"os_type"`
	AssertGroupId int32         `bson:"asset_group_id"`
	AssertLabel   []string      `bson:"asset_label"`
	Status        string        `bson:"status"`
	StatusLong    int32         `bson:"status_long"`
	RiskCount     float32       `bson:"risk_count"`
	LastCheckTm   time.Time     `bson:"last_check_tm"`
	CheckStatus   string        `bson:"check_status"`
	LastReportTm  time.Time     `bson:"last_report_tm"`
	LastOnlineTm  time.Time     `bson:"last_online_tm"`
	CreateTime    time.Time     `bson:"create_tm"`
}

func (a *Agent) CollectName() string {
	return "tb_agent"
}

type User struct {
	ID   int    `bson:"id"`
	Name string `bson:"name" mgo:"index:1"`
	Age  int    `bson:"age"`
}

func (u *User) CollectName() string {
	return "tb_user"
}
