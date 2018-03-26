package config


func (t *Config)startMetrics(){
	t.eg.Go(func() error {
		return nil
	})
}