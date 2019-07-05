package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	n, a int
	s string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("n is %v\n", n)
		fmt.Printf("a is %v\n", a)
		fmt.Printf("s is %v\n", s)
		q, _ := cmd.Flags().GetBool("q")
		fmt.Printf("q is %v\n", q)
		bbb, _ := cmd.Flags().GetInt("bbb")
		fmt.Printf("bbb is %v\n", bbb)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().IntVar(&n, "intf", 0, "Set Int")
	newCmd.Flags().StringVar(&s, "stringf", "sss", "Set String")
	newCmd.Flags().Bool("q", false, "Set Bool")

	newCmd.Flags().IntVarP(&a, "aaa", "a", 1, "Set A")
	newCmd.Flags().IntP("bbb", "b", -1, "Set B")
}
