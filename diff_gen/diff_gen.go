package main

import (
	"bytes"
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"google.golang.org/protobuf/proto"

	//"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"time"

	c "github.com/sumeetchhetri/sqldiffer/common"
	db "github.com/sumeetchhetri/sqldiffer/db"
	pb2 "github.com/sumeetchhetri/sqldiffer/protos"
	sq "github.com/sumeetchhetri/sqldiffer/sequence"
	sp "github.com/sumeetchhetri/sqldiffer/storedproc"
	spp "github.com/sumeetchhetri/sqldiffer/storedproc/storedprocparam"
	tb "github.com/sumeetchhetri/sqldiffer/table"
	co "github.com/sumeetchhetri/sqldiffer/table/column"
	cn "github.com/sumeetchhetri/sqldiffer/table/constraint"
	in "github.com/sumeetchhetri/sqldiffer/table/index"
	tr "github.com/sumeetchhetri/sqldiffer/table/trigger"
	vw "github.com/sumeetchhetri/sqldiffer/view"
)

var opts struct {
	SourceSchemaFile   string `short:"f" long:"src-scf" description:"The source schema file name" required:"true"`
	TargetSchemaFile   string `short:"d" long:"tgt-scf" description:"The target schema file name" required:"true"`
	SingleDiffFile     bool   `short:"n" long:"one-file" description:"Whether a single diff file is needed or multiple?"`
	TargetDatabaseType string `short:"t" long:"tgt-dbt" description:"The target database type" required:"true" choice:"postgres" choice:"oracle" choice:"mysql" choice:"sqlserver"`
	TargetDatabaseName string `short:"m" long:"tgt-dbn" description:"The target database name" required:"true"`
	TargetSchemaName   string `short:"s" long:"tgt-dbs" description:"The target database schema name"`
	ReverseDiffNeeded  bool   `short:"r" long:"rdiff" description:"Whether a reverse diff is needed or not?"`
	DiffFileName       string `short:"o" long:"out-file" description:"Diff File Name"`
	DiffOptions        string `short:"p" long:"options" description:"Extra Diff options"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println("Failed to parse args: %!", err)
			os.Exit(2)
		}
	}

	action := c.SchemaDiffAction{
		SourceSchemaFile:   &opts.SourceSchemaFile,
		TargetSchemaFile:   &opts.TargetSchemaFile,
		SingleDiffFile:     opts.SingleDiffFile,
		TargetDatabaseType: &opts.TargetDatabaseType,
		TargetDatabaseName: &opts.TargetDatabaseName,
		TargetSchemaName:   &opts.TargetSchemaName,
		ReverseDiffNeeded:  opts.ReverseDiffNeeded,
		DiffFileName:       &opts.DiffFileName,
		DiffOptions:        &opts.DiffOptions,
	}
	generate(&action)
}

func generate(action *c.SchemaDiffAction) {
	if action.DiffFileName == nil || *action.DiffFileName == "" {
		action.DiffFileName = proto.String("diff_" + (fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))))
	}
	if *action.TargetDatabaseType == "postgres" {
		action.Db = &db.PgDb{}
		action.Table = &tb.PgTable{}
		action.Column = &co.PgColumn{}
		action.Constraint = &cn.PgConstraint{}
		action.Index = &in.PgIndex{}
		action.Trigger = &tr.PgTrigger{SchemaName: *action.TargetSchemaName}
		action.StoredProcedure = &sp.PgStoredProcedure{SchemaName: *action.TargetSchemaName}
		action.StoredProcedureParam = &spp.PgStoredProcedureParam{}
		action.Sequence = &sq.PgSequence{}
		action.View = &vw.PgView{}
		action.DuplicateProcNamesAllowed = true
		if action.SchemaName == nil || *action.SchemaName == "" {
			action.SchemaName = proto.String("public")
		}
	} else if *action.TargetDatabaseType == "oracle" {
		action.Db = &db.OrclDb{}
		action.Table = &tb.OrclTable{}
		action.Column = &co.OrclColumn{}
		action.Constraint = &cn.OrclConstraint{}
		action.Index = &in.OrclIndex{}
		action.Trigger = &tr.OrclTrigger{SchemaName: *action.TargetSchemaName}
		action.StoredProcedure = &sp.OrclStoredProcedure{SchemaName: *action.TargetSchemaName}
		action.StoredProcedureParam = &spp.OrclStoredProcedureParam{}
		action.Sequence = &sq.OrclSequence{}
		action.View = &vw.OrclView{}
	} else if *action.TargetDatabaseType == "mysql" {
		action.Db = &db.MysqlDb{}
		action.Table = &tb.MysqlTable{}
		action.Column = &co.MysqlColumn{}
		action.Constraint = &cn.MysqlConstraint{}
		action.Index = &in.MysqlIndex{}
		action.Trigger = &tr.MysqlTrigger{SchemaName: *action.TargetSchemaName}
		action.StoredProcedure = &sp.MysqlStoredProcedure{SchemaName: *action.TargetSchemaName}
		action.StoredProcedureParam = &spp.MysqlStoredProcedureParam{}
		action.Sequence = &sq.MysqlSequence{}
		action.View = &vw.MysqlView{}
	} else if *action.TargetDatabaseType == "sqlserver" {
		action.Db = &db.SqlsDb{}
		action.Table = &tb.SqlsTable{}
		action.Column = &co.SqlsColumn{}
		action.Constraint = &cn.SqlsConstraint{}
		action.Index = &in.SqlsIndex{}
		action.Trigger = &tr.SqlsTrigger{SchemaName: *action.TargetSchemaName}
		action.StoredProcedure = &sp.SqlsStoredProcedure{SchemaName: *action.TargetSchemaName}
		action.StoredProcedureParam = &spp.SqlsStoredProcedureParam{}
		action.Sequence = &sq.SqlsSequence{}
		action.View = &vw.SqlsView{}
	}

	b, err := ioutil.ReadFile(*action.SourceSchemaFile)
	if err != nil {
		c.Fatal("Unable to read source schema file", err)
	}

	dbsrc := pb2.Db{}
	err = proto.Unmarshal(b, &dbsrc)
	if err != nil {
		c.Fatal("Unable to un-marshal source schema file contents", err)
	}

	dbdst := pb2.Db{}
	b, err = ioutil.ReadFile(*action.TargetSchemaFile)
	if err == nil {
		err = proto.Unmarshal(b, &dbdst)
		if err != nil {
			c.Fatal("Unable to un-marshal target schema file contents", err)
		}
	} else {
		if *action.TargetDatabaseType == "postgres" {
			dbdst.SchemaName = proto.String("public")
		}
	}

	generateDiff(&dbsrc, &dbdst, action, false)
	if action.ReverseDiffNeeded {
		generateDiff(&dbdst, &dbsrc, action, true)
	}
}

func generateDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, reverseDiff bool) {
	var dbb, seqb, tabb, spcb, trgb, mscb bytes.Buffer

	if !action.SingleDiffFile {
		dbb.WriteString(action.Db.Preface(dbdst))
		seqb.WriteString(action.Db.Preface(dbdst))
		tabb.WriteString(action.Db.Preface(dbdst))
		spcb.WriteString(action.Db.Preface(dbdst))
		trgb.WriteString(action.Db.Preface(dbdst))
		mscb.WriteString(action.Db.Preface(dbdst))
	}

	generateTableDiff(dbsrc, dbdst, action, &tabb, &mscb)
	generateStoreProcedureDiff(dbsrc, dbdst, action, &spcb)
	generateSequenceDiff(dbsrc, dbdst, action, &seqb)
	generateTriggerDiff(dbsrc, dbdst, action, &trgb)
	generateViewDiff(dbsrc, dbdst, action, &tabb)

	if !reverseDiff {
		if dbdst.Name == nil {
			dbdst.Name = proto.String(*action.TargetDatabaseName)
			dbb.WriteString(action.Db.Create(dbdst))
		}
		if dbdst.SchemaName == nil {
			dbdst.SchemaName = proto.String(*action.TargetSchemaName)
			dbb.WriteString(action.Db.Connect(dbdst))
			dbb.WriteString(action.Db.CreateSchema(dbdst))
		}
	}

	if !action.SingleDiffFile {
		if reverseDiff {
			os.Remove(*action.DiffFileName + "_db_r.sql")
			os.Remove(*action.DiffFileName + "_seq_r.sql")
			os.Remove(*action.DiffFileName + "_tab_r.sql")
			os.Remove(*action.DiffFileName + "_spc_r.sql")
			os.Remove(*action.DiffFileName + "_trg_r.sql")
			os.Remove(*action.DiffFileName + "_msc_r.sql")

			if dbb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_db_r.sql", dbb.Bytes(), 0644)
			}
			if seqb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_seq_r.sql", seqb.Bytes(), 0644)
			}
			if tabb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_tab_r.sql", tabb.Bytes(), 0644)
			}
			if spcb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_spc_r.sql", spcb.Bytes(), 0644)
			}
			if trgb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_trg_r.sql", trgb.Bytes(), 0644)
			}
			if mscb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_msc_r.sql", mscb.Bytes(), 0644)
			}
		} else {
			os.Remove(*action.DiffFileName + "_db.sql")
			os.Remove(*action.DiffFileName + "_seq.sql")
			os.Remove(*action.DiffFileName + "_tab.sql")
			os.Remove(*action.DiffFileName + "_spc.sql")
			os.Remove(*action.DiffFileName + "_trg.sql")
			os.Remove(*action.DiffFileName + "_msc.sql")

			if dbb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_db.sql", dbb.Bytes(), 0644)
			}
			if seqb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_seq.sql", seqb.Bytes(), 0644)
			}
			if tabb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_tab.sql", tabb.Bytes(), 0644)
			}
			if spcb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_spc.sql", spcb.Bytes(), 0644)
			}
			if trgb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_trg.sql", trgb.Bytes(), 0644)
			}
			if mscb.Len() > 0 {
				ioutil.WriteFile(*action.DiffFileName+"_msc.sql", mscb.Bytes(), 0644)
			}
		}
	} else {
		dbb.WriteString(seqb.String())
		dbb.WriteString(tabb.String())
		dbb.WriteString(spcb.String())
		dbb.WriteString(trgb.String())
		dbb.WriteString(mscb.String())

		if reverseDiff {
			os.Remove(*action.DiffFileName + "_r.sql")
			ioutil.WriteFile(*action.DiffFileName+"_r.sql", dbb.Bytes(), 0644)
		} else {
			os.Remove(*action.DiffFileName + ".sql")
			ioutil.WriteFile(*action.DiffFileName+".sql", dbb.Bytes(), 0644)
		}
	}
}

func generateTableDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, tabb *bytes.Buffer, mscb *bytes.Buffer) {
	for _, tsrc := range dbsrc.Tables {
		var tdst *pb2.Table
		for _, tmpb := range dbdst.Tables {
			if *tsrc.Name == *tmpb.Name {
				tdst = tmpb
				break
			}
		}
		donem := make(map[string]bool)
		if tdst == nil {
			tabb.WriteString(action.Table.GenerateNew(tsrc, nil))
			for _, csrc := range tsrc.Indexes {
				mscb.WriteString(action.Index.GenerateNew(csrc, nil))
			}
			for _, csrc := range tsrc.Constraints {
				_, ok := donem[*csrc.Name+*csrc.Definition]
				if !ok {
					donem[*csrc.Name+*csrc.Definition] = true
					mscb.WriteString(action.Constraint.GenerateNew(csrc, nil))
				}
			}
		} else if !c.TableEq(tsrc, tdst) {
			for _, csrc := range tsrc.Columns {
				var cdst *pb2.Column
				for _, tmp := range tdst.Columns {
					if *csrc.Name == *tmp.Name {
						cdst = tmp
						break
					}
				}
				if cdst == nil {
					tabb.WriteString(action.Column.GenerateNew(csrc, nil))
				} else if !c.ColumnEq(csrc, cdst) {
					tabb.WriteString(action.Column.GenerateUpd(csrc, cdst))
				}
			}
			for _, cdst := range tsrc.Columns {
				var csrc *pb2.Column
				for _, tmp := range tsrc.Columns {
					if *cdst.Name == *tmp.Name {
						csrc = tmp
						break
					}
				}
				if csrc == nil {
					tabb.WriteString(action.Column.GenerateDel(cdst, nil))
				}
			}

			for _, csrc := range tsrc.Constraints {
				var cdst *pb2.Constraint
				for _, tmp := range tdst.Constraints {
					if *csrc.Name == *tmp.Name {
						cdst = tmp
						break
					}
				}
				if cdst == nil {
					_, ok := donem[*csrc.Name+*csrc.Definition]
					if !ok {
						donem[*csrc.Name+*csrc.Definition] = true
						mscb.WriteString(action.Constraint.GenerateNew(csrc, nil))
					}
				} else if !c.ConstraintEq(csrc, cdst) {
					_, ok := donem[*csrc.Name+*csrc.Definition]
					if !ok {
						donem[*csrc.Name+*csrc.Definition] = true
						mscb.WriteString(action.Constraint.GenerateDel(cdst, nil))
						mscb.WriteString(action.Constraint.GenerateNew(csrc, nil))
					}
				}
			}
			for _, cdst := range tdst.Constraints {
				var csrc *pb2.Constraint
				for _, tmpb := range tsrc.Constraints {
					if *cdst.Name == *tmpb.Name {
						csrc = tmpb
						break
					}
				}
				if csrc == nil {
					mscb.WriteString(action.Constraint.GenerateDel(cdst, nil))
				}
			}

			for _, csrc := range tsrc.Indexes {
				var cdst *pb2.Index
				for _, tmp := range tdst.Indexes {
					if *csrc.Name == *tmp.Name {
						cdst = tmp
						break
					}
				}
				if cdst == nil {
					mscb.WriteString(action.Index.GenerateNew(csrc, nil))
				} else if !c.IndexEq(csrc, cdst) {
					mscb.WriteString(action.Index.GenerateDel(cdst, nil))
					mscb.WriteString(action.Index.GenerateNew(csrc, nil))
				}
			}
			for _, cdst := range tdst.Indexes {
				var csrc *pb2.Index
				for _, tmpb := range tsrc.Indexes {
					if *cdst.Name == *tmpb.Name {
						csrc = tmpb
						break
					}
				}
				if csrc == nil {
					mscb.WriteString(action.Index.GenerateDel(cdst, nil))
				}
			}
		}
	}
	for _, tdst := range dbdst.Tables {
		var tsrc *pb2.Table
		for _, tmpb := range dbsrc.Tables {
			if *tdst.Name == *tmpb.Name {
				tsrc = tmpb
				break
			}
		}
		if tsrc == nil {
			tabb.WriteString(action.Table.GenerateDel(tdst, nil))
		}
	}
}

func generateStoreProcedureDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, spcb *bytes.Buffer) {
	//procVisited := make(map[string]bool)
	for _, psrc := range dbsrc.StoredProcs {
		var pdst *pb2.StoredProcedure
		for _, tmpb := range dbdst.StoredProcs {
			if *tmpb.Name == *psrc.Name {
				pdst = tmpb
				//ok1 := procVisited[*psrc.Name]
				//ok2 := procVisited[*pdst.Name]
				if *dbsrc.Type == pb2.Db_ORACLE || (*dbsrc.Type != pb2.Db_ORACLE && *psrc.NumParams == *pdst.NumParams) {
					//procVisited[*psrc.Name] = true
					//procVisited[*pdst.Name] = true
					break
				}
				pdst = nil
			} else {
				pdst = nil
			}
		}
		if pdst == nil {
			spcb.WriteString(action.StoredProcedure.GenerateNew(psrc, nil))
		} else if !c.StoredProcedureEq(psrc, pdst) {
			spcb.WriteString(action.StoredProcedure.GenerateDel(pdst, nil))
			spcb.WriteString(action.StoredProcedure.GenerateNew(psrc, nil))
		}
	}
	for _, pdst := range dbdst.StoredProcs {
		var psrc *pb2.StoredProcedure
		for _, tmpb := range dbsrc.StoredProcs {
			if *tmpb.Name == *pdst.Name {
				psrc = tmpb
				break
			}
		}
		if psrc == nil {
			spcb.WriteString(action.StoredProcedure.GenerateDel(pdst, nil))
		}
	}
}

func generateSequenceDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, seqb *bytes.Buffer) {
	for _, ssrc := range dbsrc.Sequences {
		var sdst *pb2.Sequence
		for _, tmpb := range dbdst.Sequences {
			if *ssrc.Name == *tmpb.Name {
				sdst = tmpb
				if c.SequenceEq(ssrc, sdst) {
					break
				}
			}
		}
		if sdst == nil {
			seqb.WriteString(action.Sequence.GenerateNew(ssrc, nil))
		} else if !c.SequenceEq(ssrc, sdst) {
			seqb.WriteString(action.Sequence.GenerateDel(sdst, nil))
			seqb.WriteString(action.Sequence.GenerateNew(ssrc, nil))
		}
	}
	for _, tdst := range dbdst.Sequences {
		var tsrc *pb2.Sequence
		for _, tmpb := range dbsrc.Sequences {
			if *tdst.Name == *tmpb.Name {
				tsrc = tmpb
				break
			}
		}
		if tsrc == nil {
			seqb.WriteString(action.Sequence.GenerateDel(tdst, nil))
		}
	}
}

func generateTriggerDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, trgb *bytes.Buffer) {
	for _, tsrc := range dbsrc.Tables {
		var tdst *pb2.Table
		for _, tmpb := range dbdst.Tables {
			if *tsrc.Name == *tmpb.Name {
				tdst = tmpb
				break
			}
		}
		if tdst == nil {
			for _, csrc := range tsrc.Triggers {
				trgb.WriteString(action.Trigger.GenerateNew(csrc, nil))
			}
		} else {
			for _, csrc := range tsrc.Triggers {
				var cdst *pb2.Trigger
				for _, tmp := range tdst.Triggers {
					if *csrc.Name == *tmp.Name {
						cdst = tmp
						if c.TriggerEq(csrc, cdst) {
							break
						}
					}
				}
				if cdst == nil {
					trgb.WriteString(action.Trigger.GenerateNew(csrc, nil))
				} else if !c.TriggerEq(csrc, cdst) {
					trgb.WriteString(action.Trigger.GenerateDel(cdst, nil))
					trgb.WriteString(action.Trigger.GenerateNew(csrc, nil))
				}
			}
			for _, cdst := range tdst.Triggers {
				var csrc *pb2.Trigger
				for _, tmpb := range tsrc.Triggers {
					if *tmpb.Name == *cdst.Name {
						csrc = tmpb
						break
					}
				}
				if csrc == nil {
					trgb.WriteString(action.Trigger.GenerateDel(cdst, nil))
				}
			}
		}
	}
}

func generateViewDiff(dbsrc *pb2.Db, dbdst *pb2.Db, action *c.SchemaDiffAction, tabb *bytes.Buffer) {
	var delViews []*pb2.View
	var addViews []*pb2.View
	for _, vsrc := range dbsrc.Views {
		var vdst *pb2.View
		for _, tmpb := range dbdst.Views {
			if *vsrc.Name == *tmpb.Name {
				vdst = tmpb
				break
			}
		}
		if vdst == nil {
			addViews = append(addViews, vsrc)
		} else if !c.ViewEq(vsrc, vdst) {
			delViews = append(delViews, vdst)
			addViews = append(addViews, vsrc)
		}
	}
	for _, vwc := range delViews {
		for _, vw := range delViews {
			vwpt := regexp.MustCompile("[\t ,]+" + *vw.Name + "[\t ,]+(?i)")
			if vwpt.MatchString(*vwc.Definition) {
				*vw.Weight = *vw.Weight + 1
			}
		}
	}
	sort.Slice(delViews, func(i, j int) bool {
		return *delViews[i].Weight > *delViews[j].Weight
	})
	for l, r := 0, len(delViews)-1; l < r; l, r = l+1, r-1 {
		delViews[l], delViews[r] = delViews[r], delViews[l]
	}
	for _, vwc := range addViews {
		for _, vw := range addViews {
			vwpt := regexp.MustCompile("[\t ,]+" + *vw.Name + "[\t ,]+(?i)")
			if vwpt.MatchString(*vwc.Definition) {
				*vw.Weight = *vw.Weight + 1
			}
		}
	}
	sort.Slice(addViews, func(i, j int) bool {
		return *addViews[i].Weight > *addViews[j].Weight
	})
	for l, r := 0, len(addViews)-1; l < r; l, r = l+1, r-1 {
		addViews[l], addViews[r] = addViews[r], addViews[l]
	}
	for _, vwc := range delViews {
		tabb.WriteString(action.View.GenerateDel(vwc, nil))
	}
	for _, vwc := range addViews {
		tabb.WriteString(action.View.GenerateNew(vwc, nil))
	}
}
