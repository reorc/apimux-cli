package command

import "testing"

func TestParseSchemaBoundParamsAppliesCompatAliases(t *testing.T) {
	tests := []struct {
		name        string
		capability  string
		args        []string
		wantParam   string
		wantValue   any
		commandPath string
	}{
		{
			name:        "reddit search sort uppercases to canonical lowercase",
			capability:  "reddit.search",
			args:        []string{"--query", "desk setup", "--sort", "TOP"},
			wantParam:   "sort",
			wantValue:   "top",
			commandPath: "apimux reddit search",
		},
		{
			name:        "douyin search videos accepts legacy numeric enums",
			capability:  "douyin.search_videos",
			args:        []string{"--keyword", "desk setup", "--sort-type", "2", "--publish-time", "180"},
			wantParam:   "sort_type",
			wantValue:   "latest",
			commandPath: "apimux douyin search_videos",
		},
		{
			name:        "tiktok search videos accepts legacy numeric enums",
			capability:  "tiktok.search_videos",
			args:        []string{"--keyword", "desk setup", "--sort-by", "2", "--publish-time", "90"},
			wantParam:   "sort_by",
			wantValue:   "date",
			commandPath: "apimux tiktok search_videos",
		},
		{
			name:        "xiaohongshu search notes accepts legacy provider enums",
			capability:  "xiaohongshu.search_notes",
			args:        []string{"--keyword", "desk setup", "--sort-strategy", "popularity_descending", "--note-type", "2"},
			wantParam:   "sort_strategy",
			wantValue:   "likes",
			commandPath: "apimux xiaohongshu search_notes",
		},
		{
			name:        "amazon query aba keywords accepts page-index alias",
			capability:  "amazon.query_aba_keywords",
			args:        []string{"--page-index", "3"},
			wantParam:   "page",
			wantValue:   3,
			commandPath: "apimux amazon query_aba_keywords",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := testSchemas()[tt.capability]
			got, err := parseSchemaBoundParams(tt.args, spec, tt.commandPath)
			if err != nil {
				t.Fatalf("parseSchemaBoundParams() error = %v", err)
			}
			if got[tt.wantParam] != tt.wantValue {
				t.Fatalf("param %s = %#v, want %#v", tt.wantParam, got[tt.wantParam], tt.wantValue)
			}
		})
	}
}

func TestParseSchemaBoundParamsMapsCompatAliasSidecars(t *testing.T) {
	spec := testSchemas()["douyin.search_videos"]
	got, err := parseSchemaBoundParams(
		[]string{"--keyword", "desk setup", "--sort-type", "2", "--publish-time", "180"},
		spec,
		"apimux douyin search_videos",
	)
	if err != nil {
		t.Fatalf("parseSchemaBoundParams() error = %v", err)
	}
	if got["publish_time"] != "6m" {
		t.Fatalf("publish_time = %#v, want %#v", got["publish_time"], "6m")
	}

	spec = testSchemas()["xiaohongshu.search_notes"]
	got, err = parseSchemaBoundParams(
		[]string{"--keyword", "desk setup", "--sort-strategy", "popularity_descending", "--note-type", "2"},
		spec,
		"apimux xiaohongshu search_notes",
	)
	if err != nil {
		t.Fatalf("parseSchemaBoundParams() error = %v", err)
	}
	if got["note_type"] != "video" {
		t.Fatalf("note_type = %#v, want %#v", got["note_type"], "video")
	}
}
