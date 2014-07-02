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

	if len(b.Inputs) > 0 {
		newInputs := make([]Input, len(a.Inputs))
		copy(newInputs, a.Inputs)

		for _, bi := range b.Inputs {
			for i, ai := range newInputs {
				if ai.Name == bi.Name {
					ai.DestinationPath = bi.DestinationPath
					newInputs[i] = ai
				}
			}
		}

		a.Inputs = newInputs
	}

	if b.Run.Path != "" {
		a.Run = b.Run
	}

	return a
}
