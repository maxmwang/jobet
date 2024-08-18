package discord

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/internal/helpers"
	"github.com/maxmwang/jobet/internal/proto"
)

type squashedJob struct {
	company   string
	title     string
	updatedAt int64
	count     int
}

func (s squashedJob) String() string {
	sb := strings.Builder{}
	if s.updatedAt == 0 {
		sb.WriteString(fmt.Sprintf("%45s", s.company))
	} else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := time.Unix(s.updatedAt, 0)
		sb.WriteString(fmt.Sprintf("%23s: %20s", t.In(loc).Format(time.DateTime+" MST"), s.company))
	}

	if s.count == 1 {
		sb.WriteString(":     ")
	} else {
		sb.WriteString(fmt.Sprintf(": %2dx ", s.count))
	}

	sb.WriteString(fmt.Sprintf("%v", s.title))

	return sb.String()
}

type squashedJobs struct {
	jobs []squashedJob
}

func newSquashedJobs(batch *proto.ScrapeBatch) *squashedJobs {
	s := &squashedJobs{jobs: make([]squashedJob, 0)}

	for _, v := range batch.Jobs {
		s.addJob(v)
	}

	return s
}

func (s *squashedJobs) addJob(job *proto.ScrapeBatch_Job) {
	for i, v := range s.jobs {
		if v.company == job.Company && v.title == job.Title {
			s.jobs[i].count++
			s.jobs[i].updatedAt = max(s.jobs[i].updatedAt, job.UpdatedAt)
			return
		}
	}
	s.jobs = append(s.jobs, squashedJob{
		company:   job.Company,
		title:     job.Title,
		updatedAt: job.UpdatedAt,
		count:     1,
	})
}

func (s *squashedJobs) String() string {
	slices.SortFunc(s.jobs, func(a, b squashedJob) int {
		if a.updatedAt > b.updatedAt {
			return 1
		}
		return -1
	})

	sb := strings.Builder{}
	for _, j := range s.jobs {
		if helpers.JobIsTarget(j.title) {
			sb.WriteString(j.String())
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
