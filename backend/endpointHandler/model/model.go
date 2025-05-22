package model

type Project struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Self string `json:"self"`
}

type DBProject struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Key   string `db:"key" json:"key"`
	Self  string `db:"url" json:"url"`
}

type UIProject struct {
	ID        int    `json:"id" db:"id"` // или string, если id::text в запросе
	Key       string `json:"key" db:"key"`
	Name      string `json:"name" db:"name"`
	Self      string `json:"self" db:"self"`
	Existence bool   `json:"existence" db:"existence"`
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

type PageInfo struct {
	PageCount     int `json:"pageCount"`
	CurrentPage   int `json:"currentPage"`
	ProjectsCount int `json:"projectsCount"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	PageInfo PageInfo  `json:"pageInfo"`
}
