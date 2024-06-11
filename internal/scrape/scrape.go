package scrape

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/output"
)

const (
	siteAshby      string = "ashby"
	siteGreenhouse string = "greenhouse"
	siteLever      string = "lever"
)

func scrape(c db.Company) (jobs []output.Job, err error) {
	switch c.Site {
	case siteAshby:
		return ashby(c)
	case siteGreenhouse:
		return greenhouse(c)
	case siteLever:
		return lever(c)
	}
	return nil, nil
}

func ashby(company db.Company) (jobs []output.Job, err error) {
	query := strings.NewReader(fmt.Sprintf(`{"operationName":"ApiJobBoardWithTeams","variables":{"organizationHostedJobsPageName":"%s"},"query":"query ApiJobBoardWithTeams($organizationHostedJobsPageName: String!) {\n  jobBoard: jobBoardWithTeams(\n    organizationHostedJobsPageName: $organizationHostedJobsPageName\n  ) {\n    jobPostings {\n      title\n     }\n  }\n}"}`, company.Name))
	res, err := http.Post("https://jobs.ashbyhq.com/api/non-user-graphql?op=ApiJobBoardWithTeams", "application/json", query)

	body, err := checkThenDecode[struct {
		Data struct {
			JobBoard struct {
				JobPostings []struct {
					Title string `json:"title"`
				} `json:"jobPostings"`
			} `json:"jobBoard"`
		} `json:"data"`
	}](company, res, err)
	if err != nil {
		return nil, err
	}

	for _, j := range body.Data.JobBoard.JobPostings {
		jobs = append(jobs, output.Job{
			Company: company,
			Title:   j.Title,
		})
	}

	return jobs, nil
}

func greenhouse(company db.Company) (jobs []output.Job, err error) {
	res, err := http.Get(fmt.Sprintf("https://api.greenhouse.io/v1/boards/%s/jobs", company.Name))

	body, err := checkThenDecode[struct {
		Jobs []struct {
			Title     string `json:"title"`
			UpdatedAt string `json:"updated_at"`
		} `json:"jobs"`
	}](company, res, err)
	if err != nil {
		return nil, err
	}

	for _, resJob := range body.Jobs {
		j := output.Job{
			Company: company,
			Title:   resJob.Title,
		}

		parsedTime, err := time.Parse(time.RFC3339, resJob.UpdatedAt)
		if err == nil {
			j.UpdatedAt = parsedTime
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}

func lever(company db.Company) (jobs []output.Job, err error) {
	res, err := http.Get(fmt.Sprintf("https://api.lever.co/v0/postings/%s?limit=999", company.Name))

	body, err := checkThenDecode[[]struct {
		Title     string `json:"text"`
		UpdatedAt int    `json:"createdAt"`
	}](company, res, err)
	if err != nil {
		return nil, err
	}

	for _, j := range body {
		jobs = append(jobs, output.Job{
			Company:   company,
			Title:     j.Title,
			UpdatedAt: time.UnixMilli(int64(j.UpdatedAt)),
		})
	}

	return jobs, nil
}

func checkThenDecode[T any](company db.Company, res *http.Response, reqError error) (body T, err error) {
	if reqError != nil {
		return body, fmt.Errorf("failed to request site=%s for company=%s: %w", company.Site, company.Name, reqError)
	}
	if res.StatusCode != http.StatusOK {
		return body, fmt.Errorf("invalid site=%s for company=%s: code=%d", company.Site, company.Name, res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return body, fmt.Errorf("failed to decode response from site=%s for company=%s: %w", company.Site, company.Name, err)
	}

	return body, nil
}

// TODO: https://www.polymer.co/
