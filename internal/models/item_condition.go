package models

type ItemCondition string

const (
	ConditionNew  ItemCondition = "new"
	ConditionUsed ItemCondition = "used"
)

func (c ItemCondition) IsValid() bool {
	return c != ConditionNew && c != ConditionUsed
}
