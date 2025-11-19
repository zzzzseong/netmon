package commands

import "fmt"

// Command는 명령어 인터페이스입니다
type Command interface {
	// Name은 명령어 이름을 반환합니다
	Name() string
	// Description은 명령어 설명을 반환합니다
	Description() string
	// Usage는 명령어 사용법을 반환합니다
	Usage() string
	// Execute는 명령어를 실행합니다
	Execute(args []string) error
}

// Registry는 명령어 레지스트리입니다
type Registry struct {
	commands map[string]Command
	order    []Command // 등록 순서 유지
}

// NewRegistry는 새로운 명령어 레지스트리를 생성합니다
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
		order:    make([]Command, 0),
	}
}

// Register는 명령어를 등록합니다
func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
	r.order = append(r.order, cmd)
}

// Get은 명령어를 가져옵니다
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// List는 등록된 모든 명령어를 등록 순서대로 반환합니다
func (r *Registry) List() []Command {
	return r.order
}


// Execute는 명령어를 실행합니다
func (r *Registry) Execute(name string, args []string) error {
	cmd, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
	}
	return cmd.Execute(args)
}

