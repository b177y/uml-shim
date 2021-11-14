package main

import "github.com/spf13/cobra"

var directory string

var UMLShimCLI = &cobra.Command{
	Use:                   "uml-shim [options] KERNELCMD",
	Short:                 "uml-shim is a tool for running and managing a UserMode Linux instance",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runShim(args[0], args[1:])
	},
}

func init() {
	UMLShimCLI.Flags().StringVarP(&directory, "directory", "d", "", "directory to place connection socket, logs and exit code")
}
