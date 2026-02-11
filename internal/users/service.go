package users

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	snapshotSDK snapshotSDK
}

func NewService(snapshotSDK snapshotSDK) *Service {
	return &Service{
		snapshotSDK: snapshotSDK,
	}
}

func (s *Service) GetUsersInfo(ctx context.Context, addresses []string) ([]Info, error) {
	converted := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		converted = append(converted, strings.ToLower(addr))
	}

	list, err := s.snapshotSDK.ListUsers(ctx, converted)
	if err != nil {
		return nil, fmt.Errorf("s.snapshotSDK.ListStatements: %w", err)
	}

	result := make([]Info, 0, len(list))
	for _, user := range list {
		var about string
		if val := user.GetAbout(); val != nil {
			about = *val
		}

		result = append(result, Info{
			Address: strings.ToLower(user.GetID()),
			About:   about,
		})
	}

	return result, nil
}
