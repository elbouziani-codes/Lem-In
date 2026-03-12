package main

import (
	"os"
	"testing"
)

// Test edge cases
func TestInvalidFormat(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Missing start",
			content: "3\n2 1 2\n##end\n0 2 2\n2-0\n",
		},
		{
			name:    "Missing ants",
			content: "##start\n1 1 1\n##end\n2 2 2\n1-2\n",
		},
		{
			name:    "No path possible",
			content: "3\n##start\n1 1 1\n2 2 2\n3 3 3\n##end\n4 4 4\n1-2\n3-4\n",
		},
		{
			name:    "Duplicate room",
			content: "3\n##start\n1 1 1\n1 2 2\n##end\n2 3 3\n1-2",
		},
		{
			name:    "Link to self",
			content: "3\n##start\n1 1 1\n##end\n2 2 2\n1-1\n1-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "lemin_test_*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(file.Name())

			file.WriteString(tt.content)
			file.Close()

			graph, _, err := ParseInput(file.Name())
			if err == nil {
                // Check if paths fails
                _, pathErr := FindOptimalPaths(graph)
                if pathErr == nil {
				    t.Errorf("expected error for case: %s", tt.name)
                }
			}
		})
	}
}

func TestValidGraph(t *testing.T) {
	content := `3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5`

	file, err := os.CreateTemp("", "lemin_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(content)
	file.Close()

	graph, _, err := ParseInput(file.Name())
	if err != nil {
		t.Fatalf("unexpected error parsing valid graph: %v", err)
	}

	paths, err := FindOptimalPaths(graph)
	if err != nil {
		t.Fatalf("unexpected error finding optimal paths: %v", err)
	}

	if len(paths) != 2 {
		t.Errorf("expected 2 disjoint paths, got %d", len(paths))
	}
}
