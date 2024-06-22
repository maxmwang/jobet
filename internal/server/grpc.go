package server

import (
	"context"
	"net"
	"strings"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
	"github.com/maxmwang/jobet/internal/scrape"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type JobetServer struct {
	proto.UnimplementedJobetServer

	scrapers map[string]scrape.Scraper
	q        *db.Queries
}

func NewJobetServer(q *db.Queries) *JobetServer {
	return &JobetServer{
		scrapers: map[string]scrape.Scraper{
			scrape.SiteAshby:      scrape.NewAshbyScraper(),
			scrape.SiteGreenhouse: scrape.NewGreenhouseScraper(),
			scrape.SiteLever:      scrape.NewLeverScraper(),
		},
		q: q,
	}
}

func (s *JobetServer) Start() {
	lis, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterJobetServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *JobetServer) Probe(ctx context.Context, req *proto.ProbeRequest) (*proto.ProbeReply, error) {
	existingSites := make(map[string]struct{})
	existing, err := s.q.GetCompaniesByName(ctx, req.Name)
	if err == nil {
		for _, c := range existing {
			existingSites[c.Site] = struct{}{}
		}
	}

	res := &proto.ProbeReply{
		Results: make([]*proto.ProbeReply_ProbeSiteResult, 0),
	}
	for site, scraper := range s.scrapers {
		added := false
		_, exists := existingSites[site]

		jobs, err := scraper.Scrape(req.Name)
		if err != nil {
			log.Warn().
				Str("name", req.Name).
				Str("site", site).
				Err(err).
				Msg("could not scrape")
			continue
		}
		target := 0
		for _, j := range jobs {
			if isTarget(j.Title) {
				target++
			}
		}
		if !req.Dry && !exists && len(jobs) > 0 {
			err := s.q.AddCompany(ctx, db.AddCompanyParams{
				Name:     req.Name,
				Alias:    req.Alias,
				Site:     site,
				Priority: req.Priority,
			})
			if err == nil {
				added = true
			}
		}
		res.Results = append(res.Results, &proto.ProbeReply_ProbeSiteResult{
			Site:   site,
			Added:  added,
			Exists: exists,
			Count:  int32(len(jobs)),
			Target: int32(target),
		})
	}

	return res, nil
}

func isTarget(title string) bool {
	if strings.Index(title, "Intern") == strings.Index(title, "Internal") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Index(title, "Intern") == strings.Index(title, "International") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Contains(title, "Software") && strings.Contains(title, "Intern") {
		return true
	}
	if strings.Contains(title, "Platform") && strings.Contains(title, "Intern") {
		return true
	}

	return false
}
