package gitops

import "github.com/spf13/cobra"

// TextFilter filters text
type TextFilter struct {
	Includes []string
	Excludes []string
}

func (o *TextFilter) AddFlags(cmd *cobra.Command, optionPrefix, name string) {
	cmd.Flags().StringSliceVar(&o.Includes, optionPrefix+"-include", []string{}, "text strings in the "+name+" to be included when synchronising")
	cmd.Flags().StringSliceVar(&o.Includes, optionPrefix+"-exclude", []string{}, "text strings in the "+name+" to be excluded when synchronising")
}
