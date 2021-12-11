package constraint

import (
	"bytes"
	sql "database/sql"
	"fmt"

	c "github.com/sumeetchhetri/sqldiffer/common"

	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
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
	return fmt.Sprintf(`SELECT conrelid::regclass::varchar as relname, conname, pg_get_constraintdef(con.oid), 
    	'','','',0,'','','' FROM pg_catalog.pg_constraint con
	INNER JOIN pg_catalog.pg_class rel
		ON rel.oid = con.conrelid
	INNER JOIN pg_catalog.pg_namespace nsp
		ON nsp.oid = connamespace
    WHERE contype in ('f', 'p','c','u') and nsp.nspname = ANY (current_schemas(false)) limit %d offset %d`, args[1].(int), args[0].(int))
}

//FromResult -
func (db *PgConstraint) FromResult(rows *sql.Rows, context interface{}) *pb2.Table {
	args := context.([]interface{})
	args = append(args, "postgresql")
	return c.GetConstraintFromRow(rows, args)
}
