package domain

import "time"

type Statistics struct {
	TasksCreated      int
	TasksCompleted    int
	CompletedRate     *float64
	AvgCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	completedRate *float64,
	avgCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:      tasksCreated,
		TasksCompleted:    tasksCompleted,
		CompletedRate:     completedRate,
		AvgCompletionTime: avgCompletionTime,
	}
}
