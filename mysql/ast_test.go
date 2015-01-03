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
	testGenStatement(t, "ALTER TABLE `hoge` ADD `foo` INT(10) UNSIGNED DEFAULT NULL", &AlterTableStatement{TableNameIdentifier{Name: "hoge", Database: ""}, []AlterSpecification{
		&AlterSpecificationAddColumn{ColumnNameIdentifier{"foo"}, ColumnDefinition{
			&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, false},
			true,
			false,
			&DefaultDefinitionNull{},
		}},
	}})
}

func TestGenColumnDefinition(t *testing.T) {
	testGenColumnDefinition(t, "INT DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 0, false, false}, true, false, &DefaultDefinitionNull{}})
	testGenColumnDefinition(t, "INT(10) UNSIGNED DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, false}, true, false, &DefaultDefinitionNull{}})
	testGenColumnDefinition(t, "INT(10) UNSIGNED ZEROFILL DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionNumber{DATATYPE_INT, 10, true, true}, true, false, &DefaultDefinitionNull{}})
	testGenColumnDefinition(t, "DATE ", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATE}, true, false, &DefaultDefinitionEmpty{}})
	testGenColumnDefinition(t, "DATE DEFAULT \"2015/01/04\"", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATE}, true, false, &DefaultDefinitionString{"2015/01/04"}})
	testGenColumnDefinition(t, "DATE DEFAULT CURRENT_TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATE}, true, false, &DefaultDefinitionCurrentTimestamp{}})
	testGenColumnDefinition(t, "DATE DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", ColumnDefinition{&DataTypeDefinitionSimple{DATATYPE_DATE}, true, false, &DefaultDefinitionCurrentTimestamp{true}})

	testGenColumnDefinition(t, "DECIMAL(10, 2) UNSIGNED ZEROFILL DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 10, 2, true, true}, true, false, &DefaultDefinitionNull{}})
	testGenColumnDefinition(t, "DECIMAL(10) UNSIGNED DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 10, 0, true, false}, true, false, &DefaultDefinitionNull{}})
	testGenColumnDefinition(t, "DECIMAL DEFAULT NULL", ColumnDefinition{&DataTypeDefinitionFraction{DATATYPE_DECIMAL, 0, 0, false, false}, true, false, &DefaultDefinitionNull{}})
}

func testGenStatement(t *testing.T, expected string, input Statement) {
	result := input.ToQuery()
	if result != expected {
		t.Errorf("Failed test to generage SQL\n\tInput: %+#v \n\tExpect\t: \"%s\"\n\tBut Got\t: \"%s\"", input, expected, result)
	}
}

func testGenColumnDefinition(t *testing.T, expected string, input ColumnDefinition) {
	specAddColumn := AlterSpecificationAddColumn{ColumnNameIdentifier{"foo"}, ColumnDefinition{}}
	specAddColumn.ColumnDefinition = input
	statement := AlterTableStatement{TableNameIdentifier{Name: "hoge", Database: ""}, []AlterSpecification{}}
	statement.AlterSpecifications = append(statement.AlterSpecifications, &specAddColumn)
	testGenStatement(t, "ALTER TABLE `hoge` ADD `foo` "+expected, &statement)
}
