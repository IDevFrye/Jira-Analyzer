package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	handlerErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	configreader "github.com/jiraconnector/internal/configReader"
	myErr "github.com/jiraconnector/internal/connector/errors"
	"github.com/jiraconnector/internal/structures"
)

type JiraConnector struct {
	cfg    *configreader.JiraConfig
	client *http.Client
}

type JiraConnectorInterface interface {
	GetAllProjects() ([]structures.JiraProject, error)
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	GetProjectIssues(project string) ([]structures.JiraIssue, error)
}

func NewJiraConnector(config configreader.Config) *JiraConnector {
	return &JiraConnector{
		cfg:    &config.JiraCfg,
		client: &http.Client{},
	}
}

func (con *JiraConnector) GetAllProjects() ([]structures.JiraProject, error) {
	url := fmt.Sprintf("%s/rest/api/2/project", con.cfg.Url)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrReadResponseBody, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	var projects []structures.JiraProject
	if err = json.Unmarshal(body, &projects); err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrUnmarshalAns, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	return projects, nil
}

func (con *JiraConnector) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	allProjects, err := con.GetAllProjects()
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrGetProjects, err)
		log.Println(ansErr)
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
		return &structures.ResponseProject{}, nil
	}
	end := start + limit
	if end > totalProjects {
		end = totalProjects
	}

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
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrGetIssues, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	if totalIssues == 0 {
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
				log.Println("stop go thread")
				return
			default:
				issues, err := con.getIssuesForOneThread(issueStart, project)
				if err != nil {
					ansErr := fmt.Errorf("%w :: %w", myErr.ErrGetIssues, err)
					errChan <- ansErr
					log.Println(ansErr)
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
			log.Println(err)
			cancel()
		}
	}()

	wg.Wait()

	log.Println("Got all issues")
	return allIssues, nil
}

func (con *JiraConnector) getIssuesForOneThread(startAt int, project string) ([]structures.JiraIssue, error) {
	url := fmt.Sprintf(
		"%s/rest/api/2/search?jql=project=%s&expand=changelog&startAt=%d&maxResult=%d",
		con.cfg.Url, project, startAt, con.cfg.IssueInOneReq)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrReadResponseBody, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrUnmarshalAns, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	return issues.Issues, nil
}

func (con *JiraConnector) getTotalIssues(project string) (int, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=0&", con.cfg.Url, project)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrReadResponseBody, err)
		log.Println(ansErr)
		return 0, ansErr
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		ansErr := fmt.Errorf("%w :: %w", myErr.ErrUnmarshalAns, err)
		log.Println(ansErr)
		return 0, ansErr
	}

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
		return nil, fmt.Errorf("%w :: %w", myErr.ErrMakeRequest, err)
	}

	for {
		resp, err = con.client.Do(req)

		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
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
