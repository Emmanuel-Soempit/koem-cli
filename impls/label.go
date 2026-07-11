package impls

import (
	"fmt"

	"github.com/spf13/viper"
)

type Label struct {
	Name  string
	Ports []string
}

type Labels struct {
	Labels []Label `mapstructure:"lables"`
}

func (l *Label) AddLabel(name string, ports []string) error {

	if len(ports) < 2 {
		return fmt.Errorf("port range must have at least 2 ports")
	}
	if len(ports) > 2 {
		return fmt.Errorf("port range must have at most 2 ports - minimum and maximum")
	}
	l.Name = name
	l.Ports = ports
	return nil
}

func (l *Label) Save() error {
	viper.Set("lables.name", l.Name)
	viper.Set("lables.ports", l.Ports)
	viper.WriteConfig()
	return nil
}

func (l *Labels) Load() error {
	if err := viper.Unmarshal(l); err != nil {
		return err
	}
	return nil
}
