package app

type Options struct {
	// Configuration
	ConfigPath string
	ConfigName string

	// Output & execution modes
	JSONOutput bool
	DryRun     bool

	// Lifecycle control
	Force bool

	// Service selection
	BackendOnly  bool
	FrontendOnly bool

	// Service ports
	BackendPort  int
	FrontendPort int

	//Version
	ShowVersion bool
	CheckUpdate bool
}
