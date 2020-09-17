package inmemory

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/TopicSubscriberTypes"
)

type inMemoryPubSub struct{

}

func NewPubSubClient() *inMemoryPubSub {
	return &inMemoryPubSub{}
}

func (inMemoryPubSub) CreateTopicOnPubSub(projectID, topicID string) (*Entities.Topic, error) {
	return &Entities.Topic{
		TopicName:           topicID,
		ChunkSequenceNumber: 0,
		TopicType:           "",
	},nil
}

func (inMemoryPubSub) CreatePullSubscriberForTopicOnPubSub(projectID, subID string, topicID string) (*Entities.TopicSubscriber, error) {
	return &Entities.TopicSubscriber{
		Name:            subID,
		SubscriberType:  TopicSubscriberTypes.PullSubscriber,
		NotificationUrl: "",
		TopicName:       topicID,
	},nil
}

func (inMemoryPubSub) CreatePushSubscriberForTopicOnPubSub(projectID, subID string, topicID string, endpoint string) (*Entities.TopicSubscriber, error) {
	return &Entities.TopicSubscriber{
		Name:            subID,
		SubscriberType:  TopicSubscriberTypes.PushSubscriber,
		NotificationUrl: "",
		TopicName:       topicID,
	},nil
}

func (inMemoryPubSub) PublishMessageToPubSubTopic(data []byte, projectID, topicID string) (string, error) {
	return shared.GenerateUniqueId("IN-MEM-MSG-"),nil
}
