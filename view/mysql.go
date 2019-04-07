package view

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "sqldiffer/protos"
)

//MysqlView -
type MysqlView struct {
}

//GenerateNew -
func (db *MysqlView) GenerateNew(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE VIEW ")
	b.WriteString(*vw.Name)
	b.WriteString(" AS ")
	b.WriteString(*vw.Definition)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlView) GenerateUpd(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(vw, context))
	b.WriteString(db.GenerateNew(vw, context))
	return b.String()
}

//GenerateDel -
func (db *MysqlView) GenerateDel(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP VIEW ")
	b.WriteString(*vw.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *MysqlView) Query(context interface{}) string {
	return fmt.Sprintf(`select table_name,view_definition from information_schema.views 
	where table_schema = '%s'`, context.(string))
}

//FromResult -
func (db *MysqlView) FromResult(rows *sql.Rows, context interface{}) *pb2.View {
	return c.GetViewFromRow(rows, context)
}
