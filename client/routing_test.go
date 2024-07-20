package client_test

import (
	. "go-p2p/client"
	"reflect"
	"testing"
)

func TestNewRoutingTable(t *testing.T) {
	// func NewRoutingTable(ID string) RoutingTable
	tc := []struct {
		a  string
		ex RoutingTable
		ac RoutingTable
	}{
		{
			a: "1010",
			ex: RoutingTable{
				{
					ID: "0",
				},
				{
					ID: "11",
				},
				{
					ID: "100",
				},
				{
					ID: "1011",
				},
			},
		},
		{
			a: "0011",
			ex: RoutingTable{
				{
					ID: "1",
				},
				{
					ID: "01",
				},
				{
					ID: "000",
				},
				{
					ID: "0010",
				},
			},
		},
	}
	for _, c := range tc {
		ac := NewRoutingTable(c.a)
		if !reflect.DeepEqual(ac, c.ex) {
			t.Errorf("Expected %v, but got %v", c.ex, ac)
		}
	}
}

func TestConstructRoutingTable(t *testing.T) {
	// func ConstructRoutingTable(rt, peerRt RoutingTable) RoutingTable
	tc := []struct {
		rt     RoutingTable
		peerRt RoutingTable
		ex     RoutingTable
	}{
		{
			rt: RoutingTable{
				{
					ID: "0",
				},
				{
					ID: "11",
				},
				{
					ID: "100",
				},
			},
			peerRt: RoutingTable{
				{
					ID: "1",
				},
				{
					ID: "01",
				},
				{
					ID: "000",
				},
			},
			ex: RoutingTable{
				{
					ID: "0",
				},
				{
					ID: "11",
				},
				{
					ID: "100",
				},
			},
		},
	}

	for _, c := range tc {
		ac := ConstructRoutingTable(c.rt, c.peerRt)
		if !reflect.DeepEqual(ac, c.ex) {
			t.Errorf("Expected %v, but got %v", c.ex, ac)
		}
	}
}
