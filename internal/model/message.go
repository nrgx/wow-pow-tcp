package model

import "github.com/nrgx/wow-pow-tcp/internal/pow"

type Type int

const (
	ChallengeRequest Type = iota
	ChallengeResponse
	SolutionRequest
	SolutionResponse
	InvalidRequest
	InvalidSolution
)

type Message struct {
	Type      Type         `json:"type"`
	Challenge pow.Hashcash `json:"data,omitempty"`
	Reward    string       `json:"reward,omitempty"`
}
