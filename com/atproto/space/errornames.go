package space

// The names of the XRPC errors the "com.atproto.space.*" methods may answer with — see [xrpc.ErrorName].
const (
	ErrorNameAccountNotAuthorized     = "AccountNotAuthorized"
	ErrorNameAppNotAuthorized         = "AppNotAuthorized"
	ErrorNameBlobNotFound             = "BlobNotFound"
	ErrorNameFutureRev                = "FutureRev"
	ErrorNameInvalidClientAttestation = "InvalidClientAttestation"
	ErrorNameInvalidCursor            = "InvalidCursor"
	ErrorNameInvalidDelegationToken   = "InvalidDelegationToken"
	ErrorNameInvalidRepoFile          = "InvalidRepoFile"
	ErrorNameInvalidSpaceMac          = "InvalidSpaceMac"
	ErrorNameInvalidSpaceMeta         = "InvalidSpaceMeta"
	ErrorNameNotAuthorized            = "NotAuthorized"
	ErrorNameRecordAlreadyExists      = "RecordAlreadyExists"
	ErrorNameRecordNotFound           = "RecordNotFound"
	ErrorNameRegistrationNotConfirmed = "RegistrationNotConfirmed"
	ErrorNameRepoDeactivated          = "RepoDeactivated"
	ErrorNameRepoMismatch             = "RepoMismatch"
	ErrorNameRepoNotFound             = "RepoNotFound"
	ErrorNameRepoSuspended            = "RepoSuspended"
	ErrorNameRepoTakendown            = "RepoTakendown"
	ErrorNameSpaceAlreadyExists       = "SpaceAlreadyExists"
	ErrorNameSpaceDeleted             = "SpaceDeleted"
	ErrorNameSpaceNotFound            = "SpaceNotFound"
	ErrorNameUnsupportedPolicy        = "UnsupportedPolicy"
)
