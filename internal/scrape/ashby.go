package scrape

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
)

type ashbyScraper struct {
	name string
}

func newAshbyScraper() Scraper {
	return ashbyScraper{
		name: string(db.SiteAshby),
	}
}

func (s ashbyScraper) Scrape(companyName, site string) ([]*proto.ScrapeBatch_Job, error) {
	if site != s.name {
		return nil, nil
	}
	return s.ScrapeAll(companyName)
}

func (s ashbyScraper) ScrapeAll(companyName string) ([]*proto.ScrapeBatch_Job, error) {
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
	}](companyName, s.name, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]*proto.ScrapeBatch_Job, 0)
	for _, j := range body.Data.JobBoard.JobPostings {
		jobs = append(jobs, &proto.ScrapeBatch_Job{
			Company: companyName,
			Title:   strings.Trim(j.Title, " "),
		})
	}

	return jobs, nil
}

func (s ashbyScraper) Sites() []string {
	return []string{
		s.name,
	}
}
