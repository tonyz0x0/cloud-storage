package config

const (
	// AsyncTransferEnable : File Asynchronous Transfer(Default is Synchronous)
	AsyncTransferEnable = false
	// TransExchangeName : exchange name
	TransExchangeName = "uploadserver.trans"
	// TransOSSQueueName : oss queue name
	TransOSSQueueName = "uploadserver.trans.oss"
	// TransOSSErrQueueName : if oss is failed, write to another queue
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	// TransOSSRoutingKey : routingkey
	TransOSSRoutingKey = "oss"
)

var (
	// RabbitURL : rabbitmq service url
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
)
