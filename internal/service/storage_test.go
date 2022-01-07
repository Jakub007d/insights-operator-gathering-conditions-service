package service_test

import (
	"encoding/json"
	"testing"

	"github.com/redhatinsights/insights-operator-conditional-gathering/internal/service"
	"github.com/stretchr/testify/assert"
)

const (
	validRulesFile = "rules.json"
	rulesFolder    = "testdata"
)

func TestStorage(t *testing.T) {
	type testCase struct {
		name          string
		rulesFile     string
		expectedRules service.Rules
	}

	testCases := []testCase{
		{
			name:          "file exists and is empty",
			rulesFile:     "empty.json",
			expectedRules: service.Rules{},
		},
		{
			name:          "file doesn't exit",
			rulesFile:     "not-found.json",
			expectedRules: service.Rules{},
		},
		{
			name:          "file exists and is valid",
			rulesFile:     validRulesFile,
			expectedRules: validRules,
		},
		{
			name:          "reading from 'directory' instead of file",
			rulesFile:     "",
			expectedRules: service.Rules{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := service.NewStorage(
				service.StorageConfig{
					RulesPath: rulesFolder,
				})
			runStorageTest(t, storage, tc.rulesFile, tc.expectedRules)
		})
	}

	t.Run("use cache with a previous read file", func(t *testing.T) {
		// run Find anothertime to test the cache function
		storage := service.NewStorage(
			service.StorageConfig{
				RulesPath: rulesFolder,
			})
		for i := 0; i < 2; i++ {
			runStorageTest(t, storage, validRulesFile, validRules)
		}
	})
}

func runStorageTest(t *testing.T, storage *service.Storage, rulesFile string, expectedRules service.Rules) {
	var rules service.Rules
	data := storage.Find(rulesFile)
	if len(data) == 0 {
		rules = service.Rules{}
	} else {
		err := json.Unmarshal(data, &rules)
		assert.NoError(t, err)
	}
	assert.Equal(t, expectedRules, rules)
}
