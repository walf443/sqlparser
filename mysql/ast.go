package mysql

type (
	Statement interface {
		statement()
	}

	Identifier interface {
		identifier()
	}
)

type (
	DropTableStatement struct {
		TableNames []TableNameIdentifier
	}
	DropDatabaseStatement struct {
		DatabaseName DatabaseNameIdentifier
	}
)

func (x *DropTableStatement) statement() {}
func (x *DropDatabaseStatement) statement() {}

type (
	TableNameIdentifier struct {
		Name string
		Database string
	}
	DatabaseNameIdentifier struct {
		Name string
	}
)

func (x *TableNameIdentifier) identifier() {}
func (x *DatabaseNameIdentifier) identifier() {}
