package main

type Projects struct {
	Org   string `json:"org"`
	Infos Infos  `json:"projects"`
}

func (p Projects) Equal(otherProjects Projects) bool {
	if p.Org != otherProjects.Org {
		return false
	}

	if len(p.Infos) != len(otherProjects.Infos) {
		return false
	}

	for _, info := range p.Infos {
		found := false

		for _, otherInfo := range otherProjects.Infos {
			if otherInfo.Name == info.Name {
				found = true
				if info == otherInfo {
					break
				} else {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}
