package progressbar

import "github.com/schollz/progressbar/v3"

// this is used to consolidate the configuration options for a progress bar

// BarOptions is a helper struct to consolidate the options and set defaults
type BarOptions struct {
	Max          int
	EnableColors bool
	ShowBytes    bool
	Width        int
	Description  string
}

// GetProgressBar gets a progress bar with the ability to set overrides
func GetProgressBar(options *BarOptions) *progressbar.ProgressBar {
	if options == nil {
		options = &BarOptions{}
	}
	if options.Width == 0 {
		options.Width = 100
	}
	if options.Description == "" {
		options.Description = "Working..."
	}

	return progressbar.NewOptions(options.Max,
		progressbar.OptionEnableColorCodes(options.EnableColors),
		progressbar.OptionShowBytes(options.ShowBytes),
		progressbar.OptionSetWidth(options.Width),
		progressbar.OptionSetDescription(options.Description),
	)
}
