package server

import (
	"bufio"
	"encoding/json"
	"math/rand/v2"
	"net"

	"github.com/nrgx/wow-pow-tcp/internal/model"
	"github.com/nrgx/wow-pow-tcp/internal/pow"
	"github.com/rs/zerolog/log"
)

// e.g. quotes or words of wisdom
var rewards = []string{
	`Make it work, make it right, make it fast. – Kent Beck`,
	`Programmers seem to be changing the world. It would be a relief, for them and for all of us, if they knew something about it. – Ellen Ullman`,
	`Most good programmers do programming not because they expect to get paid or get adulation by the public, but because it is fun to program. – Linus Torvalds`,
	`Programming is learned by writing programs. ― Brian Kernighan`,
	`Computers are fast; developers keep them slow. – Anonymous`,
}

// Start starts server
func Start(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("error listening to host:port")

	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Warn().Err(err).Msg("error closing listener")
		}
	}()
	log.Info().Msgf("started server at %s", address)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Err(err).Msg("error accepting connection; skipping")
			continue
		}
		go handle(conn)
	}
}

// handles incoming requests
func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Error().Err(err).Msg("error reading from connection")
			return
		}

		var msg model.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Error().Err(err).Msg("error unmarshalling payload")
			return
		}

		result := process(msg, conn.RemoteAddr().String())

		payload, err := json.Marshal(&result)
		if err != nil {
			log.Error().Err(err).Msg("error preparing result")
			return
		}

		if _, err := conn.Write(append(payload, '\n')); err != nil {
			log.Error().Err(err).Msg("error sending payload")
			return
		}
	}
}

// process processes requests by their types
func process(msg model.Message, client string) model.Message {
	var result model.Message
	switch msg.Type {
	case model.ChallengeRequest:
		log.Info().Msgf("client %s requested challenge", client)
		// give challenge
		challenge := pow.NewHashcash(client)
		result.Type = model.ChallengeResponse
		result.Challenge = challenge
		return result
	case model.SolutionRequest:
		log.Info().Msgf("client %s sent a solution", client)
		// client came up with solution
		solution := msg.Challenge
		// verify solution by solving on server
		verified, err := msg.Challenge.Solve()
		if err != nil {
			log.Error().Err(err).Msg("error verifying solution")
			result.Type = model.InvalidSolution
			return result
		}

		// client didn't find nonce
		if solution.Counter != verified.Counter {
			result.Type = model.InvalidSolution
			return result
		}

		// client found nonce
		result.Type = model.SolutionResponse
		result.Reward = rewards[rand.IntN(len(rewards))]
		return result
	default:
		// unknown request type received
		result.Type = model.InvalidRequest
		return result
	}
}
