package view

import (
	sql "database/sql"
	//"fmt"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
)

//SqlsView -
type SqlsView struct {
}

//GenerateNew -
func (db *SqlsView) GenerateNew(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString(*vw.Definition)
	b.WriteString("\n")
	return b.String()
}

//GenerateUpd -
func (db *SqlsView) GenerateUpd(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(vw, context))
	b.WriteString(db.GenerateNew(vw, context))
	return b.String()
}

//GenerateDel -
func (db *SqlsView) GenerateDel(vw *pb2.View, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString("DROP VIEW [")
	b.WriteString(*vw.Name)
	b.WriteString("];\n")
	return b.String()
}

//Query -
func (db *SqlsView) Query(context interface{}) string {
	return `select object_name(o.object_id),definition
		from sys.objects     o
		join sys.sql_modules m on m.object_id = o.object_id
		and o.type      = 'V'`
}

//FromResult -
func (db *SqlsView) FromResult(rows *sql.Rows, context interface{}) *pb2.View {
	return c.GetViewFromRow(rows, context)
}
