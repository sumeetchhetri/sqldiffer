package db

import (
	//sql "database/sql"
	"fmt"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//SqlsDb -
type SqlsDb struct {
}

//GenerateURL -
func (db *SqlsDb) GenerateURL(det *c.SchemaDiffAction) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		*det.User, *det.Password, *det.Host, det.Port, *det.DatabaseName)
}

//Version -
func (db *SqlsDb) Version() string {
	return ""
}

//Preface -
func (db *SqlsDb) Preface(dbe *pb2.Db) string {
	var b bytes.Buffer
	b.WriteString("--")
	b.WriteString("\n-- SQL Server database dump\n\n")
	return b.String()
}

//Create -
func (db *SqlsDb) Create(dbe *pb2.Db) string {
	return fmt.Sprintf("\nCREATE DATABASE %s;\n", *dbe.Name)
}

//Connect -
func (db *SqlsDb) Connect(dbe *pb2.Db) string {
	return fmt.Sprintf(`USE [%s] \nGO\n`, *dbe.Name)
}

//CreateSchema -
func (db *SqlsDb) CreateSchema(dbe *pb2.Db) string {
	return ""
}
