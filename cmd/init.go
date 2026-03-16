package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

var template string

var initCmd = &cobra.Command{
	Use:     "init [framework] [name]",
	Aliases: []string{"new"},
	Short:   "Create a new project scaffold",
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var framework string
		var projectName string

		if len(args) > 0 {
			framework = args[0]
		}
		if len(args) > 1 {
			projectName = args[1]
		}

		// 1. Interactive Form (if missing arguments)
		if framework == "" || projectName == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Choose a framework/language").
						Options(
							huh.NewOption("Go", "go"),
							huh.NewOption("Node.js", "node"),
							huh.NewOption("React (Vite)", "react"),
							huh.NewOption("Next.js", "next"),
							huh.NewOption("Laravel", "laravel"),
							huh.NewOption("Nuxt", "nuxt"),
							huh.NewOption("Vue (Vite)", "vue"),
						).
						Value(&framework),
				).WithHide(framework != ""),
				huh.NewGroup(
					huh.NewInput().
						Title("What is your project name?").
						Value(&projectName).
						Validate(func(s string) error {
							if s == "" {
								return fmt.Errorf("project name cannot be empty")
							}
							return nil
						}),
				).WithHide(projectName != ""),
			)

			if err := form.Run(); err != nil {
				fmt.Println("Aborted.")
				return
			}
		}

		// 2. Scaffolding Logic
		fmt.Printf("🏗️  Scaffolding %s project: %s...\n", framework, projectName)
		scaffoldErr := scaffold(framework, projectName)

		if scaffoldErr != nil {
			fmt.Printf("\n%s Error: %v\n", ui.ErrorStyle.Render("✗"), scaffoldErr)
			return
		}

		fmt.Printf("\n%s Project %s created successfully!\n", ui.SuccessStyle.Render("✓"), projectName)
		fmt.Printf("Next steps:\n  cd %s\n", projectName)
	},
}

func scaffold(framework, name string) error {
	switch framework {
	case "go":
		if err := os.MkdirAll(name, 0755); err != nil {
			return err
		}
		return runCmd(name, "go", "mod", "init", name)

	case "node":
		if err := os.MkdirAll(name, 0755); err != nil {
			return err
		}
		return runCmd(name, "npm", "init", "-y")

	case "react":
		return runCmd("", "sh", "-c", fmt.Sprintf("npm create vite@latest %s -- --template react --no-interactive", name))

	case "next":
		// Added --no-react-compiler and ensured --yes is present
		cmdStr := fmt.Sprintf("npx --yes create-next-app@latest %s --typescript --tailwind --eslint --app --use-npm --no-src-dir --no-react-compiler --import-alias '@/*' --yes", name)
		return runCmd("", "sh", "-c", cmdStr)

	case "laravel":
		return runCmd("", "composer", "create-project", "laravel/laravel", name)

	case "nuxt":
		return runCmd("", "sh", "-c", fmt.Sprintf("npx --yes nuxi@latest init %s --packageManager npm --no-install", name))

	case "vue":
		return runCmd("", "sh", "-c", fmt.Sprintf("npm create vite@latest %s -- --template vue --no-interactive", name))

	default:
		return fmt.Errorf("unsupported framework: %s", framework)
	}
}

func runCmd(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	
	if dir != "" {
		absDir, _ := filepath.Abs(dir)
		cmd.Dir = absDir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// Inherit environment automatically
	return cmd.Run()
}

func init() {
	initCmd.Flags().StringVarP(&template, "template", "t", "", "Project template (go, node, react, next, laravel, nuxt, vue)")
	rootCmd.AddCommand(initCmd)
}
