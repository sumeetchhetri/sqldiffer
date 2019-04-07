package common

import (
	sql "database/sql"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"os"
	"reflect"
	pb2 "sqldiffer/protos"
	"strconv"
	"strings"
)

//DbIntf - The Db Interface Type
type DbIntf interface {
	GenerateURL(action *SchemaDiffAction) string
	Preface(dbe *pb2.Db) string
	Create(db *pb2.Db) string
	Connect(db *pb2.Db) string
	CreateSchema(db *pb2.Db) string
}

//TableIntf - The Db Interface Type
type TableIntf interface {
	GenerateNew(tbl *pb2.Table, context interface{}) string
	GenerateUpd(tbl *pb2.Table, context interface{}) string
	GenerateDel(tbl *pb2.Table, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Table
}

//ColumnIntf - The Db Interface Type
type ColumnIntf interface {
	GenerateNew(tbl *pb2.Column, context interface{}) string
	GenerateUpd(tbl *pb2.Column, context interface{}) string
	GenerateDel(tbl *pb2.Column, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Table
}

//ConstraintIntf - The Db Interface Type
type ConstraintIntf interface {
	GenerateNew(tbl *pb2.Constraint, context interface{}) string
	GenerateUpd(tbl *pb2.Constraint, context interface{}) string
	GenerateDel(tbl *pb2.Constraint, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Table
}

//IndexIntf - The Db Interface Type
type IndexIntf interface {
	GenerateNew(tbl *pb2.Index, context interface{}) string
	GenerateUpd(tbl *pb2.Index, context interface{}) string
	GenerateDel(tbl *pb2.Index, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Index
}

//TriggerIntf -
type TriggerIntf interface {
	GenerateNew(tbl *pb2.Trigger, context interface{}) string
	GenerateUpd(tbl *pb2.Trigger, context interface{}) string
	GenerateDel(tbl *pb2.Trigger, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Trigger
	DefineQuery(context interface{}) string
	GetDefinition(rows *sql.Rows) string
}

//ViewIntf - The Db Interface Type
type ViewIntf interface {
	GenerateNew(tbl *pb2.View, context interface{}) string
	GenerateUpd(tbl *pb2.View, context interface{}) string
	GenerateDel(tbl *pb2.View, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.View
}

//SequenceIntf - The Db Interface Type
type SequenceIntf interface {
	GenerateNew(tbl *pb2.Sequence, context interface{}) string
	GenerateUpd(tbl *pb2.Sequence, context interface{}) string
	GenerateDel(tbl *pb2.Sequence, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence
	ExQuery(name string) string
	UpdateSequence(rows *sql.Rows, s *pb2.Sequence)
}

//StoredProcedureIntf -
type StoredProcedureIntf interface {
	GenerateNew(tbl *pb2.StoredProcedure, context interface{}) string
	GenerateUpd(tbl *pb2.StoredProcedure, context interface{}) string
	GenerateDel(tbl *pb2.StoredProcedure, context interface{}) string
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) *pb2.StoredProcedure
	DefineQuery(context interface{}) string
	Definition(rows *sql.Rows) string
}

//StoredProcedureParamIntf - The Db Interface Type
type StoredProcedureParamIntf interface {
	Query(context interface{}) string
	FromResult(rows *sql.Rows, context interface{}) interface{}
}

//SchemaDiffAction -
type SchemaDiffAction struct {
	Host                      *string
	Port                      int32
	DatabaseName              *string
	User                      *string
	Password                  *string
	SchemaName                *string
	FileName                  *string
	Parallel                  bool
	IsDiffNeeded              bool
	DatabaseType              *string
	SourceSchemaFile          *string
	TargetSchemaFile          *string
	SingleDiffFile            bool
	TargetDatabaseType        *string
	TargetDatabaseName        *string
	TargetSchemaName          *string
	ReverseDiffNeeded         bool
	DiffFileName              *string
	DiffOptions               *string
	DuplicateProcNamesAllowed bool
	Db                        DbIntf
	Table                     TableIntf
	Column                    ColumnIntf
	Constraint                ConstraintIntf
	Index                     IndexIntf
	Trigger                   TriggerIntf
	View                      ViewIntf
	Sequence                  SequenceIntf
	StoredProcedure           StoredProcedureIntf
	StoredProcedureParam      StoredProcedureParamIntf
	Procs                     map[string]*pb2.StoredProcedure
	Tables                    map[string]*pb2.Table
	TablesP                   map[string]bool
	Indexes                   map[string][]*pb2.Index
	Triggers                  map[string][]*pb2.Trigger
	Constraints               map[string]*pb2.Constraint
}

//Fatal -
func Fatal(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(2)
}

//ColumnEq -
func ColumnEq(t1, t2 pb2.Column) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.Type == nil && t2.Type != nil) ||
		(t1.DefVal == nil && t2.DefVal != nil) || (t1.TableName == nil && t2.TableName != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.Type != nil && t2.Type == nil) ||
		(t1.DefVal != nil && t2.DefVal == nil) || (t1.TableName != nil && t2.TableName == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.Type != *t2.Type || (t1.DefVal != nil && t2.DefVal != nil && *t1.DefVal != *t2.DefVal) ||
		*t1.TableName != *t2.TableName {
		return false
	}
	return true
}

//TriggerEq -
func TriggerEq(t1, t2 pb2.Trigger) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.TableName == nil && t2.TableName != nil) ||
		(t1.When == nil && t2.When != nil) || (t1.Action == nil && t2.Action != nil) ||
		(t1.Function == nil && t2.Function != nil) || (t1.FunctionDef == nil && t2.FunctionDef != nil) ||
		(t1.Definition == nil && t2.Definition != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.TableName != nil && t2.TableName == nil) ||
		(t1.When != nil && t2.When == nil) || (t1.Action != nil && t2.Action == nil) ||
		(t1.Function != nil && t2.Function == nil) || (t1.FunctionDef != nil && t2.FunctionDef == nil) ||
		(t1.Definition != nil && t2.Definition == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.TableName != *t2.TableName || *t1.When != *t2.When ||
		*t1.Action != *t2.Action || *t1.FunctionDef != *t2.FunctionDef ||
		!StringEqualsIgnSpace(*t1.Function, *t2.Function) ||
		!StringEqualsIgnSpace(*t1.Definition, *t2.Definition) {
		return false
	}
	return true
}

//IndexEq -
func IndexEq(t1, t2 pb2.Index) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.TableName == nil && t2.TableName != nil) ||
		(t1.Definition == nil && t2.Definition != nil) || (t1.Props == nil && t2.Props != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.TableName != nil && t2.TableName == nil) ||
		(t1.Definition != nil && t2.Definition == nil) || (t1.Props != nil && t2.Props == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.TableName != *t2.TableName ||
		(t1.Definition != nil && t2.Definition != nil && *t1.Definition != *t2.Definition) ||
		!reflect.DeepEqual(t1.Props, t2.Props) {
		return false
	}
	return true
}

//ConstraintEq -
func ConstraintEq(t1, t2 pb2.Constraint) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.TableName == nil && t2.TableName != nil) ||
		(t1.Definition == nil && t2.Definition != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.TableName != nil && t2.TableName == nil) ||
		(t1.Definition != nil && t2.Definition == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.TableName != *t2.TableName || *t1.Definition != *t2.Definition {
		return false
	}
	return true
}

//SequenceEq -
func SequenceEq(t1, t2 pb2.Sequence) bool {
	if &t1 == &t2 {
		return true
	}
	if t1.Name == nil && t2.Name != nil {
		return false
	}
	if t1.Name != nil && t2.Name == nil {
		return false
	}
	if *t1.Name != *t2.Name {
		return false
	}
	return true
}

//ViewEq -
func ViewEq(t1, t2 pb2.View) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.Type == nil && t2.Type != nil) ||
		(t1.Definition == nil && t2.Definition != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.Type != nil && t2.Type == nil) ||
		(t1.Definition != nil && t2.Definition == nil) {
		return false
	}
	if *t1.Name != *t2.Name || (t1.Type != nil && t2.Type != nil && *t1.Type != *t2.Type) ||
		!StringEqualsIgnSpace(*t1.Definition, *t2.Definition) {
		return false
	}
	return true
}

//StoredProcedureParamEq -
func StoredProcedureParamEq(t1, t2 pb2.StoredProcedureParam) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.Type == nil && t2.Type != nil) ||
		(t1.Mode == nil && t2.Mode != nil) || (t1.Position == nil && t2.Position != nil) ||
		(t1.DefVal == nil && t2.DefVal != nil) || (t1.ProcName == nil && t2.ProcName != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.Type != nil && t2.Type == nil) ||
		(t1.Mode != nil && t2.Mode == nil) || (t1.Position != nil && t2.Position == nil) ||
		(t1.DefVal != nil && t2.DefVal == nil) || (t1.ProcName != nil && t2.ProcName == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.Type != *t2.Type || *t1.Mode != *t2.Mode ||
		*t1.Position != *t2.Position || (t1.DefVal != nil && t2.DefVal != nil && *t1.DefVal != *t2.DefVal) ||
		*t1.ProcName != *t2.ProcName {
		return false
	}
	return true
}

//StoredProcedureEq -
func StoredProcedureEq(t1, t2 pb2.StoredProcedure) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.Declaration == nil && t2.Declaration != nil) ||
		(t1.DropDeclaration == nil && t2.DropDeclaration != nil) || (t1.Definition == nil && t2.Definition != nil) ||
		(t1.NumParams == nil && t2.NumParams != nil) || (t1.Params == nil && t2.Params != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.Declaration != nil && t2.Declaration == nil) ||
		(t1.DropDeclaration != nil && t2.DropDeclaration == nil) || (t1.Definition != nil && t2.Definition == nil) ||
		(t1.NumParams != nil && t2.NumParams == nil) || (t1.Params != nil && t2.Params == nil) {
		return false
	}
	if *t1.Name != *t2.Name || *t1.Declaration != *t2.Declaration ||
		!StringEqualsIgnSpace(*t1.DropDeclaration, *t2.DropDeclaration) ||
		!StringEqualsIgnSpace(*t1.Definition, *t2.Definition) ||
		*t1.NumParams != *t2.NumParams {
		return false
	}
	if t1.Params != nil && t2.Params != nil {
		for _, v1 := range t1.Params {
			fl := false
			for _, v2 := range t2.Params {
				if StoredProcedureParamEq(*v1, *v2) {
					fl = true
					break
				}
			}
			if !fl {
				return false
			}
		}
	}
	return true
}

//TableEq -
func TableEq(t1, t2 pb2.Table) bool {
	if &t1 == &t2 {
		return true
	}
	if (t1.Name == nil && t2.Name != nil) || (t1.IsTemp == nil && t2.IsTemp != nil) ||
		(t1.Columns == nil && t2.Columns != nil) || (t1.Triggers == nil && t2.Triggers != nil) ||
		(t1.Indexes == nil && t2.Indexes != nil) || (t1.Constraints == nil && t2.Constraints != nil) {
		return false
	}
	if (t1.Name != nil && t2.Name == nil) || (t1.IsTemp != nil && t2.IsTemp == nil) ||
		(t1.Columns != nil && t2.Columns == nil) || (t1.Triggers != nil && t2.Triggers == nil) ||
		(t1.Indexes != nil && t2.Indexes == nil) || (t1.Constraints != nil && t2.Constraints == nil) {
		return false
	}
	if *t1.Name != *t2.Name || (t1.IsTemp != nil && t2.IsTemp != nil && *t1.IsTemp != *t2.IsTemp) {
		return false
	}
	if t1.Columns != nil && t2.Columns != nil {
		for _, v1 := range t1.Columns {
			fl := false
			for _, v2 := range t2.Columns {
				if ColumnEq(*v1, *v2) {
					fl = true
					break
				}
			}
			if !fl {
				return false
			}
		}
	}
	if t1.Triggers != nil && t2.Triggers != nil {
		for _, v1 := range t1.Triggers {
			fl := false
			for _, v2 := range t2.Triggers {
				if TriggerEq(*v1, *v2) {
					fl = true
					break
				}
			}
			if !fl {
				return false
			}
		}
	}
	if t1.Indexes != nil && t2.Indexes != nil {
		for _, v1 := range t1.Indexes {
			fl := false
			for _, v2 := range t2.Indexes {
				if IndexEq(*v1, *v2) {
					fl = true
					break
				}
			}
			if !fl {
				return false
			}
		}
	}
	if t1.Constraints != nil && t2.Constraints != nil {
		for _, v1 := range t1.Constraints {
			fl := false
			for _, v2 := range t2.Constraints {
				if ConstraintEq(*v1, *v2) {
					fl = true
					break
				}
			}
			if !fl {
				return false
			}
		}
	}
	return true
}

//StringEqualsIgnSpace -
func StringEqualsIgnSpace(v1, v2 string) bool {
	v1 = strings.Replace(v1, "\n", "", -1)
	v1 = strings.Replace(v1, "\r", "", -1)
	v1 = strings.Replace(v1, "[\t ]+", "", -1)
	v2 = strings.Replace(v2, "\n", "", -1)
	v2 = strings.Replace(v2, "\r", "", -1)
	v2 = strings.Replace(v2, "[\t ]+", "", -1)
	return v1 == v2
}

//GetColumnFromRow -
func GetColumnFromRow(rows *sql.Rows, context interface{}) *pb2.Table {
	sp := pb2.Column{}

	var (
		defv   sql.NullString
		nn     sql.NullString
		len    sql.NullInt64
		prec   sql.NullInt64
		scl    sql.NullInt64
		rowNum uint64
	)

	err := rows.Scan(&sp.Pos, &sp.Name, &nn, &sp.Type, &len, &prec, &scl, &defv, &sp.TableName, &rowNum)
	if err != nil {
		Fatal("Error fetching column details from database", err)
	}

	if nn.Valid && nn.String == "Y" {
		sp.Notnull = proto.Bool(true)
	}
	if prec.Valid {
		sp.Precision = proto.Int64(prec.Int64)
	}
	if scl.Valid {
		sp.Scale = proto.Int64(scl.Int64)
	}
	if len.Valid {
		sp.Length = proto.Int64(len.Int64)
	}
	if sp.Length != nil && strings.Index(*sp.Type, "(") == -1 && strings.Index(strings.ToLower(*sp.Type), "char") != -1 {
		*sp.Type = *sp.Type + "(" + strconv.FormatInt(len.Int64, 10) + ")"
	} else if sp.Length != nil && sp.Scale != nil && (strings.ToLower(*sp.Type) == "decimal" ||
		strings.ToLower(*sp.Type) == "number") {
		*sp.Type = *sp.Type + "(" + strconv.FormatInt(len.Int64, 10) + "," + strconv.FormatInt(scl.Int64, 10) + ")"
	}
	if defv.Valid {
		defv.String = strings.TrimSpace(defv.String)
		if strings.Index(defv.String, ".\"NEXTVAL\"") != -1 || strings.Index(defv.String, ".\"nextval\"") != -1 {
			defv.String = strings.Replace(defv.String, "\"", "", -1)
		}
	}
	sp.DefVal = proto.String(defv.String)

	tbls := context.(map[string]*pb2.Table)

	tsp, ok := tbls[*sp.TableName]
	if ok {
		tsp.Columns = append(tsp.Columns, &sp)
		return tsp
	}
	return nil
}

//GetSequenceFromRow -
func GetSequenceFromRow(rows *sql.Rows) *pb2.Sequence {
	sp := pb2.Sequence{}

	var (
		cache sql.NullString
		order sql.NullString
		typ   sql.NullString
	)
	err := rows.Scan(&sp.Name, &sp.DefVal, &sp.Cycle, &sp.Min, &sp.Max, &sp.Inc, &cache, &order, &typ)
	if err != nil {
		Fatal("Error fetching sequence details from database", err)
	}

	if cache.Valid {
		sp.Cache = proto.String(cache.String)
	}
	if order.Valid {
		sp.Order = proto.String(order.String)
	}
	if typ.Valid {
		sp.Type = proto.String(typ.String)
	}

	return &sp
}

//GetTableFromRow -
func GetTableFromRow(rows *sql.Rows, context interface{}) *pb2.Table {
	sp := pb2.Table{}

	var b string
	err := rows.Scan(&sp.Name, &b)
	if err != nil {
		Fatal("Error fetching table details from database", err)
	}

	if b == "Y" {
		sp.IsTemp = proto.Bool(true)
	} else {
		sp.IsTemp = proto.Bool(false)
	}

	return &sp
}

//GetTriggerFromRow -
func GetTriggerFromRow(rows *sql.Rows, context interface{}) *pb2.Trigger {
	sp := pb2.Trigger{}

	err := rows.Scan(&sp.Name, &sp.TableName, &sp.When, &sp.Action, &sp.Function, &sp.FunctionDef, &sp.Definition)
	if err != nil {
		Fatal("Error fetching details from database", err)
	}

	return &sp
}

//MergeDuplicates -
func MergeDuplicates(triggers []*pb2.Trigger) []*pb2.Trigger {
	utriggers := make([]*pb2.Trigger, 0)
	utrgs := make(map[string][]pb2.Trigger)

	for _, t := range triggers {
		_, ok := utrgs[*t.Name+*t.When+*t.Definition]
		if !ok {
			utrgs[*t.Name+*t.When+*t.Definition] = make([]pb2.Trigger, 0)
		}
		utrgs[*t.Name+*t.When+*t.Definition] = append(utrgs[*t.Name+*t.When+*t.Definition], *t)
	}
	for _, trgs := range utrgs {
		actions := make([]string, 0)
		for _, t := range trgs {
			actions = append(actions, *t.Action)
		}
		tgt := trgs[0]
		*tgt.Action = strings.Join(actions, " OR ")
		utriggers = append(utriggers, &tgt)
	}
	return utriggers
}

//GetIndexFromRow -
func GetIndexFromRow(rows *sql.Rows, context interface{}) *pb2.Index {
	sp := pb2.Index{}

	sp.Props = make(map[string]string)

	var (
		arg1 string
		arg2 string
	)

	err := rows.Scan(&sp.TableName, &sp.Name, &sp.Definition, &sp.Columns, &arg1, &arg2)
	if err != nil {
		Fatal("Error fetching details from database", err)
	}

	if arg1 != "" {
		sp.Props["Type"] = arg1
	}
	if arg2 != "" {
		sp.Props["Unique"] = arg2
	}

	return &sp
}

//GetConstraintFromRow -
func GetConstraintFromRow(rows *sql.Rows, context interface{}) *pb2.Table {
	sp := pb2.Constraint{}
	args := context.([]interface{})
	consColMap := args[1].(map[string]*pb2.Constraint)

	var rn int32
	var col string

	err := rows.Scan(&sp.TableName, &sp.Name, &sp.Definition, &sp.Condition, &sp.Type,
		&col, &rn, &sp.TargetTableName, &sp.TargetColumnName, &sp.TableView)
	if err != nil {
		Fatal("Error fetching Constraint Row details from database", err)
	}

	if col != "" {
		exsp, ok := consColMap[*sp.Name]
		if ok {
			if strings.Index(col, ",") == -1 {
				exsp.Columns = append(exsp.Columns, col)
			} else {
				cols := strings.Split(col, ",")
				for _, c := range cols {
					exsp.Columns = append(exsp.Columns, strings.TrimSpace(c))
				}
			}
			return nil
		}
		sp.Columns = make([]string, 0)
		if strings.Index(col, ",") == -1 {
			sp.Columns = append(sp.Columns, col)
		} else {
			cols := strings.Split(col, ",")
			for _, c := range cols {
				sp.Columns = append(sp.Columns, strings.TrimSpace(c))
			}
		}
	}

	tbls := args[0].(map[string]*pb2.Table)

	tsp, ok := tbls[*sp.TableName]
	if ok {
		tsp.Constraints = append(tsp.Constraints, &sp)
		return tsp
	}
	return nil
}

//GetProcedureFromRow -
func GetProcedureFromRow(rows *sql.Rows, context interface{}) *pb2.StoredProcedure {
	sp := pb2.StoredProcedure{}
	schemaName := context.(string)

	err := rows.Scan(&sp.Name, &sp.Definition, &sp.Declaration, &sp.DropDeclaration, &sp.NumParams)
	if err != nil {
		Fatal("Error fetching details from database", err)
	}

	tmp := *sp.Definition
	utmp := strings.ToUpper(tmp)
	if strings.Index(utmp, "CREATE OR REPLACE FUNCTION "+strings.ToUpper(schemaName)+".") == 0 {
		eprf := "CREATE OR REPLACE FUNCTION " + strings.ToUpper(schemaName) + "."
		tmp = "CREATE OR REPLACE FUNCTION " + tmp[len(eprf):]
	} else if strings.Index(utmp, "CREATE FUNCTION "+strings.ToUpper(schemaName)+".") == 0 {
		eprf := "CREATE FUNCTION " + strings.ToUpper(schemaName) + "."
		tmp = "CREATE FUNCTION " + tmp[len(eprf):]
	}
	*sp.Definition = tmp //strings.Replace(tmp, "\r\n", "\n", -1)
	return &sp
}

//GetProcedureParamsFromRow -
func GetProcedureParamsFromRow(rows *sql.Rows, context interface{}) interface{} {
	args := context.([]interface{})
	pgsa := pb2.StoredProcedureParam{}

	var (
		pos      sql.NullInt64
		name     sql.NullString
		typ      sql.NullString
		typadd1  sql.NullInt64
		typadd2  sql.NullInt64
		mode     sql.NullString
		procName sql.NullString
		rowNmum  uint64
		dtyp     string
	)

	err := rows.Scan(&pos, &name, &typ, &typadd1, &typadd2, &mode, &procName, &rowNmum, &dtyp)
	if err != nil {
		Fatal("Error fetching details from database", err)
	}

	if dtyp == "PG" {
		spName := args[0].(string)
		columns := make([]*pb2.StoredProcedureParam, 0)

		ar := strings.Split(name.String, ",")
		for _, v := range ar {
			if strings.TrimSpace(v) == "" {
				continue
			}
			pgsa := pb2.StoredProcedureParam{}
			pgsa.Mode = proto.String("IN")
			if strings.Index(strings.ToLower(v), "inout ") == 0 {
				pgsa.Mode = proto.String("INOUT")
				v = v[4:]
			} else if strings.Index(strings.ToLower(v), "out ") == 0 {
				pgsa.Mode = proto.String("OUT")
				v = v[6:]
			} else {
				v = strings.TrimSpace(v)
			}

			pnm := ""
			if strings.Index(v, " ") != -1 {
				pnm = v[0:strings.Index(v, " ")]
			}
			v = v[strings.Index(v, " ")+1:]
			pgsa.Name = proto.String(pnm)
			pgsa.Position = proto.Int32(int32(len(columns) + 1))

			typ := v
			if strings.Index(strings.ToLower(v), " default ") != -1 {
				typ = v[0:strings.Index(strings.ToLower(v), " default ")]
				v = v[strings.Index(strings.ToLower(v), " default ")+9:]
				pgsa.DefVal = proto.String(v)
			}
			pgsa.Type = proto.String(typ)
			pgsa.ProcName = proto.String(spName)
			columns = append(columns, &pgsa)
		}
		return columns
	}

	pgsa.Position = proto.Int32(int32(pos.Int64))
	pgsa.Name = proto.String(name.String)
	pgsa.Type = proto.String(typ.String)
	if dtyp == "OR" {
		if typadd1.Valid {
			*pgsa.Type = *pgsa.Type + "(" + strconv.FormatInt(typadd1.Int64, 10)
			if typadd2.Valid {
				*pgsa.Type = *pgsa.Type + ", " + strconv.FormatInt(typadd2.Int64, 10)
			}
			*pgsa.Type = *pgsa.Type + ")"
		}
	}
	pgsa.Mode = proto.String(mode.String)
	pgsa.ProcName = proto.String(procName.String)

	return &pgsa
}

//GetViewFromRow -
func GetViewFromRow(rows *sql.Rows, context interface{}) *pb2.View {
	sp := pb2.View{}

	err := rows.Scan(&sp.Name, &sp.Definition)
	if err != nil {
		Fatal("Error fetching details from database", err)
	}
	return &sp
}

//GetSQLServerPreQuery -
func GetSQLServerPreQuery() string {
	return "\nSET ANSI_NULLS ON\nGO\nSET QUOTED_IDENTIFIER ON\nGO\n"
}
