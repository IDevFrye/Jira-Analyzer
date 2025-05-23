package errors

import "errors"

var (
	ErrOpenDb = errors.New("can't open database")

	ErrInsertProject = errors.New("can't insert project")
	ErrSelectProject = errors.New("can't select project")
	ErrPushProject   = errors.New("can't while push project")

	ErrInsertAuthor = errors.New("can't insert author")
	ErrSelectAuthor = errors.New("can't select author")

	ErrInsertIssue = errors.New("can't insert Issue")
	ErrPushIssue   = errors.New("can't push Issue")

	ErrInsertStatusChange = errors.New("can't insert status change")

	ErrTranBegin = errors.New("error transaction begin")
	ErrTranClose = errors.New("error transaction close")
)
