package db

import "github.com/concourse/atc"

type ConfigDBWithDefaults struct {
	ConfigDB
}

func (configDB ConfigDBWithDefaults) GetConfig() (atc.Config, error) {
	config, err := configDB.getUnderlyingConfig()
	if err != nil {
		return atc.Config{}, err
	}

	triggerDefault := true

	for _, job := range config.Jobs {
		for i, input := range job.Inputs {
			if input.Name == "" {
				job.Inputs[i].Name = input.Resource
			}

			if input.Trigger == nil {
				job.Inputs[i].Trigger = &triggerDefault
			}
		}

		for i, output := range job.Outputs {
			if output.PerformOn == nil {
				job.Outputs[i].PerformOn = []atc.OutputCondition{"success"}
			}
		}
	}

	return config, nil
}

func (configDB ConfigDBWithDefaults) JobIsPublic(jobName string) (bool, error) {
	config, err := configDB.getUnderlyingConfig()
	if err != nil {
		return false, err
	}

	return config.JobIsPublic(jobName)
}

func (configDB ConfigDBWithDefaults) getUnderlyingConfig() (atc.Config, error) {
	return configDB.ConfigDB.GetConfig()
}
