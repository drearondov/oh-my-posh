package segments

import (
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"
)

func TestTaskwarriorTemplate(t *testing.T) {
	type ResolveSymlink struct {
		Path string
		Err  error
	}
	cases := []struct {
		Case         string
		Taskrc       ResolveSymlink
		Expected     string
		Context      string
		Projects     string
		ActiveTasks  string
		PendingTasks string
		TotalTasks   string
	}{
		{
			Case:     "No active context, active task",
			Taskrc:   ResolveSymlink{Path: "/home/vscode/.taskrc", Err: nil},
			Expected: "  No context    hobby, personal, photography, work    3    15/26",
			Context: `Name     Definition                           Active
hobby    project:hobby or project:photography no
personal project:personal                     no
work     project:work                         no

Use 'task context none' to unset the current context.`,
			Projects: `Project     Tasks
hobby           1
personal        2
photography    10
work            2

4 projects (15 tasks)`,
			ActiveTasks: `ID Started    Active Age  P Project  Description
 3 2024-09-30   3d   3d   M personal Workout routine
 1 2024-09-30   3d   3d   H work     Prepare Q4 objectives
 2 2024-09-30   3d   3d   M work     Organize project files

3 tasks`,
			PendingTasks: "15",
			TotalTasks:   "26",
		},
		{
			Case:     "Hobby context active, no active task",
			Expected: "  hobby    hobby, photography    11/13",
			Context:  "hobby",
			Projects: `Project     Tasks
hobby           1
photography    10

2 projects (11 tasks)`,
			ActiveTasks:  "No matches.",
			PendingTasks: "11",
			TotalTasks:   "13",
		},
	}

	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("RunCommand", "task", "_get", "rc.context").Return(tc.Context)
		env.On("RunCommand", "task", "projects").Return(tc.Projects)
		env.On("RunCommand", "task", "active").Return(tc.ActiveTasks)
		env.On("RunCommand", "task", "status:pending", "count").Return(tc.PendingTasks)
		env.On("RunCommand", "task", "count").Return(tc.TotalTasks)
	}
}
