package server

import (
	"context"
	"net"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/helpers"
	"github.com/maxmwang/jobet/internal/proto"
	"github.com/maxmwang/jobet/internal/scraper"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type JobetServer struct {
	proto.UnimplementedJobetServer

	scraper scraper.Scraper
	q       *db.Queries
}

func NewJobetServer(q *db.Queries) *JobetServer {
	return &JobetServer{
		scraper: scraper.NewDefault(),
		q:       q,
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
	for _, site := range s.scraper.Sites() {
		added := false
		_, exists := existingSites[site]

		jobs, err := s.scraper.Scrape(req.Name, site)
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
			if helpers.JobIsTarget(j.Title) {
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
