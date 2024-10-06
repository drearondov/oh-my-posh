package segments

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime"
)

type Taskwarrior struct {
	props properties.Properties
	env   runtime.Environment

	TasksInfo
}

type TasksInfo struct {
	Context      string
	Projects     string
	ActiveTasks  int
	PendingTasks string
	TotalTasks   string
	Tags         string
	Error        string
}

const (
	GetContext      properties.Property = "get_context"
	GetActiveTasks  properties.Property = "get_active_tasks"
	GetProjects     properties.Property = "get_projects"
	GetPendingTasks properties.Property = "get_pending_tasks"
	GetTotalTasks   properties.Property = "get_total_tasks"
	GetTags         properties.Property = "get_tags"
	Separator       properties.Property = "separator"
)

func (t *Taskwarrior) Template() string {
	return "{{ .Context }}  {{ .Projects }}{{ if .ActiveTasks }}  \uf058  {{ .ActiveTasks }}{{ end }}{{ if or (.PendingTasks) (.TotalTasks)}}  \uf03a  {{ if .PendingTasks }}{{ .PendingTasks }}/{{ end }}{{ if .TotalTasks }}{{ .TotalTasks }}{{ end }}{{ end }} "
}

func (t *Taskwarrior) Enabled() bool {
	return true
}

func (t *Taskwarrior) Init(props properties.Properties, env runtime.Environment) {
	t.props = props
	t.env = env

	if t.hasTaskrc() {
		t.loadTasksInfo()
	}
}

func (t *Taskwarrior) hasTaskrc() bool {
	home := t.env.Home()
	taskrc := filepath.Join(home, ".taskrc")

	_, err := os.Stat(taskrc)

	switch {
	case err == nil:
		return true
	case os.IsNotExist(err):
		taskrcEnv := t.env.Getenv("TASKRC")
		if len(taskrcEnv) == 0 {
			t.Error = "No .taskrc found"
			return false
		}
	default:
		fmt.Println("Error occurred:", err)
		t.Error = "Error loading .taskrc"
	}

	return true
}

func (t *Taskwarrior) loadTasksInfo() {

	// Getting the Context
	if t.props.GetBool(GetContext, true) {
		contextOutput, err := t.env.RunCommand("task", "_get", "rc.context")
		if err == nil {
			t.Context = fmt.Sprintf("\uf187  %s", contextOutput)
		}

		if contextOutput == "" {
			t.Context = "\uf187  No context"
		}
	}

	// Getting the Projects
	if t.props.GetBool(GetProjects, true) {
		projectsOutput, err := t.env.RunCommand("task", "projects")
		if err == nil {
			lines := strings.Split(projectsOutput, "\n")
			var projects []string

			for _, line := range lines[1 : len(lines)-2] { // Skip header and summary lines
				var project string
				var tasks int

				if _, err := fmt.Sscanf(line, "%s %d", &project, &tasks); err != nil {
					continue // Skip this line if there's an error
				}
				projects = append(projects, project)
			}

			// Join the project names with a comma and print them
			projectsList := strings.Join(projects, ", ")

			t.Projects = fmt.Sprintf("\uf502  %s", projectsList)
		}
	}

	// Getting number of Active Tasks
	if t.props.GetBool(GetActiveTasks, true) {
		activeTasks, err := t.env.RunCommand("task", "active")

		scanner := bufio.NewScanner(strings.NewReader(activeTasks))

		if err == nil {
			var lastLine string

			for scanner.Scan() {
				lastLine = scanner.Text() // Store each line
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading output:", err)
				return
			}

			var firstNumber int
			_, err = fmt.Sscanf(lastLine, "%d", &firstNumber)
			if err != nil {
				fmt.Println("Error parsing number:", err)
				return
			}

			t.ActiveTasks = firstNumber
		}

	}

	// Getting Pending tasks
	if t.props.GetBool(GetPendingTasks, true) {
		pendingOutput, err := t.env.RunCommand("task", "status:pending", "count")
		if err == nil {
			t.PendingTasks = pendingOutput
		}
	}

	// Getting Total tasks
	if t.props.GetBool(GetTotalTasks, true) {
		totalOutput, err := t.env.RunCommand("task", "count")
		if err == nil {
			t.TotalTasks = totalOutput
		}
	}
}
