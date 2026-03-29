package integration

import "os"

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

func getBackendBaseURL() string {
	if url := os.Getenv("BACKEND_BASE_URL"); url != "" {
		return url
	}
	return "http://backend:8000/api/v1"
}

func findProjectByName(projects []JiraProject, name string) (JiraProject, bool) {
	for _, p := range projects {
		if p.Name == name {
			return p, true
		}
	}
	return JiraProject{}, false
}

func findProjectByKey(projects []JiraProject, key string) (JiraProject, bool) {
	for _, p := range projects {
		if p.Key == key {
			return p, true
		}
	}
	return JiraProject{}, false
}

func filterProjectsByNames(projects []JiraProject, names ...string) []JiraProject {
	nameSet := make(map[string]struct{}, len(names))
	for _, n := range names {
		nameSet[n] = struct{}{}
	}

	out := make([]JiraProject, 0, len(names))
	for _, p := range projects {
		if _, ok := nameSet[p.Name]; ok {
			out = append(out, p)
		}
	}
	return out
}

func filterProjectsByKeys(projects []JiraProject, keys ...string) []JiraProject {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	out := make([]JiraProject, 0, len(keys))
	for _, p := range projects {
		if _, ok := keySet[p.Key]; ok {
			out = append(out, p)
		}
	}
	return out
}

func pickFirstSyncedProject(dbProjects []JiraProject, candidates ...JiraProject) (JiraProject, bool) {
	for _, c := range candidates {
		for _, p := range dbProjects {
			if p.Key == c.Key || p.Name == c.Name {
				return p, true
			}
		}
	}
	return JiraProject{}, false
}

func pickSyncedProjects(dbProjects []JiraProject, candidates ...JiraProject) []JiraProject {
	out := make([]JiraProject, 0, len(candidates))
	for _, c := range candidates {
		for _, p := range dbProjects {
			if p.Key == c.Key || p.Name == c.Name {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
