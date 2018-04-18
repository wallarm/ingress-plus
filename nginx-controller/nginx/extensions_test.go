package nginx

import "testing"

func TestParseLBMethod(t *testing.T) {
	var testsWithValidInput = []struct {
		input    string
		expected string
	}{
		{"", "least_conn"},
		{"least_conn", "least_conn"},
		{"round_robin", ""},
		{"ip_hash", "ip_hash"},
		{"hash $request_id", "hash $request_id"},
	}

	var invalidInput = []string{
		"blabla",
		"least_time header",
		"hash123",
	}

	for _, test := range testsWithValidInput {
		result, err := ParseLBMethod(test.input)
		if err != nil {
			t.Errorf("TestParseLBMethod(%q) returned an error for valid input", test.input)
		}

		if result != test.expected {
			t.Errorf("TestParseLBMethod(%q) returned %q expected %q", test.input, result, test.expected)
		}
	}

	for _, input := range invalidInput {
		_, err := ParseLBMethod(input)
		if err == nil {
			t.Errorf("TestParseLBMethod(%q) does not return an error for invalid input", input)
		}
	}
}

func TestParseLBMethodForPlus(t *testing.T) {
	var testsWithValidInput = []struct {
		input    string
		expected string
	}{
		{"", "least_conn"},
		{"least_conn", "least_conn"},
		{"round_robin", ""},
		{"ip_hash", "ip_hash"},
		{"hash $request_id", "hash $request_id"},
		{"least_time header", "least_time header"},
		{"least_time last_byte", "least_time last_byte"},
		{"least_time header inflight", "least_time header inflight"},
		{"least_time last_byte inflight", "least_time last_byte inflight"},
	}

	var invalidInput = []string{
		"blabla",
		"hash123",
		"least_time header inflight",
	}

	for _, test := range testsWithValidInput {
		result, err := ParseLBMethod(test.input)
		if err != nil {
			t.Errorf("TestParseLBMethod(%q) returned an error for valid input", test.input)
		}

		if result != test.expected {
			t.Errorf("TestParseLBMethod(%q) returned %q expected %q", test.input, result, test.expected)
		}
	}

	for _, input := range invalidInput {
		_, err := ParseLBMethod(input)
		if err == nil {
			t.Errorf("TestParseLBMethod(%q) does not return an error for invalid input", input)
		}
	}
}
