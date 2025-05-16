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
