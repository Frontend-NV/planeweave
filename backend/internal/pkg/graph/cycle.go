package graph

import "github.com/planeweave/planeweave/internal/domain"

func WouldCreateCycle(edges []domain.GraphEdge, fromID, toID string) bool {
	if fromID == toID {
		return true
	}

	adj := make(map[string][]string)
	for _, e := range edges {
		adj[e.FromTaskID] = append(adj[e.FromTaskID], e.ToTaskID)
	}
	adj[fromID] = append(adj[fromID], toID)

	visited := make(map[string]bool)
	var dfs func(node string) bool
	dfs = func(node string) bool {
		if node == fromID {
			return true
		}
		if visited[node] {
			return false
		}
		visited[node] = true
		for _, next := range adj[node] {
			if dfs(next) {
				return true
			}
		}
		return false
	}
	return dfs(toID)
}
