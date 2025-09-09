package server

import (
	// "sync"

	"github.com/GuPoroca/HexTest/pkg/typeDefines"
)

// Global in-memory state (unexported)
var (
	// mu             sync.RWMutex
	currentProject typeDefines.Project
)

//vai dar em nada po relaxa ai

// // Read a copy
// func GetProject() typeDefines.Project {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return currentProject
// }
//
// // Replace the whole project
// func SetProject(p typeDefines.Project) {
// 	mu.Lock()
// 	currentProject = p
// 	mu.Unlock()
// }
//
// // Mutate in place safely
// func UpdateProject(fn func(p *typeDefines.Project)) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	fn(&currentProject)
// }
