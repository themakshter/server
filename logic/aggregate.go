package logic

import (
	"errors"
	impact "github.com/impactasaurus/server"
)

func mean(in []float32) float32 {
	s := sum(in)
	return s / float32(len(in))
}

func sum(in []float32) float32 {
	var total float32 = 0
	for _, value := range in {
		total += value
	}
	return total
}

func aggregate(in []float32, aggregation impact.Aggregation) (float32, error) {
	switch aggregation {
	case impact.MEAN:
		return mean(in), nil
	case impact.SUM:
		return sum(in), nil
	default:
		return 0, errors.New("Unknown aggregation")
	}
}

func GetCategoryAggregates(m impact.Meeting, os impact.OutcomeSet) ([]impact.CategoryAggregate, error) {
	out := make([]impact.CategoryAggregate, 0, len(os.Categories))
	for _, c := range os.Categories {
		vals := make([]float32, 0, len(m.Answers))
		for _, a := range m.Answers {
			if a.Type != impact.INT {
				continue
			}
			q := os.GetQuestion(a.QuestionID)
			if q.CategoryID == c.ID {
				f, err := a.ToFloat()
				if err != nil {
					return nil, err
				}
				vals = append(vals, f)
			}
		}
		ag, err := aggregate(vals, c.Aggregation)
		if err != nil {
			return nil, err
		}
		if len(vals) > 0 {
			out = append(out, impact.CategoryAggregate{
				CategoryID: c.ID,
				Aggregate:  ag,
			})
		}
	}
	return out, nil
}
