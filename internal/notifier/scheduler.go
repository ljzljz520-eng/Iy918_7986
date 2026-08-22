package notifier

import (
	"context"
	"errors"
	"sort"
)

type Job struct {
	ID       string
	Request  Request
	Priority int
}

type Queue struct {
	jobs []Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make([]Job, 0)}
}

func (q *Queue) Enqueue(job Job) error {
	if job.ID == "" || job.Request.RecordID == "" {
		return errors.New("job identity is required")
	}
	q.jobs = append(q.jobs, job)
	sort.SliceStable(q.jobs, func(i, j int) bool {
		if q.jobs[i].Priority == q.jobs[j].Priority {
			return q.jobs[i].ID < q.jobs[j].ID
		}
		return q.jobs[i].Priority > q.jobs[j].Priority
	})
	return nil
}

func (q *Queue) Len() int {
	return len(q.jobs)
}

func (q *Queue) Peek() (Job, error) {
	if len(q.jobs) == 0 {
		return Job{}, errors.New("queue is empty")
	}
	return q.jobs[0], nil
}

func (q *Queue) Pop() (Job, error) {
	job, err := q.Peek()
	if err != nil {
		return Job{}, err
	}
	q.jobs = q.jobs[1:]
	return job, nil
}

func (s *Service) Drain(ctx context.Context, queue *Queue) []Result {
	results := make([]Result, 0)
	for queue.Len() > 0 {
		if ctx.Err() != nil {
			break
		}
		job, err := queue.Pop()
		if err != nil {
			break
		}
		results = append(results, s.Notify(ctx, job.Request))
	}
	return results
}

func CountTimeouts(results []Result) int {
	count := 0
	for _, result := range results {
		if IsTimeout(result) {
			count++
		}
	}
	return count
}
