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
	start := args[0].(int)
	end := args[1].(int)
	return fmt.Sprintf(`SELECT outer.* FROM (SELECT position,argument_name,data_type,data_length,data_precision,in_out,object_name,rownum rn,'OR' 
		FROM user_arguments) outer where outer.rn >= %d and outer.rn < %d`, start, end)
}

//FromResult -
func (db *OrclStoredProcedureParam) FromResult(rows *sql.Rows, context interface{}) interface{} {
	return c.GetProcedureParamsFromRow(rows, context)
}
