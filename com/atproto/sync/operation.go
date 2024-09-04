package sync

// internalOperation is the type for the 'Operation' field of [internalSubscriptionMessageHeader].
type internalOperation int64

const (
	internalOperationRegular internalOperation = 1
	internalOperationError   internalOperation = -1
)
