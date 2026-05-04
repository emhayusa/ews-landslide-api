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
	gorilla "github.com/gorilla/websocket"
)

type StreamService struct {
	config         *config.Config
	mqtt           mqtt.Client
	clients        map[*websocket.Conn]bool
	mu             sync.Mutex
	monitoringRepo  repositories.MonitoringRepository
	stationRepo     repositories.StationRepository
	deformationRepo repositories.DeformationRepository
}

type MQTTMessage struct {
	TS     string `json:"ts"`
	Params struct {
		CurahHujanDaily  float64 `json:"curah_hujan_daily"`
		Baterai          float64 `json:"Baterai"`
		Solar            float64 `json:"Solar"`
		Alarm            int     `json:"alarm"`
		CurahHujanHourly float64 `json:"curah_hujan_hourly"`
	} `json:"params"`
}

type DeformationMessage struct {
	TS       string  `json:"ts"`
	Distance float64 `json:"distance"`
	Offset   float64 `json:"offset"`
}

func NewStreamService(cfg *config.Config, mRepo repositories.MonitoringRepository, sRepo repositories.StationRepository, dRepo repositories.DeformationRepository) *StreamService {
	s := &StreamService{
		config:          cfg,
		clients:         make(map[*websocket.Conn]bool),
		monitoringRepo:  mRepo,
		stationRepo:     sRepo,
		deformationRepo: dRepo,
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
		if st.Name == "Rover Ungaran" || st.StationID == "UG-001" || st.StationID == "UNGR" {
			targetStationID = st.ID
			break
		}
	}
	if targetStationID == 0 && len(stations) > 0 {
		targetStationID = stations[0].ID
	}

	if targetStationID != 0 {
		// 2. Parse Timestamp
		parsedTS, err := time.Parse("2006/01/02 15:04:05", mqttData.TS)
		if err != nil {
			parsedTS = time.Now()
		}

		// 3. Save to database
		monitoring := &models.Monitoring{
			StationReferID:   targetStationID,
			Timestamp:        parsedTS,
			CurahHujanDaily:  mqttData.Params.CurahHujanDaily,
			Baterai:          mqttData.Params.Baterai,
			Solar:            mqttData.Params.Solar,
			Alarm:            mqttData.Params.Alarm,
			CurahHujanHourly: mqttData.Params.CurahHujanHourly,
			RawPayload:       string(msg.Payload()),
		}

		if err := s.monitoringRepo.Create(monitoring); err != nil {
			log.Printf("[DB] Failed to save monitoring data: %v\n", err)
		} else {
			// log.Printf("[DB] Saved monitoring data for station %d\n", targetStationID)
		}
	}

	// 3. Broadcast to WS with station info
	broadcastData := map[string]interface{}{
		"type":         "weather",
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

func (s *StreamService) StartDeformationStream(ref, obs string) {
	url := fmt.Sprintf("ws://36.92.41.75:8000/ws/data?ref=%s&obs=%s", ref, obs)
	log.Printf("[WS-Client] Connecting to %s\n", url)

	go func() {
		for {
			c, _, err := gorilla.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Printf("[WS-Client] Dial error: %v, retrying in 5s...\n", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("[WS-Client] Connected to %s\n", url)

			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Printf("[WS-Client] Read error: %v\n", err)
					break
				}

				var data DeformationMessage
				if err := json.Unmarshal(message, &data); err != nil {
					log.Printf("[WS-Client] JSON Unmarshal error: %v\n", err)
					continue
				}

				// Parse TS
				parsedTS, err := time.Parse(time.RFC3339, data.TS)
				if err != nil {
					parsedTS = time.Now()
				}

				// Save to database
				deformation := &models.Deformation{
					TS:       parsedTS,
					Distance: data.Distance,
					Offset:   data.Offset,
					RefCode:  ref,
					ObsCode:  obs,
				}

				if err := s.deformationRepo.Create(deformation); err != nil {
					log.Printf("[DB] Failed to save deformation data: %v\n", err)
				}

				// Broadcast to internal clients
				s.Broadcast(map[string]interface{}{
					"type":       "deformation",
					"station_id": obs,
					"ref_code":   ref,
					"data":       data,
				})
			}

			c.Close()
			log.Println("[WS-Client] Connection closed, retrying in 5s...")
			time.Sleep(5 * time.Second)
		}
	}()
}
