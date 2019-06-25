package constraint

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//PgConstraint -
type PgConstraint struct {
}

//GenerateNew -
func (db *PgConstraint) GenerateNew(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE ")
	b.WriteString(*cn.TableName)
	b.WriteString(" ADD CONSTRAINT ")
	b.WriteString(*cn.Name)
	b.WriteString(" ")
	b.WriteString(*cn.Definition)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgConstraint) GenerateUpd(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(cn, context))
	b.WriteString(db.GenerateNew(cn, context))
	return b.String()
}

//GenerateDel -
func (db *PgConstraint) GenerateDel(cn *pb2.Constraint, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nALTER TABLE ")
	b.WriteString(*cn.TableName)
	b.WriteString(" DROP CONSTRAINT ")
	b.WriteString(*cn.Name)
	b.WriteString(";\n")
	return b.String()
}

//CountQuery -
func (db *PgConstraint) CountQuery(context interface{}) string {
	return ""
}

//Query -
func (db *PgConstraint) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`select conrelid::regclass::varchar as relname, conname, pg_get_constraintdef(c.oid), 
		'','',kcu.column_name,0,ccu.table_name,ccu.column_name,''
		from pg_constraint c
		join pg_namespace n ON n.oid = c.connamespace
		inner JOIN information_schema.key_column_usage AS kcu on kcu.constraint_name = conname
		left outer JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = conname
		where contype in ('f', 'p','c','u') and n.nspname = ANY (current_schemas(false))
		order by conname, relname limit %d offset %d`, args[1].(int), args[0].(int))
}

//FromResult -
func (db *PgConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	return c.GetConstraintFromRow(rows, context)
}
