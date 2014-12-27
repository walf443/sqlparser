package mysql

type (
	Statement interface {
		statement()
	}

	Identifier interface {
		identifier()
	}

	AlterSpecification interface {
		alterspecification()
	}
)

type (
	DropTableStatement struct {
		TableNames []TableNameIdentifier
	}
	DropDatabaseStatement struct {
		DatabaseName DatabaseNameIdentifier
	}
	CreateDatabaseStatement struct {
		DatabaseName DatabaseNameIdentifier
	}
	AlterTableStatement struct {
		TableName TableNameIdentifier
		AlterSpecifications []AlterSpecification
	}
)

func (x *DropTableStatement) statement() {}
func (x *DropDatabaseStatement) statement() {}
func (x *CreateDatabaseStatement) statement() {}
func (x *AlterTableStatement) statement() {}

type (
	TableNameIdentifier struct {
		Name string
		Database string
	}
	DatabaseNameIdentifier struct {
		Name string
	}
	ColumnNameIdentifier struct {
		Name string
	}
	IndexNameIdentifier struct {
		Name string
	}
)

func (x *TableNameIdentifier) identifier() {}
func (x *DatabaseNameIdentifier) identifier() {}
func (x *ColumnNameIdentifier) identifier() {}
func (x *IndexNameIdentifier) identifier() {}

type (
	AlterSpecificationDropColumn struct {
		ColumnName ColumnNameIdentifier
	}
	AlterSpecificationDropIndex struct {
		IndexName IndexNameIdentifier
	}
)

func (x *AlterSpecificationDropColumn) alterspecification() {}
func (x *AlterSpecificationDropIndex) alterspecification() {}

