package volumetree

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

type NodeRef struct {
	Node   domain.VolumeNode
	Path   []string
	Depth  int
	Parent *domain.VolumeNode
}

type ConstraintSummary struct {
	Width    *domain.DimensionConstraint
	Height   *domain.DimensionConstraint
	Depth    *domain.DimensionConstraint
	HasFill  bool
	HasRatio bool
	HasFixed bool
}

func AxisToDimension(axis domain.SplitAxis) string {
	switch axis {
	case domain.SplitAxisX:
		return "width"
	case domain.SplitAxisY:
		return "height"
	case domain.SplitAxisZ:
		return "depth"
	default:
		return ""
	}
}

func SumSplitRatios(split domain.VolumeSplit) float64 {
	var sum float64
	for _, r := range split.Ratios {
		sum += r
	}
	return sum
}

func SummarizeConstraints(node domain.VolumeNode) ConstraintSummary {
	all := []*domain.DimensionConstraint{node.Constraints.Width, node.Constraints.Height, node.Constraints.Depth}

	summary := ConstraintSummary{
		Width:  node.Constraints.Width,
		Height: node.Constraints.Height,
		Depth:  node.Constraints.Depth,
	}

	for _, c := range all {
		if c == nil {
			continue
		}
		switch c.Mode {
		case domain.DimensionFill:
			summary.HasFill = true
		case domain.DimensionRatio:
			summary.HasRatio = true
		case domain.DimensionFixed:
			summary.HasFixed = true
		}
	}

	return summary
}

func WalkVolumeTree(root domain.VolumeNode, callback func(ref NodeRef) bool) {
	var visit func(node domain.VolumeNode, path []string, depth int, parent *domain.VolumeNode)
	visit = func(node domain.VolumeNode, path []string, depth int, parent *domain.VolumeNode) {
		ref := NodeRef{Node: node, Path: path, Depth: depth, Parent: parent}
		if !callback(ref) {
			return
		}
		for _, child := range node.Children {
			childPath := append(append([]string{}, path...), child.ID)
			visit(child, childPath, depth+1, &node)
		}
	}
	visit(root, []string{root.ID}, 0, nil)
}

func FindVolumeNodeByID(root domain.VolumeNode, id string) (NodeRef, bool) {
	var found NodeRef
	ok := false

	WalkVolumeTree(root, func(ref NodeRef) bool {
		if ref.Node.ID == id {
			found = ref
			ok = true
			return false
		}
		return true
	})

	return found, ok
}

func FlattenVolumeTree(root domain.VolumeNode) []NodeRef {
	var nodes []NodeRef
	WalkVolumeTree(root, func(ref NodeRef) bool {
		nodes = append(nodes, ref)
		return true
	})
	return nodes
}

func GetTreeDepth(root domain.VolumeNode) int {
	maxDepth := 0
	WalkVolumeTree(root, func(ref NodeRef) bool {
		if ref.Depth > maxDepth {
			maxDepth = ref.Depth
		}
		return true
	})
	return maxDepth
}

func GetNodeCount(root domain.VolumeNode) int {
	count := 0
	WalkVolumeTree(root, func(ref NodeRef) bool {
		count++
		return true
	})
	return count
}

func GetLeafNodes(root domain.VolumeNode) []domain.VolumeNode {
	var leaves []domain.VolumeNode
	WalkVolumeTree(root, func(ref NodeRef) bool {
		if len(ref.Node.Children) == 0 {
			leaves = append(leaves, ref.Node)
		}
		return true
	})
	return leaves
}

func CollectFeatures(root domain.VolumeNode) []domain.Feature {
	var features []domain.Feature
	WalkVolumeTree(root, func(ref NodeRef) bool {
		features = append(features, ref.Node.Features...)
		return true
	})
	return features
}

func CollectFronts(root domain.VolumeNode) []domain.Front {
	var fronts []domain.Front
	WalkVolumeTree(root, func(ref NodeRef) bool {
		fronts = append(fronts, ref.Node.Fronts...)
		return true
	})
	return fronts
}

func CountFeaturesByType(root domain.VolumeNode) map[string]int {
	counts := map[string]int{}
	for _, feature := range CollectFeatures(root) {
		counts[feature.Type]++
	}
	return counts
}
