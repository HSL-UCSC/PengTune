package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type App struct {
	ctx       context.Context
	ns        *server.Server
	nc        *nats.Conn
	pos_sub   *nats.Subscription
	att_sub   *nats.Subscription
	pos_gains PIDGains
	att_gains PIDGains
}

type PIDGains struct {
	Kp [3]float32 `json:"kp"`
	Ki [3]float32 `json:"ki"`
	Kd [3]float32 `json:"kd"`
}

type KnobUpdate struct {
	Knob  string  `json:"knob"`
	Value float32 `json:"value"`
}

// NewApp: creates the struct shell only
func NewApp() *App {
	return &App{}
}

// startup: starts embedded NATS and sets up subscriptions
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start embedded NATS server
	opts := &server.Options{
		Host:  "127.0.0.1",
		Port:  4222,
		Debug: true,
		Trace: true,
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Failed to create NATS server: %v", err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		log.Fatal("NATS server failed to start")
	}
	a.ns = ns
	log.Println("Embedded NATS server started")

	// Connect NATS client
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	a.nc = nc
	log.Println("Connected to embedded NATS")

	// Subscribe to pos.pid.gains
	a.pos_sub, err = nc.Subscribe("pid.pos", func(msg *nats.Msg) {
		var gains PIDGains
		if err := json.Unmarshal(msg.Data, &gains); err != nil {
			log.Printf("Failed to decode pos.gains: %v", err)
			return
		}
		a.pos_gains = gains
		log.Printf("Updated pos gains: %+v", gains)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to pid.pos: %v", err)
	}

	// Subscribe to att.pid.gains
	a.att_sub, err = nc.Subscribe("pid.att", func(msg *nats.Msg) {
		var gains PIDGains
		if err := json.Unmarshal(msg.Data, &gains); err != nil {
			log.Printf("Failed to decode att.gains: %v", err)
			return
		}
		a.att_gains = gains
		log.Printf("Updated att gains: %+v", gains)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to pid.att: %v", err)
	}
}

// shutdown: gracefully clean up
func (a *App) shutdown(ctx context.Context) {
	if a.nc != nil {
		a.nc.Close()
	}
	if a.ns != nil {
		a.ns.Shutdown()
		log.Println("NATS server shut down")
	}
}

// PublishKnob: publish gain update
func (a *App) PublishKnob(update KnobUpdate) error {
	log.Printf("Knob %s changed to %.2f\n", update.Knob, update.Value)

	topic, err := knobToTopic(update.Knob)
	if err != nil {
		log.Printf("Invalid knob ID: %s", update.Knob)
		return err
	}

	payload, err := json.Marshal(update.Value)
	if err != nil {
		log.Printf("Failed to encode gain: %v", err)
		return err
	}

	if err := a.nc.Publish(topic, payload); err != nil {
		log.Printf("Failed to publish to %s: %v", topic, err)
		return err
	}

	log.Printf("Published gain %.2f to %s", update.Value, topic)
	return nil
}

// knobToTopic: maps knob ID to NATS topic
func knobToTopic(knobID string) (string, error) {
	id := strings.ToLower(knobID)
	switch id {
	case "posp":
		return "pid.pos.p", nil
	case "posi":
		return "pid.pos.i", nil
	case "posd":
		return "pid.pos.d", nil
	case "attp":
		return "pid.att.p", nil
	case "atti":
		return "pid.att.i", nil
	case "attd":
		return "pid.att.d", nil
	default:
		return "", fmt.Errorf("unknown knob id: %s", knobID)
	}
}
