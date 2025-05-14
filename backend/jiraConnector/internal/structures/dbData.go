package structures

import "time"

type DBStatusChanges struct {
	IssueId    int
	AuthorId   int
	ChangeTime time.Time
	FromStatus string
	ToStatus   string
}

type DBAuthor struct {
	Id   int
	Name string
}

type DBProject struct {
	Id    int
	Title string
}

type DBIssue struct {
	Id          int
	ProjectId   int
	AuthorId    int
	AssigneeId  int
	Key         string
	Summary     string
	Description string
	Type        string
	Priority    string
	Status      string
	CreatedTime time.Time
	ClosedTime  time.Time
	UpdatedTime time.Time
	TimeSpent   int
}
