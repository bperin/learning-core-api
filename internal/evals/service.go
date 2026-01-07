package evals

import (
	"context"
	"errors"
	"time"

	"learning-core-api/internal/evals/results"
	"learning-core-api/internal/evals/rules"
	"learning-core-api/internal/evals/runs"
	"learning-core-api/internal/generation"

	"github.com/google/uuid"
)

// EvalAggregate captures the final decision for an eval run.
type EvalAggregate struct {
	OverallPass  bool
	OverallScore *float64
}

// Service owns evaluation aggregation and gating logic.
type Service interface {
	AggregateAndFinalize(ctx context.Context, evalRunID uuid.UUID) (*EvalAggregate, error)
}

type service struct {
	runsRepo     runs.Repository
	rulesRepo    rules.Repository
	resultsRepo  results.Repository
	artifactRepo generation.Repository
}

// NewService creates a new eval aggregation service.
func NewService(
	runsRepo runs.Repository,
	rulesRepo rules.Repository,
	resultsRepo results.Repository,
	artifactRepo generation.Repository,
) Service {
	return &service{
		runsRepo:     runsRepo,
		rulesRepo:    rulesRepo,
		resultsRepo:  resultsRepo,
		artifactRepo: artifactRepo,
	}
}

// AggregateEvalRun computes overall pass/score for an eval run.
func AggregateEvalRun(
	rulesList []rules.Rule,
	resultsList []results.Result,
) EvalAggregate {
	resultByRule := make(map[uuid.UUID]results.Result, len(resultsList))
	for _, res := range resultsList {
		resultByRule[res.RuleID] = res
	}

	var totalWeight float64
	var weightedSum float64

	for _, rule := range rulesList {
		res, ok := resultByRule[rule.ID]
		if !ok {
			return EvalAggregate{OverallPass: false}
		}

		if rule.HardFail && !res.Pass {
			return EvalAggregate{OverallPass: false}
		}

		if res.Score != nil && rule.Weight > 0 {
			weight := float64(rule.Weight)
			totalWeight += weight
			weightedSum += float64(*res.Score) * weight
		}
	}

	if totalWeight == 0 {
		return EvalAggregate{OverallPass: true}
	}

	score := weightedSum / totalWeight
	return EvalAggregate{
		OverallPass:  true,
		OverallScore: &score,
	}
}

// AggregateAndFinalize loads rules/results, aggregates, and gates the artifact.
func (s *service) AggregateAndFinalize(ctx context.Context, evalRunID uuid.UUID) (*EvalAggregate, error) {
	if evalRunID == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	run, err := s.runsRepo.GetByID(ctx, evalRunID)
	if err != nil {
		return nil, err
	}

	if run.SuiteID == uuid.Nil {
		return nil, errors.New("suite ID is required for aggregation")
	}

	rulesList, err := s.rulesRepo.ListBySuite(ctx, run.SuiteID)
	if err != nil {
		return nil, err
	}

	resultsList, err := s.resultsRepo.ListByRun(ctx, evalRunID)
	if err != nil {
		return nil, err
	}

	aggregate := AggregateEvalRun(rulesList, resultsList)

	var overallScore *float32
	if aggregate.OverallScore != nil {
		score32 := float32(*aggregate.OverallScore)
		overallScore = &score32
	}

	overallPass := aggregate.OverallPass
	_, err = s.runsRepo.UpdateResult(ctx, evalRunID, "SUCCEEDED", &overallPass, overallScore, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if aggregate.OverallPass {
		if err := s.artifactRepo.UpdateArtifactStatus(ctx, run.ArtifactID, generation.ArtifactStatusApproved, &now, nil); err != nil {
			return nil, err
		}
	} else {
		if err := s.artifactRepo.UpdateArtifactStatus(ctx, run.ArtifactID, generation.ArtifactStatusRejected, nil, &now); err != nil {
			return nil, err
		}
	}

	return &aggregate, nil
}
