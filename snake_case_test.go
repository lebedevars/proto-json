package main

import "testing"

type SnakeTest struct {
	input  string
	output string
}

var tests = []SnakeTest{
	{"a", "a"},
	{"snake", "snake"},
	{"A", "a"},
	{"ID", "id"},
	{"MOTD", "motd"},
	{"Snake", "snake"},
	{"SnakeTest", "snake_test"},
	{"SnakeID", "snake_id"},
	{"SnakeIDGoogle", "snake_id_google"},
	{"LinuxMOTD", "linux_motd"},
	{"OMGWTFBBQ", "omgwtfbbq"},
	{"omg_wtf_bbq", "omg_wtf_bbq"},
}

func TestToSnake(t *testing.T) {
	for _, test := range tests {
		if toSnakeCase(test.input) != test.output {
			t.Errorf(`ToSnake("%s"), wanted "%s", got \%s"`, test.input, test.output, toSnakeCase(test.input))
		}
	}
}
