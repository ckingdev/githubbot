package githubbot

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type Session struct {
	RoomName    string
	password    string
	conn        *websocket.Conn
	inbound     chan *PacketEvent
	outbound    chan *PacketEvent
	errChan     chan error
	msgID       int
	port        int
	secret      string
	logger      *logrus.Logger
	uptime      time.Time
	commitMsgID string
}

func (s *Session) connectOnce() error {
	s.logger.Debugln("Connecting to euphoria via TLS...")
	tlsConn, err := tls.Dial("tcp", "euphoria.io:443", &tls.Config{})
	if err != nil {
		s.logger.Warningln("Connection via TLS failed.")
		return err
	}
	roomURL, _ := url.Parse(fmt.Sprintf("wss://euphoria.io/room/%s/ws", s.RoomName))
	wsConn, _, err := websocket.NewClient(tlsConn, roomURL, http.Header{}, 4096, 4096)
	if err != nil {
		s.logger.Warningln("Upgrade of TLS connection to websocket failed.")
		return err
	}
	s.conn = wsConn
	s.logger.Debugln("Connection complete.")
	return nil
}

func (s *Session) connect() error {
	var err error
	for i := 0; i < 5; i++ {
		if err = s.connectOnce(); err == nil {
			go s.sendNick()
			return nil
		} else {
			s.logger.Infof("Error while connecting: %s\n", err)
			time.Sleep(time.Duration(i+1) * time.Second * 5)
		}
	}
	return err
}

func (s *Session) receivePacket() error {
	var packet PacketEvent
	err := s.conn.ReadJSON(&packet)
	if err != nil {
		if err := s.connect(); err != nil {
			return err
		}
		if err := s.conn.ReadJSON(&packet); err != nil {
			return nil
		}
	}
	s.inbound <- &packet
	return nil
}

func (s *Session) receiver() {
	for {
		err := s.receivePacket()
		if err != nil {
			s.logger.Fatalf("Error receiving packet: %s\n", err)
		}
	}
}

func (s *Session) sendPayload(payload interface{}, pType PacketType, packetID string) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		s.logger.Fatalf("Could not marshal payload: %s\n", err)
	}
	packet := &PacketEvent{
		ID:   packetID,
		Type: pType,
	}
	if err := packet.Data.UnmarshalJSON(rawPayload); err != nil {
		s.logger.Fatalf("Could not unmarshal raw message to packet: %s\n", err)
	}
	s.outbound <- packet
}

func (s *Session) sendAuth() {
	s.logger.Infoln("Sending auth.")
	payload := AuthCommand{
		Type:     "passcode",
		Passcode: s.password}
	s.sendPayload(payload, AuthType, strconv.Itoa(s.msgID))
	s.msgID++
}

func (s *Session) sendNick() {
	s.logger.Infoln("Sending nick.")
	payload := NickCommand{Name: "GithubBot"}
	s.sendPayload(payload, NickType, strconv.Itoa(s.msgID))
	s.msgID++
}

func (s *Session) sendMessage(text string, parent string, packetID string) {
	s.logger.Infof("Sending text message: '%s'", text)
	payload := SendCommand{
		Content: text,
		Parent:  parent,
	}
	s.sendPayload(payload, SendType, packetID)
}

func (s *Session) handlePing(p *PacketEvent) {
	s.logger.Debugln("Handling ping.")
	data, err := p.Payload()
	if err != nil {
		panic(err)
	}
	payload, ok := data.(*PingEvent)
	if !ok {
		logrus.Fatalln("Cannot assert *PingEvent as such.")
	}
	out := PingReply{UnixTime: payload.Time}
	s.sendPayload(out, PingReplyType, strconv.Itoa(s.msgID))
	s.msgID++
}

func (s *Session) handleSend(p *PacketEvent) {
	s.logger.Debugln("Handling send-event.")
	data, err := p.Payload()
	if err != nil {
		panic(err)
	}
	payload, ok := data.(*SendEvent)
	if !ok {
		logrus.Fatalln("Cannot assert *SendEvent as such.")
	}
	if payload.Content == "!uptime" {
		since := time.Since(s.uptime)
		s.sendMessage(fmt.Sprintf(
			"This bot has been up for %s.",
			since.String()),
			p.ID, strconv.Itoa(s.msgID))
		s.msgID++
	}
}

func (s *Session) handleSendReply(p *PacketEvent) {
	s.logger.Debugln("Handling send-reply...")
	data, err := p.Payload()
	if err != nil {
		panic(err)
	}
	payload, ok := data.(*SendReply)
	if !ok {
		s.logger.Fatalln("Cannot assert *SendReplyType as such.")
	}
	if strings.HasPrefix(payload.Content, ":repeat:") {
		s.commitMsgID = payload.ID
	}
}

func (s *Session) inboundHandler() {
	for {
		packet := <-s.inbound
		s.logger.Infof("Receiving packet of type '%s'\n", packet.Type)
		switch packet.Type {
		case PingEventType:
			s.handlePing(packet)
		case SendEventType:
			s.handleSend(packet)
		case SendReplyType:
			s.handleSendReply(packet)
		default:
			s.logger.Infof("Unhandled packet type '%s'", packet.Type)
		}
	}
}

func (s *Session) outboundHandler() {
	for {
		packet := <-s.outbound
		s.logger.Infof("Sending packet of type '%s'\n", packet.Type)
		err := s.conn.WriteJSON(packet)
		if err != nil {
			if err := s.connect(); err != nil {
				s.logger.Fatalf("Error sending packet: %s\n", err)
			}
			if err := s.conn.WriteJSON(packet); err != nil {
				s.logger.Fatalf("Error sending packet: %s\n", err)
			}
		}
	}
}

func NewSession(roomName, password string, port int, secret string, logger *logrus.Logger) (*Session, error) {
	inbound := make(chan *PacketEvent)
	outbound := make(chan *PacketEvent)
	errChan := make(chan error)
	s := Session{
		RoomName: roomName,
		password: password,
		inbound:  inbound,
		outbound: outbound,
		errChan:  errChan,
		msgID:    0,
		logger:   logger,
		port:     port,
		secret:   secret,
		uptime:   time.Now(),
	}
	if err := s.connect(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Session) Run() {
	if s.password != "" {
		go s.sendAuth()
	}
	go s.outboundHandler()
	go s.inboundHandler()
	go s.receiver()
	go s.sendNick()
	go s.hookServer(s.port, s.secret)
	go s.droneServer(8082)
	go s.travisServer(8085)
	<-s.errChan
}
