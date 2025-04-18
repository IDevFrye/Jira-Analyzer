package structures

type JiraProject struct {
	// response: ".../project"
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Self string `json:"self"`
}

type JiraIssues struct {
	// response: ".../search?jql=project=idproject"
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []JiraIssue `json:"issues"`
}

type JiraIssue struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	Fields    Field     `json:"fields"`
	Changelog Changelog `json:"changelog"`
}

type Field struct {
	Project     JiraProject   `json:"project"`
	Author      User          `json:"creator"`
	Assignee    User          `json:"reporter"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Type        IssueType     `json:"issuetype"`
	Priority    IssuePriority `json:"priority"`
	Status      IssueStatus   `json:"status"`
	CreatedTime string        `json:"created"`
	ClosedTime  string        `json:"resolutiondate"`
	UpdatedTime string        `json:"updated"`
	TimeSpent   int           `json:"timespent"`
}

type User struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type IssueType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IssuePriority struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type IssueStatus struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Changelog struct {
	StartAt    int       `json:"startAt"`
	MaxResults int       `json:"maxResults"`
	Total      int       `json:"total"`
	Histories  []History `json:"histories"`
}

type History struct {
	Id      string `json:"id"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Items   []Item `json:"items"`
}

type Item struct {
	Field      string `json:"field"`
	Fieldtype  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}
