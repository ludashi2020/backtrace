package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"strings"
)

type Result struct {
	i int
	s string
}

var (
	rIp   = []string{"219.141.136.12", "202.106.50.1", "221.179.155.161", "202.96.209.133", "210.22.97.1", "211.136.112.200", "58.60.188.222", "210.21.196.6", "120.196.165.24", "61.139.2.69", "119.6.6.6", "211.137.96.205", "36.111.200.100", "42.48.16.100", "39.134.254.6"}
	rName = []string{"北京电信", "北京联通", "北京移动", "上海电信", "上海联通", "上海移动", "广州电信", "广州联通", "广州移动", "成都电信", "成都联通", "成都移动", "湖南电信", "湖南联通", "湖南移动"}
	ca    = []color.Attribute{color.FgHiYellow, color.FgHiMagenta, color.FgHiBlue, color.FgHiGreen, color.FgHiCyan, color.FgHiRed, color.FgHiWhite}
	m     = map[string]string{"AS4134": "电信163 [普通线路]", "AS4809": "电信CN2 [优质线路]", "AS4837": "联通4837[普通线路]", "AS9929": "联通9929[优质线路]", "AS9808": "移动CMI [普通线路]", "AS58453": "移动CMI [普通线路]", "AS58807": "移动CMIN2[优质线路]"}
)

func trace(ch chan Result, i int) {
	hops, err := Trace(net.ParseIP(rIp[i]))
	if err != nil {
		s := fmt.Sprintf("%v %-15s %v", rName[i], rIp[i], err)
		ch <- Result{i, s}
		return
	}

	for _, h := range hops {
		for _, n := range h.Nodes {
			asn := ipAsn(n.IP.String())
			if asn == "" {
				continue
			} else {
				as := m[asn]
				var c *color.Color
				if strings.Contains(as, "[优质线路]") {
					c = color.New(color.FgHiGreen).Add(color.Bold)
				} else {
					c = color.New(color.FgHiYellow).Add(color.Bold)
				}
				s := fmt.Sprintf("%v %-15s %-23s", rName[i], rIp[i], c.Sprint(as))
				ch <- Result{i, s}
				return
				}
		}
	}

	s := fmt.Sprintf("%v %-15s %v", rName[i], rIp[i], "测试超时")
	ch <- Result{i, s}
}

func ipAsn(ip string) string {
	if isInIPRanges(ip) {
		return "AS58807"
	}

	switch {
	case strings.HasPrefix(ip, "59.43"):
		return "AS4809"
	case strings.HasPrefix(ip, "202.97"):
		return "AS4134"
	case strings.HasPrefix(ip, "218.105"), strings.HasPrefix(ip, "210.51"):
		return "AS9929"
	case strings.HasPrefix(ip, "219.158"):
		return "AS4837"
	case strings.HasPrefix(ip, "223.118"), strings.HasPrefix(ip, "223.119"), strings.HasPrefix(ip, "223.120"), strings.HasPrefix(ip, "223.121"):
		return "AS58453"
	default:
		return ""
	}
}

func isInIPRanges(ip string) bool {
	ipRanges := []string{
		"223.119.8.0/21",
		"223.119.32.0/24",
		"223.119.34.0/24",
		"223.119.35.0/24",
		"223.119.36.0/24",
		"223.119.37.0/24",
		"223.119.100.0/24",
		"223.120.128.0/17",
		"223.120.134.0/23",
		"223.120.138.0/23",
		"223.120.158.0/23",
		"223.120.164.0/22",
		"223.120.168.0/22",
		"223.120.172.0/22",
		"223.120.174.0/23",
		"223.120.184.0/22",
		"223.120.188.0/22",
		"223.120.192.0/23",
		"223.120.200.0/23",
		"223.120.210.0/23",
		"223.120.212.0/23",
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, ipRange := range ipRanges {
		_, subnet, err := net.ParseCIDR(ipRange)
		if err != nil {
			continue
		}

		if subnet.Contains(parsedIP) {
			return true
		}
	}

	return false
}
