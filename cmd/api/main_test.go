package main

import "testing"

func TestStartupURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		host string
		port string
		want string
	}{
		{
			name: "replace wildcard host with localhost",
			host: "0.0.0.0",
			port: "9100",
			want: "http://localhost:9100",
		},
		{
			name: "use explicit ipv4 host",
			host: "127.0.0.1",
			port: "9100",
			want: "http://127.0.0.1:9100",
		},
		{
			name: "normalize ipv6 host",
			host: "::1",
			port: "9100",
			want: "http://[::1]:9100",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := startupURL(testCase.host, testCase.port)
			if got != testCase.want {
				t.Fatalf("startupURL(%q, %q) = %q, want %q", testCase.host, testCase.port, got, testCase.want)
			}
		})
	}
}
