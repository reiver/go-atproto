package sync

// Operation is the type for the 'Operation' field of [SubscriptionMessageHeader].
type Operation int64

const (
	OperationRegular Operation = 1
	OperationError   Operation = -1
)
