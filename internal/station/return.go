package station

import (
	"context"

	"example.com/paleo-prep-station/internal/contracts"
	"example.com/paleo-prep-station/internal/domain"
	"example.com/paleo-prep-station/internal/repository"
)

func (s *Service) ReturnSpecimen(
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
		release, onLoan := m.OpenReleaseForSpecimen(specimenID)
		if !onLoan {
			return contracts.Conflict("specimen_not_on_loan", "标本当前不在借出状态，不能登记归还")
		}
		record, err := domain.NewReturnRecord(
			s.nextID("return"),
			release,
			specimen,
			request.Condition,
			request.Note,
			request.ReceivedBy,
			s.now(),
		)
		if err != nil {
			return err
		}
		if err := specimen.MarkReturned(record.Condition, s.now()); err != nil {
			return err
		}
		m.PutReturn(record)
		m.PutSpecimen(specimen)
		created = record
		return nil
	})
	return created, err
}
