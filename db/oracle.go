package db

import (
	"fmt"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	c "sqldiffer/common"
	pb2 "sqldiffer/protos"
)

//OrclDb -
type OrclDb struct {
}

//GenerateURL -
func (db *OrclDb) GenerateURL(det *c.SchemaDiffAction) string {
	return fmt.Sprintf("%s/%s@%s:%d/%s",
		*det.User, *det.Password, *det.Host, det.Port, *det.DatabaseName)
}

//Version -
func (db *OrclDb) Version() string {
	return ""
}

//Preface -
func (db *OrclDb) Preface(dbe *pb2.Db) string {
	var b bytes.Buffer
	b.WriteString("--")
	b.WriteString("\n-- Oracle database dump\n\n")
	return b.String()
}

//Create -
func (db *OrclDb) Create(dbe *pb2.Db) string {
	return fmt.Sprintf("\nCREATE DATABASE%s;\n", *dbe.Name)
}

//Connect -
func (db *OrclDb) Connect(dbe *pb2.Db) string {
	return ""
}

//CreateSchema -
func (db *OrclDb) CreateSchema(dbe *pb2.Db) string {
	return ""
}
