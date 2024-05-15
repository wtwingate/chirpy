package main

import "testing"

func TestFilterProfanity(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "message with no profanity",
			input: "I'm a nice, polite message",
			want:  "I'm a nice, polite message",
		},
		{
			name:  "message with profanity",
			input: "What the actual kerfuffle is this?!",
			want:  "What the actual **** is this?!",
		},
		{
			name:  "message with multiple profanities",
			input: "I'm gonna kick your fornax you son of a sharbert",
			want:  "I'm gonna kick your **** you son of a ****",
		},
		{
			name:  "message with capitalized profanity",
			input: "FORNAX YOU!",
			want:  "**** YOU!",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := filterProfanity(c.input)
			if result != c.want {
				t.Errorf("got %v, want %v", result, c.want)
			}
		})
	}
}
