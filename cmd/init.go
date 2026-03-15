package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
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
		var scaffoldErr error
		action := func() {
			scaffoldErr = scaffold(framework, projectName)
		}

		_ = spinner.New().
			Title(fmt.Sprintf("Scaffolding %s project: %s...", framework, projectName)).
			Action(action).
			Run()

		if scaffoldErr != nil {
			fmt.Printf("\n%s Error: %v\n", ui.ErrorStyle.Render("✗"), scaffoldErr)
			return
		}

		fmt.Printf("\n%s Project %s created successfully!\n", ui.SuccessStyle.Render("✓"), projectName)
		fmt.Printf("Next steps:\n  cd %s\n", projectName)
	},
}

func scaffold(framework, name string) error {
	// ... existing scaffold logic ...
	switch framework {
	case "go":
		if err := os.MkdirAll(name, 0755); err != nil {
			return err
		}
		cmd := exec.Command("go", "mod", "init", name)
		cmd.Dir = name
		return cmd.Run()

	case "node":
		if err := os.MkdirAll(name, 0755); err != nil {
			return err
		}
		cmd := exec.Command("npm", "init", "-y")
		cmd.Dir = name
		return cmd.Run()

	case "react":
		// Using Vite for React
		return exec.Command("npm", "create", "vite@latest", name, "--", "--template", "react").Run()

	case "next":
		// Non-interactive next app creation
		return exec.Command("npx", "create-next-app@latest", name, "--ts", "--eslint", "--tailwind", "--app", "--src-dir", "import-alias", "@/*", "--use-npm").Run()

	case "laravel":
		return exec.Command("composer", "create-project", "laravel/laravel", name).Run()

	case "nuxt":
		return exec.Command("npx", "nuxi@latest", "init", name).Run()

	case "vue":
		return exec.Command("npm", "create", "vite@latest", name, "--", "--template", "vue").Run()

	default:
		return fmt.Errorf("unsupported framework: %s", framework)
	}
}

func init() {
	initCmd.Flags().StringVarP(&template, "template", "t", "", "Project template (go, node, react, next, laravel, nuxt, vue)")
	rootCmd.AddCommand(initCmd)
}
