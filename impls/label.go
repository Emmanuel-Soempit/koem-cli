package impls

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
)

type Label struct {
	Name string
	Min  string
	Max  string
}

func (l *Label) AddLabel(name string, ports []string) error {
	if len(ports) < 2 {
		return fmt.Errorf("port range must have at least 2 ports")
	}
	if len(ports) > 2 {
		return fmt.Errorf("port range must have at most 2 ports - minimum and maximum")
	}

	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return fmt.Errorf("invalid port %q: must be a number", ports[0])
	}
	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return fmt.Errorf("invalid port %q: must be a number", ports[1])
	}
	if min >= max {
		return fmt.Errorf("minimum port (%d) must be less than maximum port (%d)", min, max)
	}

	if err := CheckOverlap(name, min, max); err != nil {
		return err
	}

	l.Name = name
	l.Min = ports[0]
	l.Max = ports[1]
	return nil
}

func (l *Label) Save() error {
	viper.Set(fmt.Sprintf("labels.%s.min", l.Name), l.Min)
	viper.Set(fmt.Sprintf("labels.%s.max", l.Name), l.Max)
	return viper.WriteConfig()
}

func LoadAll() ([]Label, error) {
	labelsMap := viper.GetStringMap("labels")
	labels := make([]Label, 0, len(labelsMap))
	for name := range labelsMap {
		min := viper.GetString(fmt.Sprintf("labels.%s.min", name))
		max := viper.GetString(fmt.Sprintf("labels.%s.max", name))
		labels = append(labels, Label{Name: name, Min: min, Max: max})
	}
	return labels, nil
}

func CheckOverlap(name string, min, max int) error {
	existing, err := LoadAll()
	if err != nil {
		return err
	}
	for _, l := range existing {
		if l.Name == name {
			continue
		}
		eMin, err := strconv.Atoi(l.Min)
		if err != nil {
			continue
		}
		eMax, err := strconv.Atoi(l.Max)
		if err != nil {
			continue
		}
		if min <= eMax && max >= eMin {
			return fmt.Errorf("port range %d-%d overlaps with label %q (%s-%s)", min, max, l.Name, l.Min, l.Max)
		}
	}
	return nil
}

type Reserve struct {
	AppName     string
	Production  string
	Preview     string
	Development string
}

func SaveReserve(labelName, appName string, ports []int) error {
	key := fmt.Sprintf("labels.%s.reserves.%s", labelName, appName)
	if viper.IsSet(key) {
		return fmt.Errorf("app %q is already reserved under label %q", appName, labelName)
	}
	viper.Set(key+".production", fmt.Sprintf("%d", ports[0]))
	viper.Set(key+".preview", fmt.Sprintf("%d", ports[1]))
	viper.Set(key+".development", fmt.Sprintf("%d", ports[2]))
	return viper.WriteConfig()
}

func LoadReserves(labelName string) ([]Reserve, error) {
	path := fmt.Sprintf("labels.%s.reserves", labelName)
	reservesMap := viper.GetStringMap(path)
	reserves := make([]Reserve, 0, len(reservesMap))
	for appName := range reservesMap {
		base := fmt.Sprintf("%s.%s", path, appName)
		reserves = append(reserves, Reserve{
			AppName:     appName,
			Production:  viper.GetString(base + ".production"),
			Preview:     viper.GetString(base + ".preview"),
			Development: viper.GetString(base + ".development"),
		})
	}
	return reserves, nil
}

func ClearReserves(labelName string) error {
	labelsRaw := viper.Get("labels")
	labelsMap, ok := labelsRaw.(map[string]interface{})
	if !ok {
		return nil
	}
	if labelData, ok := labelsMap[labelName]; ok {
		if labelMap, ok := labelData.(map[string]interface{}); ok {
			delete(labelMap, "reserves")
			labelsMap[labelName] = labelMap
		}
	}
	viper.Set("labels", labelsMap)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return viper.ReadInConfig()
}

func FindFreePorts(minStr, maxStr string, count int) ([]int, error) {
	minPort, err := strconv.Atoi(minStr)
	if err != nil {
		return nil, fmt.Errorf("invalid min port %q", minStr)
	}
	maxPort, err := strconv.Atoi(maxStr)
	if err != nil {
		return nil, fmt.Errorf("invalid max port %q", maxStr)
	}

	free := make([]int, 0, count)
	for port := minPort; port <= maxPort && len(free) < count; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		ln.Close()
		free = append(free, port)
	}
	return free, nil
}

func CountPortsInUse(minStr, maxStr string) (int, error) {
	minPort, err := strconv.Atoi(minStr)
	if err != nil {
		return 0, fmt.Errorf("invalid min port %q", minStr)
	}
	maxPort, err := strconv.Atoi(maxStr)
	if err != nil {
		return 0, fmt.Errorf("invalid max port %q", maxStr)
	}

	const maxConcurrency = 100
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	var count atomic.Int32

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			ln, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
			if err != nil {
				count.Add(1)
				return
			}
			ln.Close()
		}(port)
	}

	wg.Wait()
	return int(count.Load()), nil
}
