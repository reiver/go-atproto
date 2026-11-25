package space

// NSIDDefs is the NSID of the lexicon of shared definitions.
const NSIDDefs = "com.atproto.space.defs"

// The "$type" values of the members of the unions defined in "com.atproto.space.defs".
const (
	TypeRepoWrite   = NSIDDefs + "#repoWrite"
	TypeRepoRemoved = NSIDDefs + "#repoRemoved"
)

// Commit is "com.atproto.space.defs#commit" — a commit to the current state of a space repo.
//
// It is not signed; the repo host authenticates it with a MAC instead.
type Commit struct {
	// Ver is the commit format version — currently 1.
	Ver int64 `json:"ver"`

	// Repo is the DID of the repo.
	Repo string `json:"repo"`

	// Space is the space (a space-ref).
	Space string `json:"space"`

	// Rev is the repo revision (a TID).
	Rev string `json:"rev"`

	// Data is the SHA-256 digest of the repo's LtHash state.
	Data Bytes `json:"data"`
}

// RepoWrite is "com.atproto.space.defs#repoWrite" — a repo in the space advanced to a new revision.
//
// It is a member of the unions [ListWritersRepo] and [NotifySpaceUpdateUpdate], which add its "$type".
type RepoWrite struct {
	Repo    string `json:"repo"`
	RepoRev string `json:"repoRev"`
	Data    Bytes  `json:"data"`

	// SpaceRev is the space revision at which this update was sequenced.
	SpaceRev string `json:"spaceRev"`

	// PrevSpaceRev is the space revision before this update, if any.
	// It is set in notifySpaceUpdate, where a gap indicates missed notifications (which listWriters can recover).
	PrevSpaceRev string `json:"prevSpaceRev,omitempty"`
}

// RepoRemoved is "com.atproto.space.defs#repoRemoved" — a repo was removed from the writer set.
//
// It is a member of the unions [ListWritersRepo] and [NotifySpaceUpdateUpdate], which add its "$type".
type RepoRemoved struct {
	Repo string `json:"repo"`

	// SpaceRev is the space revision at which this update was sequenced.
	SpaceRev string `json:"spaceRev"`

	// PrevSpaceRev is the space revision before this update, if any.
	// It is set in notifySpaceUpdate, where a gap indicates missed notifications (which listWriters can recover).
	PrevSpaceRev string `json:"prevSpaceRev,omitempty"`
}
