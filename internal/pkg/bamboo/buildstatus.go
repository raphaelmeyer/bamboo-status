package bamboo

import (
	"bamboo/internal/pkg/blinkt"
	"log"
	"time"
)

func Run(done <-chan bool) error {
	config, err := readConfig("config.json")
	if err != nil {
		return err
	}

	password, err := getPassword(config.Username)
	if err != nil {
		return err
	}

	bamboo := NewBamboo(config.Server, config.Username, password)
	bt, close, err := blinkt.NewBlinkt()
	if err != nil {
		return err
	}

	defer func() {
		bt.Clear()
		close()
	}()

	err = bt.Clear()
	if err != nil {
		return err
	}

	for {
		for i, plan := range config.BuildPlans {
			result, err := bamboo.Status(plan)

			if err != nil {
				log.Println(err)
				bt.SetLed(uint32(i), blinkt.Yellow)
			} else if result.State == Success {
				bt.SetLed(uint32(i), blinkt.Green)
			} else if result.State == Building {
				bt.SetLed(uint32(i), blinkt.Cyan)
			} else {
				bt.SetLed(uint32(i), blinkt.Red)
			}
		}

		select {
		case <-done:
			return nil
		case <-time.After(20 * time.Second):
			continue
		}
	}

	return nil
}
