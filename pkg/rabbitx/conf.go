package rabbitx

// RabbitMQConfig 根配置结构，对应整个 YAML 配置
type RabbitMQConfig struct {
	Addr      string           `yaml:"addr" json:"addr"`           // RabbitMQ 连接地址（AMQP 协议格式）
	QoS       QoSConfig        `yaml:"qos" json:"qos"`             // QoS 流量控制配置
	Exchanges []ExchangeConfig `yaml:"exchanges" json:"exchanges"` // 交换机配置列表
	Queues    []QueueConfig    `yaml:"queues" json:"queues"`       // 队列配置列表
	Bindings  []BindingConfig  `yaml:"bindings" json:"bindings"`   // 交换机-队列绑定配置列表
	Consumers []ConsumerConfig `yaml:"consumers" json:"consumers"` // 消费者配置列表
}

// QoSConfig 流量控制配置（对应 qos 节点）
type QoSConfig struct {
	PrefetchCount int  `yaml:"prefetch-count" json:"prefetch-count"` // 预取消息数量（控制单次获取的消息数）
	PrefetchSize  int  `yaml:"prefetch-size" json:"prefetch-size"`   // 预取消息总大小（0 表示无限制）
	Global        bool `yaml:"global" json:"global"`                 // QoS 生效范围（true：全通道；false：当前消费者）
}

// ExchangeConfig 交换机配置（对应 exchanges 数组元素）
type ExchangeConfig struct {
	Name       string                 `yaml:"name" json:"name"`               // 交换机名称
	Kind       string                 `yaml:"kind" json:"kind"`               // 交换机类型（topic/direct/fanout/headers）
	Durable    bool                   `yaml:"durable" json:"durable"`         // 重启后是否持久化
	AutoDelete bool                   `yaml:"auto-delete" json:"auto-delete"` // 无绑定时是否自动删除
	Internal   bool                   `yaml:"internal" json:"internal"`       // 是否为内部交换机（仅接收内部路由消息）
	NoWait     bool                   `yaml:"no-wait" json:"no-wait"`         // 是否不等待服务器确认
	Arguments  map[string]interface{} `yaml:"arguments" json:"arguments"`     // 自定义扩展参数
}

// QueueConfig 队列配置（对应 queues 数组元素）
type QueueConfig struct {
	Name       string                 `yaml:"name" json:"name"`               // 队列名称
	Durable    bool                   `yaml:"durable" json:"durable"`         // 重启后是否持久化
	AutoDelete bool                   `yaml:"auto-delete" json:"auto-delete"` // 无消费者时是否自动删除
	Exclusive  bool                   `yaml:"exclusive" json:"exclusive"`     // 是否为独占队列（仅当前连接可用）
	NoWait     bool                   `yaml:"no-wait" json:"no-wait"`         // 是否不等待服务器确认
	Arguments  map[string]interface{} `yaml:"arguments" json:"arguments"`     // 自定义扩展参数（如 TTL、DLX 等）
}

// BindingConfig 绑定配置（对应 bindings 数组元素，补充 no-wait 和 arguments）
type BindingConfig struct {
	Exchange   string                 `yaml:"exchange" json:"exchange"`       // 绑定的交换机名称
	Queue      string                 `yaml:"queue" json:"queue"`             // 绑定的队列名称
	RoutingKey string                 `yaml:"routing-key" json:"routing-key"` // 路由键（支持通配符）
	NoWait     bool                   `yaml:"no-wait" json:"no-wait"`         // 是否不等待服务器确认绑定
	Arguments  map[string]interface{} `yaml:"arguments" json:"arguments"`     // 绑定的自定义扩展参数
}

// ConsumerConfig 消费者配置（对应 consumers 数组元素）
type ConsumerConfig struct {
	Queue     string                 `yaml:"queue" json:"queue"`         // 监听的队列名称
	Size      int                    `yaml:"size" json:"size"`           // 消费者数量. max is 10
	AutoAck   bool                   `yaml:"auto-ack" json:"auto-ack"`   // 是否自动确认消息
	Exclusive bool                   `yaml:"exclusive" json:"exclusive"` // 是否为独占消费者
	NoWait    bool                   `yaml:"no-wait" json:"no-wait"`     // 是否不等待服务器确认消费
	Arguments map[string]interface{} `yaml:"arguments" json:"arguments"` // 消费者自定义扩展参数
}
