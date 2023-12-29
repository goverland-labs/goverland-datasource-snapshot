package voting

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/rs/zerolog/log"
	"github.com/shutter-network/shutter/shlib/shcrypto"

	"github.com/goverland-labs/datasource-snapshot/internal/config"
	"github.com/goverland-labs/datasource-snapshot/internal/helpers"
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
	PrimaryType string    `json:"primaryType,omitempty"`
	VoteTypes   VoteTypes `json:"types"`
	Vote        Vote      `json:"message"`
}

func (t TypedData) ForSnapshot() TypedData {
	t.PrimaryType = ""
	t.VoteTypes.EIP712Domain = nil

	return t
}

type Domain struct {
	Name              string  `json:"name"`
	Version           string  `json:"version"`
	ChainId           *int    `json:"chainId,omitempty"`           // TODO do we need it?
	VerifyingContract *string `json:"verifyingContract,omitempty"` // TODO do we need it?
}

type VoteTypes struct {
	EIP712Domain []TypesDescription `json:"EIP712Domain,omitempty"`
	Vote         []TypesDescription `json:"Vote"`
}

type TypesDescription struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var EIP712Domain = []TypesDescription{
	{Name: "name", Type: "string"},
	{Name: "version", Type: "string"},
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

var VoteNumberTypes2 = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "bytes32"},
	{Name: "choice", Type: "uint32"},
	{Name: "reason", Type: "string"},
	{Name: "app", Type: "string"},
	{Name: "metadata", Type: "string"},
}

var VoteNumbersTypes2 = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "bytes32"},
	{Name: "choice", Type: "uint32[]"},
	{Name: "reason", Type: "string"},
	{Name: "app", Type: "string"},
	{Name: "metadata", Type: "string"},
}

var VoteStringTypes2 = []TypesDescription{
	{Name: "from", Type: "address"},
	{Name: "space", Type: "string"},
	{Name: "timestamp", Type: "uint64"},
	{Name: "proposal", Type: "bytes32"},
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
	Choice    json.RawMessage `json:"choice"`
	Reason    string          `json:"reason"`
	App       string          `json:"app"`
	Metadata  string          `json:"metadata"`
}

type TypedSignDataBuilder struct {
	cfg config.Snapshot
}

func NewTypedSignDataBuilder(cfg config.Snapshot) *TypedSignDataBuilder {
	return &TypedSignDataBuilder{
		cfg: cfg,
	}
}

func (t *TypedSignDataBuilder) Build(checksumVoter string, reason *string, choice json.RawMessage, pFragment *client.ProposalFragment) (TypedData, error) {
	isShutter := pFragment.Privacy != nil && *pFragment.Privacy == shutterPrivacy

	isTypes2 := strings.HasPrefix(pFragment.ID, "0x")

	types := helpers.Ternary(isTypes2, VoteNumberTypes2, VoteNumberTypes)
	if t.isProposalType(pFragment.Type, approvalProposalType, rankedChoiceProposalType) {
		types = helpers.Ternary(isTypes2, VoteNumbersTypes2, VoteNumbersTypes)
	}
	if isShutter || t.isProposalType(pFragment.Type, quadraticProposalType, weightedProposalType) {
		types = helpers.Ternary(isTypes2, VoteStringTypes2, VoteStringTypes)
	}

	if reason == nil {
		reason = helpers.Ptr("")
	}

	if isShutter {
		choiceForShutter, err := t.getChoiceForShutter(string(choice), pFragment)
		if err != nil {
			return TypedData{}, fmt.Errorf("failed to get choice for shutter: %w", err)
		}

		choiceStr, err := t.encodeShutterChoice(choiceForShutter, pFragment.ID)
		if err != nil {
			return TypedData{}, fmt.Errorf("failed to encode choice: %w", err)
		}

		choice = []byte(fmt.Sprintf(`"%s"`, choiceStr))
	}

	td := TypedData{
		Domain: Domain{
			Name:    "snapshot",
			Version: "0.1.4", // TODO to config
		},
		PrimaryType: "Vote",
		VoteTypes: VoteTypes{
			Vote:         types,
			EIP712Domain: EIP712Domain,
		},
		Vote: Vote{
			From:      checksumVoter,
			Space:     pFragment.Space.ID,
			Timestamp: time.Now().Unix(),
			Proposal:  pFragment.ID,
			Choice:    choice,
			Reason:    *reason,
			App:       "goverland",
			Metadata:  "{}",
		},
	}

	return td, nil
}

func (t *TypedSignDataBuilder) getChoiceForShutter(choice string, pFragment *client.ProposalFragment) (string, error) {
	if !t.isProposalType(pFragment.Type, quadraticProposalType, weightedProposalType) {
		return choice, nil
	}

	var choiceForShutter string
	if err := json.Unmarshal([]byte(choice), &choiceForShutter); err != nil {
		return "", fmt.Errorf("failed to convert choice %s for shutter: %w", choice, err)
	}

	log.Info().
		Str("choice", choice).
		Str("choiceForShutter", choiceForShutter).
		Msg("get choice for string shutter")

	return choiceForShutter, nil
}

func (t *TypedSignDataBuilder) encodeShutterChoice(choice string, proposalID string) (string, error) {
	if t.cfg.ViteShutterEonPubKey == "" {
		return "", fmt.Errorf("vite shutter eon public key is empty, please set it in config")
	}

	var eonKey shcrypto.EonPublicKey
	eonKeyBytes, err := hex.DecodeString(t.cfg.ViteShutterEonPubKey)
	if err != nil {
		return "", err
	}
	err = eonKey.GobDecode(eonKeyBytes)
	if err != nil {
		return "", err
	}

	proposalBytes, err := hex.DecodeString(strings.TrimPrefix(proposalID, "0x"))
	if err != nil {
		return "", err
	}
	epochID := shcrypto.ComputeEpochID(proposalBytes)
	if err != nil {
		return "", err
	}
	sigma, err := shcrypto.RandomSigma(rand.Reader)
	if err != nil {
		return "", err
	}

	encrypt := shcrypto.Encrypt([]byte(choice), &eonKey, epochID, sigma)
	mEncrypted := encrypt.Marshal()

	return "0x" + hex.EncodeToString(mEncrypted), nil
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
