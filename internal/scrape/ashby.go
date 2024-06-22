package scrape

import (
	"fmt"
	"net/http"
	"strings"
)

type AshbyScraper struct{}

func NewAshbyScraper() AshbyScraper {
	return AshbyScraper{}
}

func (a AshbyScraper) Scrape(companyName string) ([]Job, error) {
	query := strings.NewReader(fmt.Sprintf(`{"operationName":"ApiJobBoardWithTeams","variables":{"organizationHostedJobsPageName":"%s"},"query":"query ApiJobBoardWithTeams($organizationHostedJobsPageName: String!) {\n  jobBoard: jobBoardWithTeams(\n    organizationHostedJobsPageName: $organizationHostedJobsPageName\n  ) {\n    jobPostings {\n      title\n     }\n  }\n}"}`, companyName))
	res, err := http.Post("https://jobs.ashbyhq.com/api/non-user-graphql?op=ApiJobBoardWithTeams", "application/json", query)

	body, err := checkThenDecode[struct {
		Data struct {
			JobBoard struct {
				JobPostings []struct {
					Title string `json:"title"`
				} `json:"jobPostings"`
			} `json:"jobBoard"`
		} `json:"data"`
	}](companyName, SiteAshby, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0)
	for _, j := range body.Data.JobBoard.JobPostings {
		jobs = append(jobs, Job{
			Title: j.Title,
		})
	}

	return jobs, nil
}
