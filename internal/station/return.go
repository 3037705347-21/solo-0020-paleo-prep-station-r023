package station

import (
	"context"

	"example.com/paleo-prep-station/internal/contracts"
	"example.com/paleo-prep-station/internal/domain"
	"example.com/paleo-prep-station/internal/repository"
)

func (s *Service) RegisterReturn(
	ctx context.Context,
	specimenID string,
	request contracts.ReturnRequest,
) (domain.ReturnRecord, error) {
	if err := ctx.Err(); err != nil {
		return domain.ReturnRecord{}, err
	}

	var created domain.ReturnRecord
	err := s.repo.Mutate(func(m *repository.Mutation) error {
		specimen, ok := m.GetSpecimen(specimenID)
		if !ok {
			return notFoundSpecimen(specimenID)
		}
		release, ok := m.LatestReleaseForSpecimen(specimenID)
		if !ok {
			return contracts.Conflict("specimen_not_on_loan", "标本尚未借出，不能登记归还")
		}
		if _, exists := m.ReturnForRelease(release.ID); exists {
			return contracts.Conflict("return_already_recorded", "该放行记录已登记归还，不能重复归还")
		}
		record, err := domain.NewReturnRecord(
			s.nextID("return"),
			specimen,
			release,
			request.Condition,
			request.Notes,
			request.ReceivedBy,
			s.now(),
		)
		if err != nil {
			return err
		}
		if err := specimen.ResolveReturn(record.Condition, s.now()); err != nil {
			return err
		}
		m.PutReturn(record)
		m.PutSpecimen(specimen)
		created = record
		return nil
	})
	return created, err
}
