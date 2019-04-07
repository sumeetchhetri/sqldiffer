package storedprocparam

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//proto "github.com/golang/protobuf/proto"
	//pb2 "sqldiffer/protos"
)

//MysqlStoredProcedureParam -
type MysqlStoredProcedureParam struct {
}

//Query -
func (db MysqlStoredProcedureParam) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT ORDINAL_POSITION,PARAMETER_NAME,DTD_IDENTIFIER,0,0,PARAMETER_MODE,SPECIFIC_NAME,0,'MY' 
	FROM information_schema.parameters WHERE specific_schema = '%s'
	limit %d offset %d`, args[2].(string), args[1].(int), args[0].(int))
}

//FromResult -
func (db MysqlStoredProcedureParam) FromResult(rows *sql.Rows, context interface{}) interface{} {
	return c.GetProcedureParamsFromRow(rows, context)
}
