package model

type NodeStatus int

const (
	ACTIVE NodeStatus = iota
	DEAD
	NOT_APPLICABLE
)
