package storedprocparam

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	//pb2 "sqldiffer/protos"
)

//SqlsStoredProcedureParam -
type SqlsStoredProcedureParam struct {
}

//Query -
func (db *SqlsStoredProcedureParam) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`select
			'Param_order'  = parameter_id,
			'Parameter_name' = name, 
			'Type'   = type_name(user_type_id),
			'Length'   = max_length,  
			'Prec'   = case when type_name(system_type_id) = 'uniqueidentifier' 
						then precision  
						else OdbcPrec(system_type_id, max_length, precision) end,
			'' mode,
			object_name(object_id),
			0, 'SQ' from sys.parameters WHERE parameter_id > 0
			ORDER BY parameter_id offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *SqlsStoredProcedureParam) FromResult(rows *sql.Rows, context interface{}) interface{} {
	return c.GetProcedureParamsFromRow(rows, context)
}
