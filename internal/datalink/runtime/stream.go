package runtime

import (
	"time"

	"go-gateway/internal/datalink/schema"
)

type pointMeta struct {
	DeviceID string
	Address  string
}

type ValueEvent struct {
	DeviceID         string             `json:"device_id"`
	PointID          string             `json:"point_id"`
	Address          string             `json:"address"`
	RawValue         any                `json:"raw_value"`
	TransformedValue any                `json:"transformed_value"`
	Quality          schema.QualityFlag `json:"quality"`
	Stale            bool               `json:"stale"`
	Timestamp        time.Time          `json:"timestamp"`
}

type valueSubscriber struct {
	deviceID string
	pointIDs map[string]struct{}
	ch       chan ValueEvent
}

func (s *Service) registerPointMeta(pointID string, meta pointMeta) {
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	if s.pointMetaIndex == nil {
		s.pointMetaIndex = make(map[string]pointMeta)
	}
	s.pointMetaIndex[pointID] = meta
}

func (s *Service) lookupPointMeta(pointID string) pointMeta {
	s.pointMetaMu.RLock()
	defer s.pointMetaMu.RUnlock()
	if s.pointMetaIndex == nil {
		return pointMeta{}
	}
	return s.pointMetaIndex[pointID]
}

func (s *Service) SubscribeValueEvents(deviceID string, pointIDs []string) (<-chan ValueEvent, func()) {
	s.subscriberMu.Lock()
	defer s.subscriberMu.Unlock()

	if s.subscribers == nil {
		s.subscribers = make(map[int64]valueSubscriber)
	}

	pointSet := make(map[string]struct{}, len(pointIDs))
	for _, pointID := range pointIDs {
		pointSet[pointID] = struct{}{}
	}

	id := s.nextSubscriberID.Add(1)
	ch := make(chan ValueEvent, 32)
	s.subscribers[id] = valueSubscriber{
		deviceID: deviceID,
		pointIDs: pointSet,
		ch:       ch,
	}

	return ch, func() {
		s.subscriberMu.Lock()
		defer s.subscriberMu.Unlock()

		sub, exists := s.subscribers[id]
		if !exists {
			return
		}

		delete(s.subscribers, id)
		close(sub.ch)
	}
}

func (s *Service) broadcastValueEvent(event ValueEvent) {
	s.subscriberMu.RLock()
	defer s.subscriberMu.RUnlock()

	for _, sub := range s.subscribers {
		if sub.deviceID != "" && sub.deviceID != event.DeviceID {
			continue
		}
		if len(sub.pointIDs) > 0 {
			if _, exists := sub.pointIDs[event.PointID]; !exists {
				continue
			}
		}

		select {
		case sub.ch <- event:
		default:
		}
	}
}
