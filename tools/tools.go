package tools

// workaround for mockgen
// https://github.com/golang/mock/issues/494#issuecomment-715609711
import (
	_ "go.uber.org/mock/mockgen/model"
)
