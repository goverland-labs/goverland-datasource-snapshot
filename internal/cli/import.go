package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

const (
	ImportCommandName = "import"

	ImportTypeUnspecified = "unspecified"
	ImportTypeSpace       = "space"
	ImportTypeProposal    = "proposal"
	ImportTypeVote        = "vote"
)

type ImportType string

type Import struct {
	base

	Spaces    *db.SpaceService
	Proposals *db.ProposalService
	Votes     *db.VoteService

	votesBatch []db.Vote
}

func (c *Import) GetName() string {
	return ImportCommandName
}

func (c *Import) GetArguments() ArgumentsDetails {
	return ArgumentsDetails{
		"type":    "import type: space, proposal, vote",
		"path":    "source absolute path",
		"limit":   "number of maximum rows for processing",
		"offset":  "how many skip",
		"timeout": "time duration for <FIXME>",
	}
}

func (c *Import) ParseArgs(args ...string) (Arguments, error) {
	return c.parseArgs(c, args...)
}

func (c *Import) Execute(args Arguments) error {
	start := time.Now()
	log.Info().Msg("import started")

	importType, err := c.getImportType(args)
	if err != nil {
		return err
	}

	ttl, err := c.getInputTimeout(args)
	if err != nil {
		return err
	}

	path, err := c.getInputPath(args)
	if err != nil {
		return fmt.Errorf("path: %w", err)
	}

	limit, err := c.getLimit(args)
	if err != nil {
		return err
	}

	offset, err := c.getOffset(args)
	if err != nil {
		return err
	}

	f, _ := os.Open(path)
	reader := csv.NewReader(f)
	var idx int64 = 0
	for {
		idx++
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("ERROR: %d: %s\n", idx, err.Error())
			break
		}

		if idx <= offset {
			continue
		}

		if limit != 0 && idx > offset+limit {
			break
		}

		if idx%1000 == 0 {
			log.Info().Msgf("process %d lines", idx)
		}

		switch importType {
		case ImportTypeSpace:
			err = c.processSpace(line, ttl)
		case ImportTypeProposal:
			err = c.processProposal(line, ttl)
		case ImportTypeVote:
			err = c.processVotes(line, ttl)
		default:
			panic(fmt.Sprintf("import type is not implemented: %s", importType))
		}

		if err != nil {
			log.Error().Err(err).Msgf("upsert %s: %d: %s", importType, idx, line[0])
		}
	}

	// todo: separate commands
	if len(c.votesBatch) != 0 {
		err := c.Votes.BatchCreate(c.votesBatch)
		if err != nil {
			return fmt.Errorf("batch create: %w", err)
		}
	}

	log.Info().Msgf("import finished. Took: %v", time.Since(start))

	return nil
}

func (c *Import) processSpace(line []string, ttl *time.Duration) error {
	space := &db.Space{
		ID:        line[0],
		CreatedAt: getTimeFromString(line[5]),
		UpdatedAt: getTimeFromString(line[6]),
		Snapshot:  json.RawMessage(line[2]),
	}

	err := c.Spaces.Upsert(space)
	if err != nil {
		return fmt.Errorf("upsert space: %w", err)
	}

	if ttl != nil {
		<-time.After(*ttl)
	}

	return nil
}

func (c *Import) processProposal(line []string, ttl *time.Duration) error {
	strategies, err := helpers.Unmarshal([]*client.StrategyFragment{}, json.RawMessage(line[8]))
	if err != nil {
		return fmt.Errorf("convert strategies: %w", err)
	}

	validations, err := helpers.Unmarshal(client.ValidationFragment{}, json.RawMessage(line[9]))
	if err != nil {
		return fmt.Errorf("convert validations: %w", err)
	}

	choices, err := helpers.Unmarshal([]*string{}, json.RawMessage(line[14]))
	if err != nil {
		return fmt.Errorf("convert choices: %w", err)
	}

	scores, err := helpers.Unmarshal([]*float64{}, json.RawMessage(line[22]))
	if err != nil {
		return fmt.Errorf("convert scores: %w", err)
	}

	pf := client.ProposalFragment{
		ID:            line[0],
		Ipfs:          helpers.Ptr(line[1]),
		Author:        line[2],
		Created:       getUnixFromString(line[3]),
		Network:       line[5],
		Symbol:        line[6],
		Type:          helpers.Ptr(line[7]),
		Strategies:    strategies,
		Validation:    &validations,
		Title:         line[11],
		Body:          helpers.Ptr(line[12]),
		Discussion:    line[13],
		Choices:       choices,
		Start:         getUnixFromString(line[15]),
		End:           getUnixFromString(line[16]),
		Quorum:        getFloat64FromString(line[18]),
		Privacy:       helpers.Ptr(line[19]),
		Snapshot:      helpers.Ptr(line[20]),
		State:         helpers.Ptr(line[24]),
		Link:          helpers.Ptr(fmt.Sprintf("https://snapshot.org/#/%s/proposal/%s", line[4], line[0])),
		App:           helpers.Ptr(line[21]),
		Scores:        scores,
		ScoresState:   helpers.Ptr(line[24]),
		ScoresTotal:   helpers.Ptr(getFloat64FromString(line[25])),
		ScoresUpdated: helpers.Ptr(getUnixFromString(line[26])),
		Votes:         helpers.Ptr(getUnixFromString(line[27])),
		Space: &client.SpaceIdentifierFragment{
			ID: line[4],
		},
	}

	snapshot, _ := json.Marshal(pf)

	pr := &db.Proposal{
		ID:            pf.ID,
		SpaceID:       pf.Space.ID,
		CreatedAt:     getTimeFromString(line[3]),
		UpdatedAt:     time.Now(),
		Snapshot:      snapshot,
		VoteProcessed: false,
	}

	err = c.Proposals.Upsert(pr)
	if err != nil {
		return fmt.Errorf("upsert proposal: %w", err)
	}

	if ttl != nil {
		<-time.After(*ttl)
	}

	return nil
}

func (c *Import) processVotes(line []string, ttl *time.Duration) error {
	vpByStrategy, err := helpers.Unmarshal([]float64{}, json.RawMessage(line[11]))
	if err != nil {
		return fmt.Errorf("convert vp by strategy: %w", err)
	}

	vote := db.Vote{
		ID:           line[0],
		Ipfs:         line[1],
		CreatedAt:    getTimeFromString(line[3]),
		UpdatedAt:    time.Now(),
		Voter:        line[2],
		SpaceID:      line[4],
		ProposalID:   line[5],
		Choice:       prepareChoice(line[6]),
		Reason:       line[8],
		App:          line[9],
		Vp:           getFloat64FromString(line[10]),
		VpByStrategy: vpByStrategy,
		VpState:      line[12],
	}

	c.votesBatch = append(c.votesBatch, vote)

	if len(c.votesBatch) >= 1000 {
		err := c.Votes.BatchCreate(c.votesBatch)
		if err != nil {
			return fmt.Errorf("batch create: %w", err)
		}

		if ttl != nil {
			<-time.After(*ttl)
		}

		c.votesBatch = make([]db.Vote, 0, 1200)
	}

	return nil
}

func prepareChoice(str string) json.RawMessage {
	raw := json.RawMessage(str)

	var number int64
	if _, err := helpers.Unmarshal(number, raw); err == nil {
		return raw
	}

	if _, err := helpers.Unmarshal(map[string]int{}, raw); err == nil {
		return raw
	}

	if _, err := helpers.Unmarshal([]string{}, raw); err == nil {
		return raw
	}

	if _, err := helpers.Unmarshal([]int{}, raw); err == nil {
		return raw
	}

	// should be string
	val, _ := json.Marshal(str)

	return val
}

func getTimeFromString(val string) time.Time {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("convert %s to int", val)
		return time.Now()
	}

	return time.Unix(i, 0)
}

func getUnixFromString(val string) int64 {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("convert %s to int", val)
		return 0
	}

	return i
}
func getFloat64FromString(val string) float64 {
	i, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Error().Err(err).Msgf("convert %s to float", val)
		return 0
	}

	return i
}

func (c *Import) getImportType(args Arguments) (ImportType, error) {
	switch args.Get("type") {
	case ImportTypeSpace:
		return ImportTypeSpace, nil
	case ImportTypeProposal:
		return ImportTypeProposal, nil
	case ImportTypeVote:
		return ImportTypeVote, nil
	default:
		return ImportTypeUnspecified, fmt.Errorf("type has wrong format: %s", args.Get("type"))
	}
}

func (c *Import) getInputTimeout(args Arguments) (*time.Duration, error) {
	val := args.Get("timeout")
	if val == "" {
		return nil, nil
	}

	ttl, err := time.ParseDuration(val)
	if err != nil {
		return nil, err
	}

	return &ttl, nil
}

func (c *Import) getInputPath(args Arguments) (string, error) {
	src := args.Get("path")
	if src == "" {
		return "", errors.New("path is required")
	}

	if _, err := os.Stat(src); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("path: %w", err)
	}

	return src, nil
}

func (c *Import) getLimit(args Arguments) (int64, error) {
	limit := args.Get("limit")
	if limit == "" {
		return 0, nil
	}

	return strconv.ParseInt(limit, 0, 64)
}

func (c *Import) getOffset(args Arguments) (int64, error) {
	offset := args.Get("offset")
	if offset == "" {
		return 0, nil
	}

	return strconv.ParseInt(offset, 0, 64)
}
