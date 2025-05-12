package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	handlerErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	myErr "github.com/jiraconnector/internal/connector/errors"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
)

type JiraConnector struct {
	cfg    *config.JiraConfig
	client *http.Client
	log    *slog.Logger
}

func NewJiraConnector(config *config.Config, log *slog.Logger) *JiraConnector {
	return &JiraConnector{
		cfg:    &config.JiraCfg,
		client: &http.Client{},
		log:    log,
	}
}

func (con *JiraConnector) GetAllProjects() ([]structures.JiraProject, error) {
	url := fmt.Sprintf("%s/rest/api/2/project", con.cfg.Url)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		con.log.Error("err retry request", logger.Err(err), "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrReadResponseBody, err)
		con.log.Error(ansErr.Error(), "url", url)
		return nil, ansErr
	}

	var projects []structures.JiraProject
	if err = json.Unmarshal(body, &projects); err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrUnmarshalAns, err)
		con.log.Error(ansErr.Error(), "url", url)
		return nil, ansErr
	}

	con.log.Info("success get all projects from", "url", url)
	return projects, nil
}

func (con *JiraConnector) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	allProjects, err := con.GetAllProjects()
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrGetProjects, err)
		con.log.Error(ansErr.Error(), "search", search, "page", page, "limit", limit)
		return nil, ansErr
	}

	var pageProjects []structures.JiraProject
	for _, proj := range allProjects {
		if search == "" || containsSearchProject(proj.Name, search) {
			pageProjects = append(pageProjects, proj)
		}
	}

	totalProjects := len(pageProjects)
	start := (page - 1) * limit
	if start >= totalProjects {
		con.log.Info("success get projects with params", "search", search, "page", page, "limit", limit)
		return &structures.ResponseProject{}, nil
	}
	end := start + limit
	if end > totalProjects {
		end = totalProjects
	}

	con.log.Info("success get projects with params", "search", search, "page", page, "limit", limit)

	return &structures.ResponseProject{
			Projects: pageProjects[start:end],
			PageInfo: structures.PageInfo{
				PageCount:     int(math.Ceil(float64(totalProjects) / float64(limit))),
				CurrentPage:   page,
				ProjectsCount: totalProjects,
			},
		},
		nil
}

func (con *JiraConnector) GetProjectIssues(project string) ([]structures.JiraIssue, error) {
	//get all issues for this project
	totalIssues, err := con.getTotalIssues(project)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrGetIssues, err)
		con.log.Error(ansErr.Error(), "project", project)
		return nil, ansErr
	}

	if totalIssues == 0 {
		con.log.Info("success got all issues", "project", project)
		return []structures.JiraIssue{}, nil
	}

	//create common source = map for results
	var allIssues []structures.JiraIssue
	threadCount := con.cfg.ThreadCount
	issueReq := con.cfg.IssueInOneReq

	//create go routines
	var wg sync.WaitGroup
	var issuesMux sync.Mutex

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)

	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		//find start index for get issues
		issueStart := i*issueReq + 1
		if issueStart > totalIssues {
			issueStart = totalIssues
		}

		go func() {

			defer wg.Done()
			//TODO: add count of request (?)
			select {
			case <-ctx.Done():
				con.log.Error("stop go thread", "project", project)
				return
			default:
				issues, err := con.getIssuesForOneThread(issueStart, project)
				if err != nil {
					ansErr := fmt.Errorf("%w: %w", myErr.ErrGetIssues, err)
					errChan <- ansErr
					con.log.Error(ansErr.Error(), "project", project)
					return
				}

				issuesMux.Lock()
				defer issuesMux.Unlock()
				allIssues = append(allIssues, issues...)

			}
		}()
	}

	go func() {
		if err := <-errChan; err != nil {
			con.log.Error("err chan", logger.Err(err), "project", project)
			cancel()
		}
	}()

	wg.Wait()

	con.log.Info("success got all issues", "project", project)
	return allIssues, nil
}

func (con *JiraConnector) getIssuesForOneThread(startAt int, project string) ([]structures.JiraIssue, error) {
	url := fmt.Sprintf(
		"%s/rest/api/2/search?jql=project=%s&expand=changelog&startAt=%d&maxResult=%d",
		con.cfg.Url, project, startAt, con.cfg.IssueInOneReq)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		con.log.Error("error retry request", logger.Err(err), "project", project, "startAt", startAt)
		return nil, myErr.ErrGetIssues
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrReadResponseBody, err)
		con.log.Error(ansErr.Error(), "project", project, "startAt", startAt)
		return nil, ansErr
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrUnmarshalAns, err)
		con.log.Error(ansErr.Error(), "project", project, "startAt", startAt)
		return nil, ansErr
	}

	con.log.Info("success get issue for thread", "project", project, "startAt", startAt)
	return issues.Issues, nil
}

func (con *JiraConnector) getTotalIssues(project string) (int, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=0&", con.cfg.Url, project)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		con.log.Error("error retry request", logger.Err(err), "project", project)
		return 0, myErr.ErrGetIssues
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrReadResponseBody, err)
		con.log.Error(ansErr.Error(), "project", project)
		return 0, ansErr
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrUnmarshalAns, err)
		con.log.Error(ansErr.Error(), "project", project)
		return 0, ansErr
	}

	con.log.Info("success got all issues", "project", project)
	return issues.Total, nil
}

func (con *JiraConnector) retryRequest(method, url string) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	timeSleep := con.cfg.MinSleep

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myErr.ErrMakeRequest, err)
		con.log.Error(ansErr.Error(), "method", method, "url", url)
		return nil, ansErr
	}

	for {
		resp, err = con.client.Do(req)

		if resp == nil || resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
			ansErr := fmt.Errorf("%w", handlerErr.ErrNoProject)
			con.log.Error(ansErr.Error(), "method", method, "url", url)
			return nil, handlerErr.ErrNoProject
		}

		// if everything ok - return resp
		if err == nil && resp.StatusCode < 300 {
			return resp, nil
		}
		time.Sleep(time.Duration(timeSleep))
		timeSleep *= 2

		if timeSleep > con.cfg.MaxSleep {
			break
		}
	}

	// if in cycle we didn't do response - return err
	return nil, myErr.ErrMaxTimeRequest

}

func containsSearchProject(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
