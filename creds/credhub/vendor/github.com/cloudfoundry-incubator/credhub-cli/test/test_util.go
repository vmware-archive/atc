package test

import "os"

func CleanEnv() {
	os.Unsetenv("CREDHUB_SECRET")
	os.Unsetenv("CREDHUB_CLIENT")
}
