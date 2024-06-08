package scheduler

import messagepackage "proj3/messagePackage"

// Run the correct version based on the Mode field of the configuration value
func Schedule(config messagepackage.Config) {
	if config.Mode == "s" || config.Mode == "p" || config.Mode == "chunk" {
		CreateImagesTaskQueueAndRun(config)
	} else {
		panic("Invalid scheduling scheme given.")
	}
}
