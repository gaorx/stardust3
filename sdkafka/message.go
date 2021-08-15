package sdkafka

type MessagePriority int

const (
	PN          MessagePriority = 0
	P1          MessagePriority = 1
	P2          MessagePriority = 2
	P3          MessagePriority = 3
	P4          MessagePriority = 4
	PriorityNum                 = 4
)

type Message struct {
	Topic    string
	Value    interface{}
	Priority MessagePriority
}

func (msg *Message) KafkaTopic() string {
	return topicOfPriority(msg.Topic, msg.Priority)
}
