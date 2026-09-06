package internal

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

type Player struct {
	Id    string
	State struct {
		Live                    bool
		Paused                  bool
		InterpolatedCurrentTime float64
	}
	Type string
}

func APIQuery(client *graphql.Client, query any, variables map[string]interface{}) (any, error) {
	err := client.Query(context.Background(), query, variables)

	if err != nil {
		return nil, err
	}

	return query, nil
}

func APIMutation(client *graphql.Client, mutation any, variables map[string]interface{}) (any, error) {
	err := client.Mutate(context.Background(), mutation, variables)

	if err != nil {
		return nil, err
	}

	return mutation, nil
}

func APISubscribe(client *graphql.SubscriptionClient, query any, variables map[string]interface{}, handler func([]byte, error) error) (string, error) {
	subscriptionId, err := client.Subscribe(&query, variables, handler)

	if err != nil {
		return "", err
	}

	return subscriptionId, err
}

func GetPlayers(client *graphql.Client) ([]*Player, error) {
	var query struct {
		Players []Player
	}

	_, err := APIQuery(client, &query, nil)

	var players []*Player

	for i := range query.Players {
		players = append(players, &query.Players[i])
	}

	return players, err
}

func GetPlayer(client *graphql.Client, playerId string) (Player, error) {
	var query struct {
		Player Player `graphql:"player(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(playerId),
	}

	_, err := APIQuery(client, &query, variables)

	return query.Player, err
}

func GetAdditionalPlayers(client *graphql.Client) ([]*Player, error) {
	players, err := GetPlayers(client)

	if err != nil {
		return nil, err
	}

	var additional []*Player

	for _, player := range players {
		if player.Type == "ADDITIONAL" {
			additional = append(additional, player)
		}
	}

	return additional, nil
}

func PausePlayer(client *graphql.Client, playerId string, paused bool) error {
	var mutation struct {
		PlayerSetPaused bool `graphql:"playerSetPaused(id: $id, paused: $paused)"`
	}

	variables := map[string]interface{}{
		"id":     graphql.ID(playerId),
		"paused": paused,
	}

	_, err := APIMutation(client, &mutation, variables)

	return err
}

func PausePlayers(client *graphql.Client, players []*Player, paused bool) error {
	for _, player := range players {
		err := PausePlayer(client, player.Id, paused)

		if err != nil {
			return err
		}

		player.State.Paused = paused
	}

	return nil
}

func SyncPlayersToPlayer(client *graphql.Client, playerId string) error {
	var mutation struct {
		PlayerSync bool `graphql:"playerSync(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(playerId),
	}

	_, err := APIMutation(client, mutation, variables)

	return err
}

func SeekPlayerTo(client *graphql.Client, playerId string, time float64) error {
	var mutation struct {
		PlayerSeekTo bool `graphql:"playerSeekTo(id: $id, absolute: $time)"`
	}

	variables := map[string]interface{}{
		"id":   graphql.ID(playerId),
		"time": time,
	}

	_, err := APIMutation(client, mutation, variables)

	return err
}
