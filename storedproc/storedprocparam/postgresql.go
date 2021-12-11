package storedprocparam

import (
	sql "database/sql"
	"fmt"
	//proto "github.com/golang/protobuf/proto"
	c "github.com/sumeetchhetri/sqldiffer/common"
	//pb2 "github.com/sumeetchhetri/sqldiffer/protos"
	//"strings"
)

//PgStoredProcedureParam -
type PgStoredProcedureParam struct {
}

//Query -
func (db *PgStoredProcedureParam) Query(context interface{}) string {
	args := context.([]interface{})
	spName := args[0].(string)
	spNumPars := args[1].(int32)
	return fmt.Sprintf(`SELECT 0,pg_catalog.pg_get_function_arguments(p.oid),'',0,0,'','',0,'PG'
					FROM pg_catalog.pg_proc p 
					JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace 
					WHERE n.nspname not like 'pg%%' 
					and n.nspname not like 'information_%%'  
					and p.proname = '%s' and p.pronargs = %d`, spName, spNumPars)
}

//FromResult -
func (db *PgStoredProcedureParam) FromResult(rows *sql.Rows, context interface{}) interface{} {
	return c.GetProcedureParamsFromRow(rows, context)
}
