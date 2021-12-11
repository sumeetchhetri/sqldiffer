package view

import (
	"bytes"
	sql "database/sql"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//OrclView -
type OrclView struct {
}

//GenerateNew -
func (db *OrclView) GenerateNew(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nCREATE VIEW \"")
	b.WriteString(*vw.Name)
	b.WriteString("\" AS ")
	b.WriteString(*vw.Definition)
	b.WriteString(";\n/")
	return b.String()
}

//GenerateUpd -
func (db *OrclView) GenerateUpd(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(vw, context))
	b.WriteString(db.GenerateNew(vw, context))
	return b.String()
}

//GenerateDel -
func (db *OrclView) GenerateDel(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP VIEW ")
	b.WriteString(*vw.Name)
	b.WriteString(";\n/")
	return b.String()
}

//Query -
func (db *OrclView) Query(context interface{}) string {
	return "select view_name,text from user_views"
}

//FromResult -
func (db *OrclView) FromResult(rows *sql.Rows, context interface{}) *pb2.View {
	return c.GetViewFromRow(rows, context)
}
