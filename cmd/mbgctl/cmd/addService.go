/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.ibm.com/mbg-agent/cmd/mbgctl/state"
	"github.ibm.com/mbg-agent/pkg/protocol"

	httpAux "github.ibm.com/mbg-agent/pkg/protocol/http/aux_func"
)

// updateCmd represents the update command
var addServiceCmd = &cobra.Command{
	Use:   "addService",
	Short: "Add local service to the MBG",
	Long:  `Add local service to the MBG and save it also in the state of the mbgctl`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceId, _ := cmd.Flags().GetString("id")
		serviceIp, _ := cmd.Flags().GetString("ip")
		description, _ := cmd.Flags().GetString("description")

		state.UpdateState()
		state.AddService(serviceId, serviceIp, description)
		addServiceReq(serviceId)

	},
}

func init() {
	rootCmd.AddCommand(addServiceCmd)
	addServiceCmd.Flags().String("id", "", "service id field")
	addServiceCmd.Flags().String("ip", "", "service ip to connect")
	addServiceCmd.Flags().String("description", "", "Service description to connect")
}

func addServiceReq(serviceId string) {
	log.Printf("Start addService %v to ", serviceId)
	s := state.GetService(serviceId)
	mbgIP := state.GetMbgIP()
	svcExp := s.Service
	log.Printf("Service %v", s)

	address := state.GetAddrStart() + mbgIP + "/service"
	j, err := json.Marshal(protocol.ServiceRequest{Id: svcExp.Id, Ip: svcExp.Ip, Description: svcExp.Description})
	if err != nil {
		log.Fatal(err)
	}

	//send
	resp := httpAux.HttpPost(address, j, state.GetHttpClient())
	log.Infof(`Response message for serive %s addservice:  %s`, svcExp.Id, string(resp))
}
