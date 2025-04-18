package structures

type ResponseProject struct {
	Projects []JiraProject `json:"projects"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type PageInfo struct {
	PageCount     int `json:"pageCount"`
	CurrentPage   int `json:"currentPage"`
	ProjectsCount int `json:"projectsCount"`
}
