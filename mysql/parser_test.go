package mysql

import (
	"reflect"
	"testing"
)

func TestParseDropTableStatement(t *testing.T) {
	testStatement(t, "DROP TABLE hoge", &DropTableStatement{TableNames:[]TableNameIdentifier{TableNameIdentifier{Name:"hoge"}}})
	testStatement(t, "drop table hoge,fuga", &DropTableStatement{TableNames:[]TableNameIdentifier{TableNameIdentifier{Name:"fuga"}, TableNameIdentifier{Name:"hoge"}}})
	testStatement(t, "drop table `hoge`", &DropTableStatement{TableNames:[]TableNameIdentifier{TableNameIdentifier{Name:"hoge"}}})
	testStatement(t, "drop table hoge.fuga", &DropTableStatement{TableNames:[]TableNameIdentifier{TableNameIdentifier{Database: "hoge", Name:"fuga"}}})
}

func TestParseDropDatabaseStatement(t *testing.T) {
	testStatement(t, "DROP DATABASE hoge", &DropDatabaseStatement{DatabaseNameIdentifier{Name:"hoge"}})
	testStatement(t, "drop database `hoge`", &DropDatabaseStatement{DatabaseNameIdentifier{Name:"hoge"}})
}

func TestParseCreateDatabaseStatement(t *testing.T) {
	testStatement(t, "CREATE DATABASE hoge", &CreateDatabaseStatement{DatabaseNameIdentifier{Name:"hoge"}})
	testStatement(t, "create database `hoge`", &CreateDatabaseStatement{DatabaseNameIdentifier{Name:"hoge"}})
}

func TestParseAlterTableStatement(t *testing.T) {
	testStatement(t, "ALTER TABLE hoge", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, nil })
	testStatement(t, "alter table `hoge`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, nil })
	testStatement(t, "alter table `hoge` DROP COLUMN fuga", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropColumn{ColumnNameIdentifier{Name: "fuga"}}}})
	testStatement(t, "alter table `hoge` DROP `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropColumn{ColumnNameIdentifier{Name: "fuga"}}} })
}


func testStatement(t *testing.T, src string, expect interface{}) {
	s := new(Scanner)
	s.Init(src + ";")
	statements := Parse(s)
	if len(statements) != 1 {
		t.Errorf("Expect %q to be parsed, but %+#v", src, statements)
		return
	}
	if !reflect.DeepEqual(statements[0], expect) {
		t.Errorf("\tExpect\t: %+#v, \n\t\tBut Got\t: %+#v", expect, statements[0])
		return
	}
}

