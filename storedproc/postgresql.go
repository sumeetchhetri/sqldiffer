package storedproc

import (
	"bytes"
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	pb2 "sqldiffer/protos"
	"strings"
)

//PgStoredProcedure -
type PgStoredProcedure struct {
	SchemaName string
}

//GenerateNew -
func (db *PgStoredProcedure) GenerateNew(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\n")
	b.WriteString(*sp.Definition)
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgStoredProcedure) GenerateUpd(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(*sp.DropDeclaration)
	b.WriteString(*sp.Definition)
	return b.String()
}

//GenerateDel -
func (db *PgStoredProcedure) GenerateDel(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\n")
	b.WriteString(*sp.DropDeclaration)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *PgStoredProcedure) Query(context interface{}) string {
	return `SELECT DISTINCT quote_ident(p.proname) as function, pg_get_functiondef(p.oid),
			'CREATE OR REPLACE FUNCTION ' || p.proname  || '(' || pg_catalog.pg_get_function_arguments(p.oid) || ');\',
			'DROP FUNCTION ' || p.oid::regprocedure, pronargs 
			FROM pg_catalog.pg_proc p 
			JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace  
			WHERE n.nspname = ANY (current_schemas(false)) 
			and p.proname not like 'pgpool_%' 
			and p.proname not like 'pcp_%' 
			and p.proname not like 'dblink%' 
			and p.proname not like 'uuid_%' 
			and p.proname not like 'pg%' 
			and pg_catalog.pg_get_function_result(p.oid) <> 'trigger'`
}

//FromResult -
func (db *PgStoredProcedure) FromResult(rows *sql.Rows, context interface{}) *pb2.StoredProcedure {
	return c.GetProcedureFromRow(rows, context)
}

//DefineQuery -
func (db *PgStoredProcedure) DefineQuery(context interface{}) string {
	args := context.([]interface{})
	spName := args[0].(string)
	spNumPars := args[1].(int32)
	return fmt.Sprintf(`select pg_get_functiondef(pg_proc.oid) from pg_proc LEFT JOIN pg_namespace n ON n.oid = pg_proc.pronamespace
		where nspname = ANY (current_schemas(false)) and proname = '%s' and pronargs = %d`, spName, spNumPars)
}

//Definition -
func (db *PgStoredProcedure) Definition(rows *sql.Rows) string {
	if rows != nil {
		tmp := ""
		err := rows.Scan(&tmp)
		if err != nil {
			c.Fatal("Error fetching details from database", err)
		}
		utmp := strings.ToUpper(tmp)
		if strings.Index(utmp, "CREATE OR REPLACE FUNCTION "+strings.ToUpper(db.SchemaName)+".") == 0 {
			eprf := "CREATE OR REPLACE FUNCTION " + strings.ToUpper(db.SchemaName) + "."
			tmp = "CREATE OR REPLACE FUNCTION " + tmp[len(eprf):]
		} else if strings.Index(utmp, "CREATE FUNCTION "+strings.ToUpper(db.SchemaName)+".") == 0 {
			eprf := "CREATE FUNCTION " + strings.ToUpper(db.SchemaName) + "."
			tmp = "CREATE FUNCTION " + tmp[len(eprf):]
		}
		return strings.Replace(tmp, "\r\n", "\n", -1)
	}
	return ""
}
