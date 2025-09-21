package tui

import "github.com/GuPoroca/HexTest/pkg/typeDefines"

// Send this to the program when your runner has a new snapshot of the tree.
type TreeUpdateMsg struct {
	Project typeDefines.Project
}

// From a leaf check
type CheckUpdateMsg struct {
	SuiteName  string
	TestName   string
	AssertName string
	CheckName  string
	NewStatus  int
}

