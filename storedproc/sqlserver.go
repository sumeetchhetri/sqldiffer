package storedproc

import (
	"bytes"
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	pb2 "sqldiffer/protos"
)

//SqlsStoredProcedure -
type SqlsStoredProcedure struct {
	SchemaName string
}

//GenerateNew -
func (db *SqlsStoredProcedure) GenerateNew(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString(*sp.Definition)
	b.WriteString("\n")
	return b.String()
}

//GenerateUpd -
func (db *SqlsStoredProcedure) GenerateUpd(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(sp, context))
	b.WriteString(*sp.Definition)
	return b.String()
}

//GenerateDel -
func (db *SqlsStoredProcedure) GenerateDel(sp *pb2.StoredProcedure, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(c.GetSQLServerPreQuery())
	b.WriteString("DROP PROECURE [")
	b.WriteString(*sp.Name)
	b.WriteString("];\n")
	return b.String()
}

//Query -
func (db *SqlsStoredProcedure) Query(context interface{}) string {
	return `select obj.name as function_name, mod.definition,
			'','',0 from sys.objects obj join sys.sql_modules mod
			on mod.object_id = obj.object_id
			cross apply (select p.name + ' ' + TYPE_NAME(p.user_type_id) + ', ' 
					from sys.parameters p
					where p.object_id = obj.object_id 
							and p.parameter_id != 0 
					for xml path ('') ) par (parameters)
			left join sys.parameters ret
				on obj.object_id = ret.object_id
				and ret.parameter_id = 0
			where obj.type in ('FN', 'TF', 'IF')`
}

//FromResult -
func (db *SqlsStoredProcedure) FromResult(rows *sql.Rows, context interface{}) *pb2.StoredProcedure {
	return c.GetProcedureFromRow(rows, context)
}

//DefineQuery -
func (db *SqlsStoredProcedure) DefineQuery(context interface{}) string {
	return ""
}

//Definition -
func (db *SqlsStoredProcedure) Definition(rows *sql.Rows) string {
	return ""
}
