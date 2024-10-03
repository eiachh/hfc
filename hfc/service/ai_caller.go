package service

type AiCaller interface {
	callAI([]byte) (*[]byte, error)
}
