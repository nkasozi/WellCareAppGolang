package pubsubserviceproviders

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
)

type PubSubInterface interface {
	CreateTopicOnPubSub(projectID, topicID string) (*Entities.Topic, error)
	CreatePullSubscriberForTopicOnPubSub(projectID, subID string, topicID string) (*Entities.TopicSubscriber, error)
	CreatePushSubscriberForTopicOnPubSub(projectID, subID string, topicID string, endpoint string) (*Entities.TopicSubscriber, error)
	PublishMessageToPubSubTopic(data []byte, projectID, topicID string) (string, error)
}
