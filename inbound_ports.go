package main

import (
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/core"
)

func inboundPortsOf(config *core.Config) map[int]struct{} {
	ports := map[int]struct{}{}
	for _, inbound := range config.Inbound {
		settings, err := inbound.ReceiverSettings.GetInstance()
		if err != nil {
			continue
		}
		receiver, ok := settings.(*proxyman.ReceiverConfig)
		if !ok || receiver.PortList == nil {
			continue
		}
		for _, r := range receiver.PortList.Range {
			for p := r.From; p <= r.To; p++ {
				ports[int(p)] = struct{}{}
			}
		}
	}
	return ports
}
