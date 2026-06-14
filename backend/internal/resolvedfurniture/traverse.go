package resolvedfurniture

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

type VolumeRef struct {
	Volume domain.ResolvedVolume
	Path   []string
	Depth  int
	Parent *domain.ResolvedVolume
}

func WalkResolvedTree(root domain.ResolvedVolume, callback func(ref VolumeRef) bool) {
	var visit func(volume domain.ResolvedVolume, path []string, depth int, parent *domain.ResolvedVolume)
	visit = func(volume domain.ResolvedVolume, path []string, depth int, parent *domain.ResolvedVolume) {
		ref := VolumeRef{Volume: volume, Path: path, Depth: depth, Parent: parent}
		if !callback(ref) {
			return
		}
		for _, child := range volume.Children {
			childPath := append(append([]string{}, path...), child.ID)
			visit(child, childPath, depth+1, &volume)
		}
	}
	visit(root, []string{root.ID}, 0, nil)
}

func FindResolvedVolumeByID(root domain.ResolvedVolume, id string) (VolumeRef, bool) {
	var found VolumeRef
	ok := false

	WalkResolvedTree(root, func(ref VolumeRef) bool {
		if ref.Volume.ID == id {
			found = ref
			ok = true
			return false
		}
		return true
	})

	return found, ok
}

func FlattenResolvedTree(root domain.ResolvedVolume) []VolumeRef {
	var nodes []VolumeRef
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		nodes = append(nodes, ref)
		return true
	})
	return nodes
}

func GetTreeDepth(root domain.ResolvedVolume) int {
	maxDepth := 0
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		if ref.Depth > maxDepth {
			maxDepth = ref.Depth
		}
		return true
	})
	return maxDepth
}

func GetNodeCount(root domain.ResolvedVolume) int {
	count := 0
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		count++
		return true
	})
	return count
}

func GetLeafVolumes(root domain.ResolvedVolume) []domain.ResolvedVolume {
	var leaves []domain.ResolvedVolume
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		if len(ref.Volume.Children) == 0 {
			leaves = append(leaves, ref.Volume)
		}
		return true
	})
	return leaves
}

func CollectFeatures(root domain.ResolvedVolume) []domain.ResolvedFeature {
	var features []domain.ResolvedFeature
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		features = append(features, ref.Volume.Features...)
		return true
	})
	return features
}

func CollectFronts(root domain.ResolvedVolume) []domain.ResolvedFront {
	var fronts []domain.ResolvedFront
	WalkResolvedTree(root, func(ref VolumeRef) bool {
		fronts = append(fronts, ref.Volume.Fronts...)
		return true
	})
	return fronts
}
