package builds

func (a Config) Merge(b Config) Config {
	if b.Image != "" {
		a.Image = b.Image
	}

	if len(a.Params) > 0 {
		newParams := map[string]string{}

		for k, v := range a.Params {
			newParams[k] = v
		}

		for k, v := range b.Params {
			newParams[k] = v
		}

		a.Params = newParams
	} else {
		a.Params = b.Params
	}

	if len(a.Paths) > 0 {
		newPaths := map[string]string{}

		for k, v := range a.Paths {
			newPaths[k] = v
		}

		for k, v := range b.Paths {
			newPaths[k] = v
		}

		a.Paths = newPaths
	} else {
		a.Paths = b.Paths
	}

	if b.Run.Path != "" {
		a.Run = b.Run
	}

	return a
}
