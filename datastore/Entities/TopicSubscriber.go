package Entities

import (
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/TopicSubscriberTypes"
)

type TopicSubscriber struct {
	Name string
	SubscriberType TopicSubscriberTypes.TopicSubscriberType
	NotificationUrl string
	TopicName string
}

func (ts TopicSubscriber) IsValid() error {
	if ts.SubscriberType == TopicSubscriberTypes.PushSubscriber && len(ts.NotificationUrl)<=0{
		return fmt.Errorf("Push Susbscriber with Name [%v] must have a Notification Url", ts.Name)
	}
	if len(ts.TopicName)<=0{
		return fmt.Errorf("Topic Name for Subscriber [%v] cant be empty",ts.Name)
	}
	return nil
}
