package mysql

import (
	"reflect"
	"testing"
)

func TestParseDropTableStatement(t *testing.T) {
	testStatement(t, "DROP TABLE hoge", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "hoge"}}})
	testStatement(t, "drop table hoge,fuga", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "fuga"}, TableNameIdentifier{Name: "hoge"}}})
	testStatement(t, "drop table `TABLE`", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "TABLE"}}})
	testStatement(t, "drop table hoge.fuga", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Database: "hoge", Name: "fuga"}}})
}

func TestParseDropDatabaseStatement(t *testing.T) {
	testStatement(t, "DROP DATABASE hoge", &DropDatabaseStatement{DatabaseNameIdentifier{Name: "hoge"}})
	testStatement(t, "drop database `hoge`", &DropDatabaseStatement{DatabaseNameIdentifier{Name: "hoge"}})
}

func TestParseCreateDatabaseStatement(t *testing.T) {
	testStatement(t, "CREATE DATABASE hoge", &CreateDatabaseStatement{DatabaseNameIdentifier{Name: "hoge"}})
	testStatement(t, "create database `hoge`", &CreateDatabaseStatement{DatabaseNameIdentifier{Name: "hoge"}})
}

func TestCreateTableStatement(t *testing.T) {
	testStatement(t, "CREATE TABLE hoge ( id INT(10) UNSIGNED NOT NULL, PRIMARY KEY (id) ) ENGINE=InnoDB", &CreateTableStatement{TableNameIdentifier{"hoge", ""}, []CreateDefinition{
		&CreateDefinitionColumn{ColumnNameIdentifier{"id"}, ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, false}, false, false, &DefaultDefinitionEmpty{}}},
		&CreateDefinitionPrimaryIndex{[]ColumnNameIdentifier{ColumnNameIdentifier{"id"}}},
	}})
	testStatement(t, "CREATE TABLE hoge ( id INT(10) UNSIGNED NOT NULL, name VARCHAR(255) NOT NULL, PRIMARY KEY (id, name), UNIQUE INDEX name (name), INDEX (id) )", &CreateTableStatement{TableNameIdentifier{"hoge", ""}, []CreateDefinition{
		&CreateDefinitionColumn{ColumnNameIdentifier{"id"}, ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, false}, false, false, &DefaultDefinitionEmpty{}}},
		&CreateDefinitionColumn{ColumnNameIdentifier{"name"}, ColumnDefinition{&DataTypeDefinitionString{DATATYPE_VARCHAR, 255, "", ""}, false, false, &DefaultDefinitionEmpty{}}},
		&CreateDefinitionPrimaryIndex{[]ColumnNameIdentifier{ColumnNameIdentifier{"id"}, ColumnNameIdentifier{"name"}}},
		&CreateDefinitionUniqueIndex{IndexNameIdentifier{"name"}, []ColumnNameIdentifier{ColumnNameIdentifier{"name"}}},
		&CreateDefinitionIndex{IndexNameIdentifier{""}, []ColumnNameIdentifier{ColumnNameIdentifier{"id"}}},
	}})
}

func TestParseAlterTableStatement(t *testing.T) {
	testStatement(t, "ALTER TABLE hoge", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, nil})
	testStatement(t, "alter table `hoge`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, nil})

	testStatement(t, "alter table `hoge` DROP COLUMN fuga", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropColumn{ColumnNameIdentifier{Name: "fuga"}}}})
	testStatement(t, "alter table `hoge` DROP `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropColumn{ColumnNameIdentifier{Name: "fuga"}}}})

	testStatement(t, "alter table `hoge` DROP KEY `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropIndex{IndexNameIdentifier{Name: "fuga"}}}})
	testStatement(t, "alter table `hoge` DROP INDEX `fuga`", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationDropIndex{IndexNameIdentifier{Name: "fuga"}}}})

	testStatement(t, "alter table `hoge` ADD COLUMN `fuga` INT", &AlterTableStatement{TableNameIdentifier{Name: "hoge"}, []AlterSpecification{&AlterSpecificationAddColumn{ColumnNameIdentifier{Name: "fuga"}, ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}}}}})
}

func TestParseCommentStatement(t *testing.T) {
	testStatement(t, "/* hoge */", &CommentStatement{" hoge "})
	testStatement(t, "/* あいうえお */", &CommentStatement{" あいうえお "})
	testStatement(t, "/* SELECT * FROM hoge; */", &CommentStatement{" SELECT * FROM hoge; "})
}

func TestParseColumnDefinition(t *testing.T) {
	testColumnDefinition(t, "BIT", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_BIT}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "bit", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_BIT}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TINYINT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_TINYINT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "SMALLINT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_SMALLINT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "MEDIUMINT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_MEDIUMINT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "INT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL NOT NULL AUTO_INCREMENT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, false, true, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL NOT NULL DEFAULT 100 AUTO_INCREMENT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, false, true, &DefaultDefinitionString{"100"}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL NOT NULL DEFAULT '100' AUTO_INCREMENT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, false, true, &DefaultDefinitionString{"100"}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL NOT NULL DEFAULT \"100\" AUTO_INCREMENT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, false, true, &DefaultDefinitionString{"100"}})
	testColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, true, false, &DefaultDefinitionNull{}})
	testColumnDefinition(t, "INTEGER", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "BIGINT", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_BIGINT, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "REAL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_REAL, 0, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "DOUBLE", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DOUBLE, 0, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "FLOAT", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_FLOAT, 0, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "FLOAT(10, 2) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_FLOAT, 10, 2, true, true}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "DECIMAL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 0, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "DECIMAL(10, 2) UNSIGNED ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 10, 2, true, true}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "DECIMAL(10) ZEROFILL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 10, 0, false, true}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "NUMERIC", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_NUMERIC, 0, 0, false, false}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "DATE", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATE}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TIME", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_TIME}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_TIMESTAMP}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TIMESTAMP DEFAULT CURRENT_TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_TIMESTAMP}, true, false, &DefaultDefinitionCurrentTimestamp{false}})
	testColumnDefinition(t, "TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_TIMESTAMP}, true, false, &DefaultDefinitionCurrentTimestamp{true}})
	testColumnDefinition(t, "DATETIME", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATETIME}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "YEAR", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_YEAR}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "CHAR", ColumnDefinition{&DataTypeDefinitionString{DATATYPE_CHAR, 0, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "CHAR(255)", ColumnDefinition{&DataTypeDefinitionString{DATATYPE_CHAR, 255, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "VARCHAR", ColumnDefinition{&DataTypeDefinitionString{DATATYPE_VARCHAR, 0, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "VARCHAR(255)", ColumnDefinition{&DataTypeDefinitionString{DATATYPE_VARCHAR, 255, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "BINARY", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_BINARY}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "VARBINARY", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_VARBINARY}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TINYBLOB", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_TINYBLOB}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "BLOB", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_BLOB}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "MEDIUMBLOB", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_MEDIUMBLOB}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "LONGBLOB", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_LONGBLOB}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TINYTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{DATATYPE_TINYTEXT, false, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "TEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{DATATYPE_TEXT, false, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "MEDIUMTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{DATATYPE_MEDIUMTEXT, false, "", ""}, true, false, &DefaultDefinitionEmpty{}})
	testColumnDefinition(t, "LONGTEXT", ColumnDefinition{&DataTypeDefinitionTextBlob{DATATYPE_LONGTEXT, false, "", ""}, true, false, &DefaultDefinitionEmpty{}})
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
