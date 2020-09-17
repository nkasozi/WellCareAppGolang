package Entities

import "gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Enums/TopicTypes"

type Topic struct {
	TopicName           string
	ChunkSequenceNumber int
	TopicType           TopicTypes.TopicType
}
