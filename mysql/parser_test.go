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
	testColumnDefinition(t, "BIT", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_BIT }})
	testColumnDefinition(t, "TINYINT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_TINYINT, 0, false, false }})
	testColumnDefinition(t, "SMALLINT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_SMALLINT, 0, false, false }})
	testColumnDefinition(t, "MEDIUMINT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_MEDIUMINT, 0, false, false }})
	testColumnDefinition(t, "INT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_INT, 0, false, false }})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_INT, 10, true, true }})
	testColumnDefinition(t, "INTEGER", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_INT, 0, false, false }})
	testColumnDefinition(t, "BIGINT", ColumnDefinition{&DataTypeDefinitionNumber{ DATATYPE_BIGINT, 0, false, false }})
	testColumnDefinition(t, "REAL", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_REAL, 0, 0, false, false }})
	testColumnDefinition(t, "DOUBLE", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_DOUBLE, 0, 0, false, false }})
	testColumnDefinition(t, "FLOAT", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_FLOAT, 0, 0, false, false }})
	testColumnDefinition(t, "FLOAT(10, 2) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_FLOAT, 10, 2, true, true}})
	testColumnDefinition(t, "DECIMAL", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_DECIMAL, 0, 0, false, false }})
	testColumnDefinition(t, "DECIMAL(10, 2) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_DECIMAL, 10, 2, true, true }})
	testColumnDefinition(t, "DECIMAL(10) ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_DECIMAL, 10, 0, false, true }})
	testColumnDefinition(t, "NUMERIC", ColumnDefinition{&DataTypeDefinitionFraction{ DATATYPE_NUMERIC, 0, 0, false, false }})
	testColumnDefinition(t, "DATE", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_DATE }})
	testColumnDefinition(t, "TIME", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_TIME }})
	testColumnDefinition(t, "TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_TIMESTAMP }})
	testColumnDefinition(t, "DATETIME", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_DATETIME }})
	testColumnDefinition(t, "YEAR", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_YEAR }})
	testColumnDefinition(t, "CHAR", ColumnDefinition{&DataTypeDefinitionString{ DATATYPE_CHAR, 0, "", "" }})
	testColumnDefinition(t, "CHAR(255)", ColumnDefinition{&DataTypeDefinitionString{ DATATYPE_CHAR, 255, "", "" }})
	testColumnDefinition(t, "VARCHAR", ColumnDefinition{&DataTypeDefinitionString{ DATATYPE_VARCHAR, 0, "", "" }})
	testColumnDefinition(t, "VARCHAR(255)", ColumnDefinition{&DataTypeDefinitionString{ DATATYPE_VARCHAR, 255, "", "" }})
	testColumnDefinition(t, "BINARY", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_BINARY }})
	testColumnDefinition(t, "VARBINARY", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_VARBINARY }})
	testColumnDefinition(t, "TINYBLOB", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_TINYBLOB }})
	testColumnDefinition(t, "BLOB", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_BLOB }})
	testColumnDefinition(t, "MEDIUMBLOB", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_MEDIUMBLOB }})
	testColumnDefinition(t, "LONGBLOB", ColumnDefinition{&DataTypeDefinitionSimple{ DATATYPE_LONGBLOB }})
	testColumnDefinition(t, "TINYTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{ DATATYPE_TINYTEXT, false, "", ""}})
	testColumnDefinition(t, "TEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{ DATATYPE_TEXT, false, "", ""}})
	testColumnDefinition(t, "MEDIUMTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{ DATATYPE_MEDIUMTEXT, false, "", ""}})
	testColumnDefinition(t, "LONGTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{ DATATYPE_LONGTEXT, false, "", ""}})
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
		t.Errorf("Test failed about \"%s\":\n\tExpect\t: %+#v, \n\tBut Got\t: %+#v", src, expect, statements[0])
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
					t.Errorf("Test failed about \"%s\":\n\tExpect\t: %+#v, \n\tBut Got\t: %+#v", src, expect, v.ColumnDefinition)
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
