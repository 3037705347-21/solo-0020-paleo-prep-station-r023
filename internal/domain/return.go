package domain

import (
	"time"

	"example.com/paleo-prep-station/internal/contracts"
)

type ReturnCondition string

const (
	ReturnIntact  ReturnCondition = "intact"
	ReturnDamaged ReturnCondition = "damaged"
)

type ReturnRecord struct {
	ID         string          `json:"id"`
	SpecimenID string          `json:"specimen_id"`
	ReleaseID  string          `json:"release_id"`
	Condition  ReturnCondition `json:"condition"`
	Notes      string          `json:"notes"`
	ReceivedBy string          `json:"received_by"`
	CreatedAt  time.Time       `json:"created_at"`
}

func NewReturnRecord(
	id string,
	specimen Specimen,
	release ReleaseRecord,
	condition string,
	notes string,
	receivedBy string,
	now time.Time,
) (ReturnRecord, error) {
	if id == "" {
		return ReturnRecord{}, contracts.Internal("missing_generated_id", "归还标识未生成", nil)
	}
	if specimen.Status != SpecimenAvailable {
		return ReturnRecord{}, contracts.Conflict("specimen_not_on_loan", "只有已借出标本才能登记归还")
	}
	if release.ID == "" || release.SpecimenID != specimen.ID {
		return ReturnRecord{}, contracts.Internal("release_specimen_mismatch", "放行记录与标本不一致", nil)
	}
	conditionValue, err := parseReturnCondition(condition)
	if err != nil {
		return ReturnRecord{}, err
	}
	receiver, err := requiredText("归还接收人", receivedBy, 80)
	if err != nil {
		return ReturnRecord{}, err
	}
	noteValue, err := optionalText("归还备注", notes, 500)
	if err != nil {
		return ReturnRecord{}, err
	}
	return ReturnRecord{
		ID:         id,
		SpecimenID: specimen.ID,
		ReleaseID:  release.ID,
		Condition:  conditionValue,
		Notes:      noteValue,
		ReceivedBy: receiver,
		CreatedAt:  now.UTC(),
	}, nil
}

func parseReturnCondition(value string) (ReturnCondition, error) {
	switch ReturnCondition(value) {
	case ReturnIntact:
		return ReturnIntact, nil
	case ReturnDamaged:
		return ReturnDamaged, nil
	default:
		return "", contracts.Invalid("invalid_return_condition", "归还检查结论必须是 intact 或 damaged")
	}
}
