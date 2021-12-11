package db

import (
	//sql "database/sql"
	"fmt"
	//Import mysql driver
	_ "github.com/go-sql-driver/mysql"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//MysqlDb -
type MysqlDb struct {
}

//GenerateURL -
func (db *MysqlDb) GenerateURL(det *c.SchemaDiffAction) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		*det.User, *det.Password, *det.Host, det.Port, *det.DatabaseName)
}

//Version -
func (db *MysqlDb) Version() string {
	return ""
}

//Preface -
func (db *MysqlDb) Preface(dbe *pb2.Db) string {
	var b bytes.Buffer
	b.WriteString("--")
	b.WriteString("\n-- Mysql database dump\n\n")
	return b.String()
}

//Create -
func (db *MysqlDb) Create(dbe *pb2.Db) string {
	return fmt.Sprintf("\nCREATE DATABASE %s;\n", *dbe.Name)
}

//Connect -
func (db *MysqlDb) Connect(dbe *pb2.Db) string {
	return fmt.Sprintf("\n\\USE %s;\n", *dbe.Name)
}

//CreateSchema -
func (db *MysqlDb) CreateSchema(dbe *pb2.Db) string {
	return ""
}
