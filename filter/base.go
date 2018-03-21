package filter

type LogFilter interface {
	Filter(event map[string]interface{})(map[string]interface{}, error)
}
