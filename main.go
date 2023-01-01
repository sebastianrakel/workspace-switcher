package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	version string = "dirty"
	config  *WorkspaceSwitcherConfiguration
	cfgFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "workspace-switcher",
		Short: "Workspace Switcher",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints version",
		Run:   printVersion,
		Aliases: []string{
			"v",
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List workspaces",
		Run:   listWorkspaces,
		Aliases: []string{
			"l",
		},
	}

	var listAliasesCmd = &cobra.Command{
		Use:   "aliases",
		Short: "List aliases,",
		Run:   listAliases,
	}

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply workspace",
		Run:   applyWorkspaceCmd,
	}

	var rofiCmd = &cobra.Command{
		Use:   "rofi",
		Short: "Open Rofi",
		Run:   showRofi,
	}

	defaultConfigPath := path.Join(os.Getenv("HOME"), ".config", "workspace-switcher", "config.yaml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfigPath, "Path to configuration file")
	rootCmd.AddCommand(versionCmd, listCmd, listAliasesCmd, applyCmd, rofiCmd)

	cfg, err := loadConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	config = cfg

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Workspace Switcher\nVersion: %s\n", version)
}

func listWorkspaces(cmd *cobra.Command, args []string) {
	if config.Workspaces == nil {
		log.Println("no workspaces found")
		return
	}

	for key := range config.Workspaces {
		fmt.Printf("%s\n", key)
	}
}

func listAliases(cmd *cobra.Command, args []string) {
	if config.Aliases == nil {
		log.Println("no aliases found")
		return
	}

	for _, name := range config.GetWorkspaceNames() {
		fmt.Printf("%s\n", name)
	}
}

func applyWorkspaceCmd(cmd *cobra.Command, args []string) {
	workspaceName := args[0]

	applyWorkspace(workspaceName)
}

func GetDisplayName(nameOrAlias string) string {
	if alias, exists := config.Aliases[nameOrAlias]; exists {
		return alias
	}

	return nameOrAlias
}

func (w *Workspace) ExecuteDisplayCommand() {
	var displayBlocks []string

	for key, display := range w.Displays {
		displayBlocks = append(displayBlocks, display.GetDisplayCommandBlock(key)...)
	}

	displayBlocks = append(displayBlocks, "--verbose")

	execute("/usr/bin/xrandr", displayBlocks, nil)
}

func (d *Display) GetDisplayCommandBlock(displayName string) []string {
	parts := []string{}

	parts = append(parts, "--output", GetDisplayName(displayName))
	if d.Primary {
		parts = append(parts, "--primary")
	}

	if d.Rotation != "" {
		parts = append(parts, "--rotation", d.Rotation)
	}

	for _, order := range d.Order {
		parts = append(parts, fmt.Sprintf("--%s", order.Position), GetDisplayName(order.Display))
	}

	if d.Resolution == "" {
		parts = append(parts, "--auto")
	}

	return parts
}

func deactivateDisplays() {
	args := []string{}
	for _, output := range getOutputs() {
		args = append(args, "--output", output, "--off")
	}

	execute("/usr/bin/xrandr", args, nil)
}

func getOutputs() []string {
	cmd := exec.Command("xrandr")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
		panic(err)
	}

	lines := strings.Split(string(out), "\n")
	outputs := []string{}

	for _, line := range lines {
		if strings.Contains(line, "connected") {
			outputs = append(outputs, strings.Split(line, " ")[0])
		}
	}

	return outputs
}

func showRofi(cmd *cobra.Command, args []string) {
	workspace, err := openRofi(config.GetWorkspaceNames())

	if err == nil && workspace != "" {
		applyWorkspace(workspace)
	}
}

func applyWorkspace(workspaceName string) {
	workspace := config.Workspaces[workspaceName]

	deactivateDisplays()

	workspace.ExecuteDisplayCommand()

	for _, hook := range config.Hooks.Activate {
		executeString(hook)
	}

	for _, hook := range workspace.Hooks.Activate {
		executeString(hook)
	}
}
