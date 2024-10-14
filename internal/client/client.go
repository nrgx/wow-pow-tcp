package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/nrgx/wow-pow-tcp/internal/model"
	"github.com/rs/zerolog/log"
)

func Run(address string) {
	log.Info().Msgf("connecting to server %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("error dialing server")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Warn().Err(err).Msg("error closing connection to server")
		}
	}()

	for {
		log.Info().Msg("requesting")
		if err := request(conn); err != nil {
			log.Error().Err(err).Msg("error requesting server")
		}
		log.Info().Msg("request processed")
		time.Sleep(3 * time.Second)
	}
}

func request(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	// send challenge
	msg := model.Message{
		Type: model.ChallengeRequest,
	}
	payload, err := json.Marshal(&msg)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling message")
		return err
	}

	log.Info().Msg("sending challenge request")
	if _, err := conn.Write(append(payload, '\n')); err != nil {
		log.Error().Err(err).Msg("error sending challenge request")
		return err
	}

	log.Info().Msg("got challenge")
	// receive challenge
	response, err := reader.ReadBytes('\n')
	if err != nil {
		log.Error().Err(err).Msg("error reading challenge response")
		return err
	}

	log.Info().Msg("parsing challenge")
	var challenge model.Message
	if err := json.Unmarshal(response, &challenge); err != nil {
		log.Error().Err(err).Msg("error unmarshalling challenge")
		return err
	}

	log.Info().Msg("solving challenge")
	// get the solution
	solution, err := challenge.Challenge.Solve()
	if err != nil {
		log.Error().Err(err).Msg("error solving challenge")
		return err
	}

	log.Info().Msg("got the solution")
	// send it to the server
	msg = model.Message{
		Type:      model.SolutionRequest,
		Challenge: solution,
	}
	payload, err = json.Marshal(&msg)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling solution request")
		return err
	}

	log.Info().Msg("sending solution to server")
	if _, err := conn.Write(append(payload, '\n')); err != nil {
		log.Error().Err(err).Msg("error sending solution to server")
		return err
	}

	log.Info().Msg("got results from server")
	// get quote as a reward
	response, err = reader.ReadBytes('\n')
	if err != nil {
		log.Error().Err(err).Msg("error reading response")
		return err
	}

	var result model.Message
	if err := json.Unmarshal(response, &result); err != nil {
		log.Error().Err(err).Msg("error unmarshalling response")
		return err
	}

	if result.Reward == "" {
		log.Error().Err(err).Msg("didn't receive a reward")
		return fmt.Errorf("no reward")
	}

	log.Info().Msgf("got quote from server: %s", result.Reward)
	return nil
}
