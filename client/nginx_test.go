package client

import "testing"

const validStabStats = "Active connections: 1457 \nserver accepts handled requests\n 6717066 6717066 65844359 \nReading: 1 Writing: 8 Waiting: 1448 \n"

func TestParseStubStatsValidInput(t *testing.T) {
	var tests = []struct {
		input          []byte
		expectedResult StubStats
		expectedError  bool
	}{
		{
			input: []byte(validStabStats),
			expectedResult: StubStats{
				Connections: StubConnections{
					Active:   1457,
					Accepted: 6717066,
					Handled:  6717066,
					Reading:  1,
					Writing:  8,
					Waiting:  1448,
				},
				Requests: 65844359,
			},
			expectedError: false,
		},
		{
			input:         []byte("invalid-stats"),
			expectedError: true,
		},
	}

	for _, test := range tests {
		var result StubStats

		err := parseStubStats(test.input, &result)

		if err != nil && !test.expectedError {
			t.Errorf("parseStubStats() returned error for valid input %q: %v", string(test.input), err)
		}

		if !test.expectedError && test.expectedResult != result {
			t.Errorf("parseStubStats() result %v != expected %v for input %q", result, test.expectedResult, test.input)
		}
	}
}
