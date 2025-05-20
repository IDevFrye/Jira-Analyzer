package integration

type JiraProject struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Url  string `json:"self"`
}

type ResponseProject struct {
	Projects []JiraProject `json:"projects"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type PageInfo struct {
	PageCount     int `json:"pageCount"`
	CurrentPage   int `json:"currentPage"`
	ProjectsCount int `json:"projectsCount"`
}

type ResponseUpdate struct {
	Project string `json:"project"`
	Status  string `json:"status"`
}

type ProjectStats struct {
	TotalIssues        int     `json:"total_issues"`
	OpenIssues         int     `json:"open_issues"`
	ClosedIssues       int     `json:"closed_issues"`
	ReopenedIssues     int     `json:"reopened_issues"`
	ResolvedIssues     int     `json:"resolved_issues"`
	InProgressIssues   int     `json:"in_progress_issues"`
	AvgResolutionTimeH float64 `json:"avg_resolution_time_h"`
	AvgCreatedPerDay7d float64 `json:"avg_created_per_day_7d"`
}

type PriorityStats map[string]map[string]int
