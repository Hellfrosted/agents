package status

const targetSetThreshold = 8

type targetFilter struct {
	list []string
	set  map[string]struct{}
}

func newTargetFilter(targets []string) targetFilter {
	if len(targets) == 0 {
		return targetFilter{}
	}
	if len(targets) <= targetSetThreshold {
		return targetFilter{list: targets}
	}
	selected := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		selected[target] = struct{}{}
	}
	return targetFilter{set: selected}
}

func (targets targetFilter) selected(name string) bool {
	if len(targets.list) == 0 && len(targets.set) == 0 {
		return true
	}
	if len(targets.set) > 0 {
		_, ok := targets.set[name]
		return ok
	}
	for _, target := range targets.list {
		if target == name {
			return true
		}
	}
	return false
}
