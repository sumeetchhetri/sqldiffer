package view

import (
	"bytes"
	sql "database/sql"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//PgView -
type PgView struct {
}

//GenerateNew -
func (db *PgView) GenerateNew(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE VIEW ")
	b.WriteString(*vw.Name)
	b.WriteString(" AS ")
	b.WriteString(*vw.Definition)
	b.WriteString("\n")
	return b.String()
}

//GenerateUpd -
func (db *PgView) GenerateUpd(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(vw, context))
	b.WriteString(db.GenerateNew(vw, context))
	return b.String()
}

//GenerateDel -
func (db *PgView) GenerateDel(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP VIEW ")
	b.WriteString(*vw.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *PgView) Query(context interface{}) string {
	return "select table_name,view_definition from INFORMATION_SCHEMA.views WHERE table_schema = ANY (current_schemas(false))"
}

//FromResult -
func (db *PgView) FromResult(rows *sql.Rows, context interface{}) *pb2.View {
	return c.GetViewFromRow(rows, context)
}
