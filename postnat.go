package postnat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type App struct {
	listening chan bool

	config   *Config
	postgres *pgx.Conn
	nats     *nats.Conn
}

func New(config *Config) (*App, error) {
	postgres, err := pgx.Connect(context.Background(), config.dbConnStr())
	if err != nil {
		return nil, err
	}
	nats, err := nats.Connect(config.natsConnStr(),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(config.Nats.MaxReconnects),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return nil, err
	}
	app := &App{
		listening: make(chan bool),
		config:    config,
		postgres:  postgres,
		nats:      nats,
	}

	return app, nil
}

func (a *App) Run() chan bool {
	go func() {
		configuredPrefix := a.config.Topics.Prefix

		log.Info().Msg("Setting up listeners")

		if len(a.config.Topics.ListenFor) < 1 {
			log.Error().Msg("Topics to listen_for cannot be empty")
			a.listening <- false
			return
		}

		for _, topic := range a.config.Topics.ListenFor {
			sql := fmt.Sprintf("LISTEN %s;", topic)

			loggedTopic := topic
			if a.config.Topics.ReplaceUnderscore {
				loggedTopic = strings.ReplaceAll(loggedTopic, "_", ".")
			}

			if configuredPrefix != "" {
				if !strings.HasSuffix(configuredPrefix, ".") {
					loggedTopic = configuredPrefix + "." + loggedTopic
				} else {
					loggedTopic = configuredPrefix + loggedTopic
				}
			}

			log.Info().Msgf("LISTEN %s; will be published to NATS as %s", topic, loggedTopic)

			if _, err := a.postgres.Exec(context.Background(), sql); err != nil {
				log.Error().Msgf("Exec unexpectedly failed with %v: %v", sql, err)
				return
			}
		}

		for {
			notification, err := a.postgres.WaitForNotification(context.Background())
			if err != nil {
				log.Error().Err(err).Msg("Failed to handle notification")
				a.listening <- false
				break
			}

			log.Debug().
				Str("channel", notification.Channel).
				Str("topic", notification.Channel).
				Int("content-length", len(notification.Payload)).
				Msg("Publishing message to nats topic")

			topic := notification.Channel
			if a.config.Topics.ReplaceUnderscore {
				topic = strings.ReplaceAll(topic, "_", ".")
			}

			if configuredPrefix != "" {
				if !strings.HasSuffix(configuredPrefix, ".") {
					topic = configuredPrefix + "." + topic
				} else {
					topic = configuredPrefix + topic
				}
			}

			err = a.nats.Publish(topic, []byte(notification.Payload))
			if err != nil {
				log.Error().Err(err).Msg("Failed to notify NATS")
			}
		}
	}()

	return a.listening
}

func (a *App) Shutdown() {
	defer a.postgres.Close(context.Background())
	defer a.nats.Close()
}
