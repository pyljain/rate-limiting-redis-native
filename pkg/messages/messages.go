package messages

type Request struct {
	CSI    string `json:"csi"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
