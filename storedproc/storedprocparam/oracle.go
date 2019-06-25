package storedprocparam

import (
	sql "database/sql"
	"fmt"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	//pb2 "sqldiffer/protos"
	//"strconv"
)

//OrclStoredProcedureParam -
type OrclStoredProcedureParam struct {
}

//Query -
func (db *OrclStoredProcedureParam) Query(context interface{}) string {
	args := context.([]interface{})
	return fmt.Sprintf(`SELECT position,argument_name,data_type,data_length,data_precision,in_out,object_name,0 rn,'OR' 
		FROM user_arguments offset %d rows fetch next %d rows only`, args[0].(int), args[1].(int))
}

//FromResult -
func (db *OrclStoredProcedureParam) FromResult(rows *sql.Rows, context interface{}) interface{} {
	return c.GetProcedureParamsFromRow(rows, context)
}
