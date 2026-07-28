package graph_test

import (
	"testing"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/pkg/graph"
)

func TestWouldCreateCycle(t *testing.T) {
	edges := []domain.GraphEdge{
		{FromTaskID: "a", ToTaskID: "b"},
		{FromTaskID: "b", ToTaskID: "c"},
	}

	tests := []struct {
		name   string
		from   string
		to     string
		expect bool
	}{
		{"cycle c to a", "c", "a", true},
		{"self loop", "a", "a", true},
		{"safe d to a", "d", "a", false},
		{"safe c to d", "c", "d", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.WouldCreateCycle(edges, tt.from, tt.to)
			if got != tt.expect {
				t.Fatalf("WouldCreateCycle(%s,%s)=%v want %v", tt.from, tt.to, got, tt.expect)
			}
		})
	}
}
