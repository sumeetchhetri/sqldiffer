package trigger

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	"strings"
	//"fmt"
	proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//PgTrigger -
type PgTrigger struct {
	SchemaName string
}

//GenerateNew -
func (db *PgTrigger) GenerateNew(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\n")
	b.WriteString(*tr.FunctionDef)
	b.WriteString(";\nCREATE TRIGGER ")
	b.WriteString(*tr.Name)
	b.WriteString(" ")
	b.WriteString(*tr.When)
	b.WriteString(" ")
	b.WriteString(*tr.Action)
	b.WriteString(" ON ")
	b.WriteString(*tr.TableName)
	b.WriteString(" FOR EACH ROW ")
	b.WriteString(*tr.Definition)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgTrigger) GenerateUpd(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(tr, context))
	b.WriteString(db.GenerateNew(tr, context))
	return b.String()
}

//GenerateDel -
func (db *PgTrigger) GenerateDel(tr *pb2.Trigger, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP TRIGGER IF EXISTS ")
	b.WriteString(*tr.Name)
	b.WriteString(" ON ")
	b.WriteString(*tr.TableName)
	b.WriteString(";\nDROP FUNCTION IF EXISTS ")
	b.WriteString(*tr.Function)
	b.WriteString("();\n")
	return b.String()
}

//Query -
func (db *PgTrigger) Query(context interface{}) string {
	return `select trigger_name, event_object_table, action_timing, event_manipulation, 
		action_statement, '', '' from information_schema.triggers where trigger_schema = ANY (current_schemas(false))`
}

//FromResult -
func (db *PgTrigger) FromResult(rows *sql.Rows, context interface{}) *pb2.Trigger {
	tr := c.GetTriggerFromRow(rows, context)
	tr.Definition = proto.String(*tr.Function)
	if strings.Index(strings.ToLower(*tr.Definition), "execute procedure ") == 0 {
		*tr.Function = strings.Replace(strings.ToLower(*tr.Definition), "execute procedure ", "", 1)
		*tr.Function = strings.Replace(*tr.Function, "()", "", 1)
	}
	return tr
}

//DefineQuery -
func (db *PgTrigger) DefineQuery(context interface{}) string {
	function := context.(string)
	return fmt.Sprintf(`select pg_get_functiondef(pg_proc.oid) from pg_proc LEFT JOIN pg_namespace n 
		ON n.oid = pg_proc.pronamespace where nspname = ANY (current_schemas(false)) and proname = '%s' 
		and pg_catalog.pg_get_function_result(pg_proc.oid) = 'trigger'`, function)
}

//GetDefinition -
func (db *PgTrigger) GetDefinition(rows *sql.Rows) string {
	var tmp string

	err := rows.Scan(&tmp)
	if err != nil {
		c.Fatal("Error fetching details from database", err)
	}
	if strings.Index(strings.ToUpper(tmp), "CREATE OR REPLACE FUNCTION "+strings.ToUpper(db.SchemaName)+".") == 0 {
		eprf := "CREATE OR REPLACE FUNCTION " + strings.ToUpper(db.SchemaName) + "."
		tmp = "CREATE OR REPLACE FUNCTION " + tmp[len(eprf):]
	} else if strings.Index(strings.ToUpper(tmp), "CREATE FUNCTION "+strings.ToUpper(db.SchemaName)+".") == 0 {
		eprf := "CREATE FUNCTION " + strings.ToUpper(db.SchemaName) + "."
		tmp = "CREATE FUNCTION " + tmp[len(eprf):]
	}
	return strings.Replace(tmp, "\r\n", "\n", -1)
}
