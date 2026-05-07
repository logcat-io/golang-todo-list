package todo

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func (t *Todo) IsDone() bool {
	return t.Done
}

func (t *Todo) MarkDone() {
	t.Done = true
}
