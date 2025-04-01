package quote

const Namespace = "Instagram Quote Generator"
const GenerateTaskQueue = "GENERATE_TASK_QUEUE"

// Quote represents a quote entry in the database
type Quote struct {
	ID        int
	Text      string
	Author    string
	Reference string
	Used      bool
}

type GenerateTextInput struct {
}

type GenerateTextOutput struct {
	Text      string
	Author    string
	Reference string
}

type GenerateImageOutput struct {
	FileName string
}
