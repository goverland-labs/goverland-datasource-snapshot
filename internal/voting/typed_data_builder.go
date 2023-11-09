package voting

import (
	"encoding/json"
	"time"

	"github.com/goverland-labs/sdk-snapshot-go/client"
)

const (
	shutterPrivacy = "shutter"
	app            = "goverland"

	singleChoiceProposalType ProposalType = "single-choice"
	approvalProposalType     ProposalType = "approval"
	quadraticProposalType    ProposalType = "quadratic"
	rankedChoiceProposalType ProposalType = "ranked-choice"
	weightedProposalType     ProposalType = "weighted"
	basicProposalType        ProposalType = "basic"
)

type ProposalType string

type TypedData struct {
	Domain      Domain    `json:"domain"`
	PrimaryType string    `json:"primaryType"`
	VoteTypes   VoteTypes `json:"types"`
	Vote        Vote      `json:"message"`
}

type Domain struct {
	Name              string  `json:"name"`
	Version           string  `json:"version"`
	ChainId           *int    `json:"chainId,omitempty"`           // TODO
	VerifyingContract *string `json:"verifyingContract,omitempty"` // TODO
}

type VoteTypes struct {
	Vote []TypesDescription `json:"Vote"`
}

type TypesDescription struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var VoteNumberTypes = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "string"},
	{Name: "choice", Type: "uint32"},
	{Name: "reason", Type: "string"},
	{Name: "app", Type: "string"},
	{Name: "metadata", Type: "string"},
}

var VoteNumbersTypes = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "string"},
	{Name: "choice", Type: "uint32[]"},
	{Name: "reason", Type: "string"},
	{Name: "app", Type: "string"},
	{Name: "metadata", Type: "string"},
}

var VoteStringTypes = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "string"},
	{Name: "choice", Type: "string"},
	{Name: "reason", Type: "string"},
	{Name: "app", Type: "string"},
	{Name: "metadata", Type: "string"},
}

type Vote struct {
	From      string          `json:"from"`
	Space     string          `json:"space"`
	Timestamp int64           `json:"timestamp"`
	Proposal  string          `json:"proposal"`
	Choice    json.RawMessage `json:"choice"` // TODO think about type
	Reason    string          `json:"reason"`
	App       string          `json:"app"`
	Metadata  string          `json:"metadata"`
}

type TypedSignDataBuilder struct {
}

func NewTypedSignDataBuilder() *TypedSignDataBuilder {
	return &TypedSignDataBuilder{}
}

// Build Now handle only 0x... proposal format
func (t *TypedSignDataBuilder) Build(prepareParams *PrepareParams, pFragment *client.ProposalFragment) TypedData {
	isShutter := pFragment.Privacy != nil && *pFragment.Privacy == shutterPrivacy

	types := VoteNumberTypes
	if t.isProposalType(pFragment.Type, approvalProposalType, rankedChoiceProposalType) {
		types = VoteNumbersTypes
	}
	if isShutter || t.isProposalType(pFragment.Type, quadraticProposalType, weightedProposalType) {
		types = VoteStringTypes
	}

	reason := ""
	if prepareParams.Reason != nil {
		reason = *prepareParams.Reason
	}

	td := TypedData{
		Domain: Domain{
			Name:    app,
			Version: "0.5", // TODO to config
		},
		PrimaryType: "Vote",
		VoteTypes: VoteTypes{
			Vote: types,
		},
		Vote: Vote{
			From:      prepareParams.Voter,
			Space:     pFragment.Space.ID,
			Timestamp: time.Now().Unix(),
			Proposal:  pFragment.ID,
			Choice:    prepareParams.Choice,
			Reason:    reason,
			App:       app,
			Metadata:  "",
		},
	}

	return td
}

func (t *TypedSignDataBuilder) isProposalType(proposalType *string, proposalTypes ...ProposalType) bool {
	if proposalType == nil {
		return false
	}

	for _, t := range proposalTypes {
		if *proposalType == string(t) {
			return true
		}
	}

	return false
}
