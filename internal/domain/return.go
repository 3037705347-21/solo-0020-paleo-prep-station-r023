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

func ParseReturnCondition(value string) (ReturnCondition, error) {
	switch ReturnCondition(value) {
	case ReturnIntact:
		return ReturnIntact, nil
	case ReturnDamaged:
		return ReturnDamaged, nil
	default:
		return "", contracts.Invalid("invalid_condition", "归还检查结论必须是 intact 或 damaged")
	}
}

type ReturnRecord struct {
	ID         string          `json:"id"`
	ReleaseID  string          `json:"release_id"`
	SpecimenID string          `json:"specimen_id"`
	Condition  ReturnCondition `json:"condition"`
	Note       string          `json:"note"`
	ReceivedBy string          `json:"received_by"`
	CreatedAt  time.Time       `json:"created_at"`
}

func NewReturnRecord(
	id string,
	release ReleaseRecord,
	specimen Specimen,
	condition string,
	note string,
	receivedBy string,
	now time.Time,
) (ReturnRecord, error) {
	if id == "" {
		return ReturnRecord{}, contracts.Internal("missing_generated_id", "归还标识未生成", nil)
	}
	if release.ID == "" || release.SpecimenID != specimen.ID {
		return ReturnRecord{}, contracts.Internal("release_state_inconsistent", "放行记录与标本不一致", nil)
	}
	if specimen.Status != SpecimenAvailable {
		return ReturnRecord{}, contracts.Conflict("specimen_not_on_loan", "只有借出中的标本可以登记归还")
	}
	conclusion, err := ParseReturnCondition(condition)
	if err != nil {
		return ReturnRecord{}, err
	}
	noteValue, err := optionalText("归还说明", note, 500)
	if err != nil {
		return ReturnRecord{}, err
	}
	receiver, err := requiredText("接收人", receivedBy, 80)
	if err != nil {
		return ReturnRecord{}, err
	}
	return ReturnRecord{
		ID:         id,
		ReleaseID:  release.ID,
		SpecimenID: specimen.ID,
		Condition:  conclusion,
		Note:       noteValue,
		ReceivedBy: receiver,
		CreatedAt:  now.UTC(),
	}, nil
}
