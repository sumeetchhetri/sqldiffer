package trigger

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//SqlsTrigger -
type SqlsTrigger struct {
	SchemaName string
}

//GenerateNew -
func (db *SqlsTrigger) GenerateNew(tr *pb2.Trigger, context interface{}) string {
	return c.GetSQLServerPreQuery() + *tr.Definition + ";\n"
}

//GenerateUpd -
func (db *SqlsTrigger) GenerateUpd(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(tr, context))
	b.WriteString(db.GenerateNew(tr, context))
	return b.String()
}

//GenerateDel -
func (db *SqlsTrigger) GenerateDel(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString("DROP TRIGGER [")
	b.WriteString(*tr.Name)
	b.WriteString("];\n")
	return b.String()
}

//Query -
func (db *SqlsTrigger) Query(context interface{}) string {
	return `SELECT sysobjects.name AS trigger_name, 
				OBJECT_NAME(parent_obj) AS table_name,'','','','',
				OBJECT_DEFINITION(id) AS trigger_definition
			FROM sysobjects 
			WHERE sysobjects.type = 'TR' `
}

//FromResult -
func (db *SqlsTrigger) FromResult(rows *sql.Rows, context interface{}) *pb2.Trigger {
	return c.GetTriggerFromRow(rows, context)
}

//DefineQuery -
func (db *SqlsTrigger) DefineQuery(context interface{}) string {
	return ""
}

//GetDefinition -
func (db *SqlsTrigger) GetDefinition(rows *sql.Rows) string {
	return ""
}
