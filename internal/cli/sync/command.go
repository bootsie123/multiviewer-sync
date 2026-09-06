package sync

import (
	"fmt"
	"math"
	"multiviewer-sync/internal"
	"net"
	"net/netip"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var (
	server  string
	client  string
	drift   int
	refresh int
	verbose bool
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Syncs the local MultiViewer client to the specified main instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := ParseHost(server)

			if err != nil {
				return fmt.Errorf("The given server \"%s\" is not valid: %w\n", server, err)
			}

			_, err = ParseHost(client)

			if err != nil {
				return fmt.Errorf("The given client \"%s\" is not valid: %w\n", client, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting MultiViewer sync...")

			verbose, _ = cmd.Flags().GetBool("verbose")

			client := graphql.NewClient(fmt.Sprintf("http://%s/api/graphql", client), nil)

			fmt.Println("Connected to client MultiViewer instance!")

			server := graphql.NewClient(fmt.Sprintf("http://%s/api/graphql", server), nil)

			fmt.Println("Connected to server MultiViewer instance!")

			allPaused := false

			for {
				serverPlayers, _ := internal.GetPlayers(server)
				clientPlayers, _ := internal.GetPlayers(client)

				if verbose {
					fmt.Printf("Found %d server players and %d client players\n", len(serverPlayers), len(clientPlayers))
				}

				var serverMainPlayers []*internal.Player

				for _, player := range serverPlayers {
					if player.Type == "ADDITIONAL" {
						serverMainPlayers = append(serverMainPlayers, player)
					}
				}

				for _, player := range serverMainPlayers {
					playerPaused := player.State.Paused

					if playerPaused != allPaused {
						if verbose {
							fmt.Printf("Syncing player pause state to %t\n", playerPaused)
						}

						err := internal.PausePlayers(server, serverPlayers, playerPaused)

						if err != nil {
							fmt.Println(err)
						}

						err = internal.PausePlayers(client, clientPlayers, playerPaused)

						if err != nil {
							fmt.Println(err)
						}

						allPaused = playerPaused

						break
					}
				}

				if len(serverMainPlayers) > 0 && len(clientPlayers) > 0 {
					serverPlayer := serverMainPlayers[0]
					clientPlayer := clientPlayers[1]

					serverTime := serverPlayer.State.InterpolatedCurrentTime

					if math.Abs(serverTime-clientPlayer.State.InterpolatedCurrentTime) > float64(drift) {
						if verbose {
							fmt.Printf("Syncing player time state to %0.4f\n", serverTime)
						}

						internal.SeekPlayerTo(client, clientPlayer.Id, serverPlayer.State.InterpolatedCurrentTime)
						internal.SyncPlayersToPlayer(client, clientPlayer.Id)
					}
				}

				time.Sleep(time.Millisecond * time.Duration(refresh))
			}
		},
	}

	cmd.Flags().StringVarP(&server, "server", "s", "localhost:10101", "IP and port of the server - the main MultiViewer instance")
	cmd.Flags().StringVarP(&client, "client", "c", "localhost:10101", "IP and port of the client - the MultiViewer instance to sync to")
	cmd.Flags().IntVarP(&drift, "drift", "d", 2, "Maximum amount of time drift allowed in seconds")
	cmd.Flags().IntVarP(&refresh, "refresh", "r", 500, "How often the client should sync to the server in milliseconds")

	return cmd
}

func ParseHost(host string) (netip.AddrPort, error) {
	host, port, err := net.SplitHostPort(host)

	if err != nil {
		return netip.AddrPort{}, err
	}

	if host == "localhost" {
		host = "127.0.0.1"
	}

	return netip.ParseAddrPort(fmt.Sprintf("%s:%s", host, port))
}
