package trigger

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"strings"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//OrclTrigger -
type OrclTrigger struct {
	SchemaName string
}

//GenerateNew -
func (db *OrclTrigger) GenerateNew(tr *pb2.Trigger, context interface{}) string {
	return *tr.Definition + ";\n/"
}

//GenerateUpd -
func (db *OrclTrigger) GenerateUpd(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(tr, context))
	b.WriteString(db.GenerateNew(tr, context))
	return b.String()
}

//GenerateDel -
func (db *OrclTrigger) GenerateDel(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP TRIGGER \"")
	b.WriteString(*tr.Name)
	b.WriteString("\";\n/")
	return b.String()
}

//Query -
func (db *OrclTrigger) Query(context interface{}) string {
	return `SELECT trigger_name, table_name, trigger_type, triggering_event, 
		trigger_body, description, dbms_metadata.get_ddl('TRIGGER', trigger_name) from user_triggers`
}

//FromResult -
func (db *OrclTrigger) FromResult(rows *sql.Rows, context interface{}) *pb2.Trigger {
	return c.GetTriggerFromRow(rows, context)
}

//DefineQuery -
func (db *OrclTrigger) DefineQuery(context interface{}) string {
	return ""
}

//GetDefinition -
func (db *OrclTrigger) GetDefinition(rows *sql.Rows) string {
	return ""
}
