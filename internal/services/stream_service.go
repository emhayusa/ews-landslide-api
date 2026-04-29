package services

import (
	"big-devops-api/internal/config"
	"big-devops-api/internal/models"
	"big-devops-api/internal/repositories"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/websocket/v2"
)

type StreamService struct {
	config         *config.Config
	mqtt           mqtt.Client
	clients        map[*websocket.Conn]bool
	mu             sync.Mutex
	monitoringRepo repositories.MonitoringRepository
	stationRepo    repositories.StationRepository
}

type MQTTMessage struct {
	TS     string `json:"ts"`
	Params struct {
		Bucket    float64 `json:"bucket"`
		Battery   float64 `json:"Baterai"`
		Solar     float64 `json:"Solar"`
		Alarm     int     `json:"alarm"`
		MaxBucket float64 `json:"max_bucket"`
		Deformasi float64 `json:"deformasi"`
	} `json:"params"`
}

func NewStreamService(cfg *config.Config, mRepo repositories.MonitoringRepository, sRepo repositories.StationRepository) *StreamService {
	s := &StreamService{
		config:         cfg,
		clients:        make(map[*websocket.Conn]bool),
		monitoringRepo: mRepo,
		stationRepo:    sRepo,
	}
	s.initMQTT()
	return s
}

func (s *StreamService) initMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", s.config.MQTTHost, s.config.MQTTPort))
	opts.SetClientID(s.config.MQTTClientID)
	opts.SetUsername(s.config.MQTTUser)
	opts.SetPassword(s.config.MQTTPass)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("[MQTT] Connected to broker")
		if token := c.Subscribe(s.config.MQTTTopic, 0, s.onMessageReceived); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] Subscribe error: %v\n", token.Error())
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v\n", err)
	}

	s.mqtt = mqtt.NewClient(opts)
	if token := s.mqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[MQTT] Connect error: %v\n", token.Error())
	}
}

func (s *StreamService) onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	var mqttData MQTTMessage
	if err := json.Unmarshal(msg.Payload(), &mqttData); err != nil {
		log.Printf("[MQTT] JSON parse error: %v\n", err)
		return
	}

	// 1. Find the target station
	// For this specific case, we'll try to find station by name 'Rover Ungaran' 
	// or fallback to the first station available for testing.
	// In production, you'd match by a field in the MQTT payload or topic.
	var targetStationID uint = 0
	stations, _ := s.stationRepo.FindAll()
	for _, st := range stations {
		if st.Name == "Rover Ungaran" || st.StationID == "UG-001" {
			targetStationID = st.ID
			break
		}
	}
	if targetStationID == 0 && len(stations) > 0 {
		targetStationID = stations[0].ID
	}

	if targetStationID != 0 {
		// 2. Save to database
		monitoring := &models.Monitoring{
			StationReferID: targetStationID,
			Timestamp:      time.Now(), // Or parse mqttData.TS if needed
			Bucket:     mqttData.Params.Bucket,
			Battery:    mqttData.Params.Battery,
			Solar:      mqttData.Params.Solar,
			Alarm:      mqttData.Params.Alarm,
			MaxBucket:  mqttData.Params.MaxBucket,
			Deformasi:  mqttData.Params.Deformasi,
			RawPayload: string(msg.Payload()),
		}

		if err := s.monitoringRepo.Create(monitoring); err != nil {
			log.Printf("[DB] Failed to save monitoring data: %v\n", err)
		} else {
			// log.Printf("[DB] Saved monitoring data for station %d\n", targetStationID)
		}
	}

	// 3. Broadcast to WS with station info
	broadcastData := map[string]interface{}{
		"id":           targetStationID,
		"station_id":   "", // Code string
		"station_name": "Rover Ungaran", // Default or found name
		"payload":      mqttData,
	}

	// Update station info if found
	for _, st := range stations {
		if st.ID == targetStationID {
			broadcastData["station_id"] = st.StationID
			broadcastData["station_name"] = st.Name
			break
		}
	}

	s.Broadcast(broadcastData)
}

func (s *StreamService) RegisterClient(c *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = true
	log.Printf("[WS] Client connected. Total clients: %d\n", len(s.clients))
}

func (s *StreamService) UnregisterClient(c *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c)
	log.Printf("[WS] Client disconnected. Total clients: %d\n", len(s.clients))
}

func (s *StreamService) Broadcast(data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("[WS] Marshal error: %v\n", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("[WS] Write error: %v\n", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}
