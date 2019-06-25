package sequence

import (
	"bytes"
	sql "database/sql"
	"strconv"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	c "sqldiffer/common"
	pb2 "sqldiffer/protos"
	//"strings"
)

//PgSequence -
type PgSequence struct {
}

//GenerateNew -
func (db *PgSequence) GenerateNew(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	//TODO
	b.WriteString("\nCREATE SEQUENCE ")
	b.WriteString(*seq.Name)
	if seq.Min != nil {
		b.WriteString(" MINVALUE ")
		b.WriteString(*seq.Min)
	}
	if seq.Max != nil {
		b.WriteString(" MAXVALUE ")
		b.WriteString(*seq.Max)
	}
	if seq.Cycle != nil && *seq.Cycle == "Y" {
		b.WriteString(" CYCLE")
	} else {
		b.WriteString(" NOCYCLE")
	}
	if seq.Cache != nil {
		b.WriteString(" CACHE ")
		b.WriteString(*seq.Cache)
	}
	if seq.Inc != nil {
		b.WriteString(" INCREMENT BY ")
		b.WriteString(*seq.Inc)
	}
	b.WriteString(";\n")
	return b.String()
}

//GenerateUpd -
func (db *PgSequence) GenerateUpd(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(seq, context))
	b.WriteString(db.GenerateNew(seq, context))
	return b.String()
}

//GenerateDel -
func (db *PgSequence) GenerateDel(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nDROP SEQUENCE ")
	b.WriteString(*seq.Name)
	b.WriteString(";\n")
	return b.String()
}

//Query -
func (db *PgSequence) Query(context interface{}) string {
	args := context.([]interface{})
	dbvers, _ := strconv.Atoi(args[1].(string))
	if dbvers < 100000 {
		return `SELECT sequence_name, 
			start_value,
			CASE 
				WHEN cycle_option = 'YES' THEN 'Y' 
				ELSE 'N' 
			END AS cycle, 
			minimum_value, 
			maximum_value, 
			INCREMENT, 
			NULL, 
			NULL, 
			NULL 
		FROM information_schema.sequences`
	}
	return `SELECT sequence_name, 
		start_value,
		CASE 
			WHEN cycle_option = 'YES' THEN 'Y' 
			ELSE 'N' 
		END AS cycle, 
		minimum_value, 
		maximum_value, 
		INCREMENT, 
		seqcache, 
		NULL, 
		NULL 
	FROM information_schema.sequences i
	inner join pg_sequence p on i.sequence_name::regclass = p.seqrelid`
}

//FromResult -
func (db *PgSequence) FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence {
	return c.GetSequenceFromRow(rows)
}

//ExQuery -
func (db *PgSequence) ExQuery(context interface{}) string {
	args := context.([]interface{})
	dbvers, _ := strconv.Atoi(args[1].(string))
	if dbvers < 100000 {
		return "select cache_value from " + args[0].(string)
	}
	return ""
}

//UpdateSequence -
func (db *PgSequence) UpdateSequence(rows *sql.Rows, s *pb2.Sequence) {
	rows.Scan(&s.Cache)
}
