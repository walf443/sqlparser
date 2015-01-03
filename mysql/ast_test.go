package mysql

import (
	"testing"
)

func TestGenDropTableStatement(t *testing.T) {
	testGenStatement(t, "DROP TABLE `hoge`", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "hoge"}}})
	testGenStatement(t, "DROP TABLE `fuga`, `hoge`", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "fuga"}, TableNameIdentifier{Name: "hoge"}}})
	testGenStatement(t, "DROP TABLE `TABLE`", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Name: "TABLE"}}})
	testGenStatement(t, "DROP TABLE `hoge`.`fuga`", &DropTableStatement{TableNames: []TableNameIdentifier{TableNameIdentifier{Database: "hoge", Name: "fuga"}}})
}

func TestGenDropDatabaseStatement(t *testing.T) {
	testGenStatement(t, "DROP DATABASE `hoge`", &DropDatabaseStatement{DatabaseName: DatabaseNameIdentifier{Name: "hoge"}})
}

func TestGenCreateDatabaseStatement(t *testing.T) {
	testGenStatement(t, "CREATE DATABASE `hoge`", &CreateDatabaseStatement{DatabaseName: DatabaseNameIdentifier{Name: "hoge"}})
}

func TestGenAlterStatement(t *testing.T) {
	testGenStatement(t, "ALTER TABLE `hoge` DROP `foo`", &AlterTableStatement{TableNameIdentifier{Name: "hoge", Database: ""}, []AlterSpecification{
		&AlterSpecificationDropColumn{ColumnNameIdentifier{"foo"}},
	}})
	testGenStatement(t, "ALTER TABLE `hoge` DROP INDEX `foo`", &AlterTableStatement{TableNameIdentifier{Name: "hoge", Database: ""}, []AlterSpecification{
		&AlterSpecificationDropIndex{IndexNameIdentifier{"foo"}},
	}})
}

func testGenStatement(t *testing.T, expected string, input Statement) {
	result := input.ToQuery()
	if result != expected {
		t.Errorf("Failed test to generage SQL\n\tInput: %+#v \n\tExpect\t: \"%s\"\n\tBut Got\t: \"%s\"", input, expected, result)
	}
}
