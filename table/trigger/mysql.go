package trigger

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "sqldiffer/protos"
)

//MysqlTrigger -
type MysqlTrigger struct {
	SchemaName string
}

//GenerateNew -
func (db *MysqlTrigger) GenerateNew(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDELIMITER $$\nCREATE TRIGGER ")
	b.WriteString(*tr.Name)
	b.WriteString(" ")
	b.WriteString(*tr.When)
	b.WriteString(" ")
	b.WriteString(*tr.Action)
	b.WriteString(" ON ")
	b.WriteString(*tr.TableName)
	b.WriteString(" FOR EACH ROW ")
	b.WriteString(*tr.Definition)
	b.WriteString("\nEND $$\nDELIMETER ;\n")
	return b.String()
}

//GenerateUpd -
func (db *MysqlTrigger) GenerateUpd(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(tr, context))
	b.WriteString(db.GenerateNew(tr, context))
	return b.String()
}

//GenerateDel -
func (db *MysqlTrigger) GenerateDel(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP TRIGGER IF EXISTS ")
	b.WriteString(*tr.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *MysqlTrigger) Query(context interface{}) string {
	return fmt.Sprintf(`select trigger_name, event_object_table, action_timing, event_manipulation, '', '',
		action_statement from information_schema.triggers
		where trigger_schema = '%s'`, context.(string))
}

//FromResult -
func (db *MysqlTrigger) FromResult(rows *sql.Rows, context interface{}) *pb2.Trigger {
	return c.GetTriggerFromRow(rows, context)
}

//DefineQuery -
func (db *MysqlTrigger) DefineQuery(context interface{}) string {
	return ""
}

//GetDefinition -
func (db *MysqlTrigger) GetDefinition(rows *sql.Rows) string {
	return ""
}
