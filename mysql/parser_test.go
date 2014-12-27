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

	testStatement(t, "alter table `hoge` DROP KEY `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropIndex{IndexNameIdentifier{Name: "fuga"}}} })
	testStatement(t, "alter table `hoge` DROP INDEX `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropIndex{IndexNameIdentifier{Name: "fuga"}}} })

	testStatement(t, "alter table `hoge` ADD COLUMN `fuga` INT", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationAddColumn{ColumnNameIdentifier{Name: "fuga"}, ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 0, false, false}}}}})
}

func TestParseColumnDefinition(t *testing.T) {
	testColumnDefinition(t, "INT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_INT, 0, false, false }})
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
		t.Errorf("Test failed abount \"%s\":\n\tExpect\t: %+#v, \n\tBut Got\t: %+#v", src, expect, statements[0])
		return
	}
}

func testColumnDefinition(t *testing.T, src string, expect interface{}) {
	s := new(Scanner)
	s.Init("ALTER TABLE hoge ADD COLUMN fuga " + src + ";")
	statements := Parse(s)
	if len(statements) != 1 {
		t.Errorf("Expect %q to be parsed, but %+#v", src, statements)
		return
	}
	if v, ok := statements[0].(*AlterTableStatement); ok {
		if len(v.AlterSpecifications) == 1 {
			if v, ok := v.AlterSpecifications[0].(*AlterSpecificationAddColumn); ok {
				if !reflect.DeepEqual(v.ColumnDefinition, expect) {
					t.Errorf("Test failed abount \"%s\":\n\tExpect\t: %+#v, \n\tBut Got\t: %+#v", src, expect, v.ColumnDefinition)
				}
			} else {
				t.Errorf("Expect %q to be parsed, but %+#v", src, v)
			}
		} else {
			t.Errorf("Expect %q to be parsed, but %+#v", src, v)
		}
	} else {
		t.Errorf("statement should be AlterTableStatement\n")
	}
}
