package prober

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5"

	"github.com/maxmwang/jobet/internal/config"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
	"github.com/maxmwang/jobet/internal/scrape"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedProberServer

	scraper  scrape.Scraper
	dbClient *db.Queries

	port string
}

func NewServer(ctx context.Context, cfg config.Config, scraper scrape.Scraper, dbClient *db.Queries) *Server {
	return &Server{
		scraper:  scraper,
		dbClient: dbClient,
		port:     cfg.PollerPort,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", "localhost:"+s.port)
	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	server := grpc.NewServer()
	proto.RegisterProberServer(server, s)

	if err = server.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to register server")
	}
	return nil
}

func (s *Server) Probe(ctx context.Context, req *proto.ProbeRequest) (*proto.ProbeReply, error) {
	if req.Site == "" {
		return s.probeAll(ctx, req)
	}
	return s.probe(ctx, req)
}

func (s *Server) probeAll(ctx context.Context, req *proto.ProbeRequest) (*proto.ProbeReply, error) {
	existing := make(map[string]db.Company)
	if rows, err := s.dbClient.GetCompaniesByName(ctx, req.Company); err == nil {
		for _, c := range rows {
			existing[string(c.Site)] = c
		}
	}

	res := &proto.ProbeReply{
		Results: make([]*proto.ProbeReply_Result, 0),
	}
	for _, site := range s.scraper.Sites() {
		_, exists := existing[site]

		jobs, err := s.scraper.Scrape(req.Company, site)
		if err != nil {
			log.Warn().
				Str("company", req.Company).
				Str("site", site).
				Err(err).
				Msg("failed to scrape for probe")
			continue
		}
		res.Results = append(res.Results, &proto.ProbeReply_Result{
			Site:     site,
			Exists:   exists,
			Priority: existing[site].Priority,
			Count:    int32(len(jobs)),
		})
	}

	return res, nil
}

func (s *Server) probe(ctx context.Context, req *proto.ProbeRequest) (*proto.ProbeReply, error) {
	rows, err := s.dbClient.GetCompanyByNameAndSite(ctx, db.GetCompanyByNameAndSiteParams{
		Name: req.Company,
		Site: db.Site(req.Site),
	})
	exists := !errors.Is(err, pgx.ErrNoRows)

	res := &proto.ProbeReply{
		Results: make([]*proto.ProbeReply_Result, 0),
	}
	jobs, err := s.scraper.Scrape(req.Company, req.Site)
	if err != nil {
		log.Warn().
			Str("company", req.Company).
			Str("site", req.Site).
			Err(err).
			Msg("failed to scrape for probe")
	} else {
		res.Results = append(res.Results, &proto.ProbeReply_Result{
			Site:     req.Site,
			Exists:   exists,
			Priority: rows.Priority,
			Count:    int32(len(jobs)),
		})
	}

	return res, nil
}
