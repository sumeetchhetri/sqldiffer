package sequence

import (
	sql "database/sql"
	c "sqldiffer/common"
	//"fmt"
	//proto "github.com/golang/protobuf/proto"
	"bytes"
	pb2 "sqldiffer/protos"
)

//SqlsSequence -
type SqlsSequence struct {
}

//GenerateNew -
func (db *SqlsSequence) GenerateNew(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nCREATE SEQUENCE [")
	b.WriteString(*seq.Name)
	b.WriteString("]")
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
func (db *SqlsSequence) GenerateUpd(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString(db.GenerateDel(seq, context))
	b.WriteString(db.GenerateNew(seq, context))
	return b.String()
}

//GenerateDel -
func (db *SqlsSequence) GenerateDel(seq *pb2.Sequence, context interface{}) string {
	var b bytes.Buffer
	b.WriteString("\nGO\nDROP SEQUENCE [")
	b.WriteString(*seq.Name)
	b.WriteString("];\n")
	return b.String()
}

//Query -
func (db *SqlsSequence) Query(context interface{}) string {
	return `SELECT
			name,
			cast(start_value AS NUMERIC)   AS start_value,
			is_cycling,
			minimum_value,
			maximum_value,
			cast(increment AS NUMERIC)     AS increment,
			cache_size,'',''
		FROM sys.sequences`
}

//FromResult -
func (db *SqlsSequence) FromResult(rows *sql.Rows, context interface{}) *pb2.Sequence {
	return c.GetSequenceFromRow(rows)
}

//ExQuery -
func (db *SqlsSequence) ExQuery(context interface{}) string {
	return ""
}

//UpdateSequence -
func (db *SqlsSequence) UpdateSequence(rows *sql.Rows, s *pb2.Sequence) {
}
