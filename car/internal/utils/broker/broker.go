package broker

type Broker struct {
	Addr   string
	Topics Topics
}

type Topics struct {
	Metrics MetricTopics
	Images  ImageTopics
}

type ImageTopics struct {
	Save   string
	Delete string
}

type MetricTopics struct {
	DecreaseActiveRentsAmount string
	IncreaseActiveRentsAmount string
	IncreaseRentsCancel       string
	NotifyBrandDemand         string
}

type SaveImageMessage struct {
	Value []byte `json:"value"`
	UUID  string `json:"UUID"`
	Idx   int    `json:"idx"`
}
