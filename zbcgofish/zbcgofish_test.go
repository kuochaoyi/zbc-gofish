package zbcgofish_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

const brokerAddr string = "0.0.0.0:26500"
const BpmnPid string = "order-process"
const jobType string = "payment-service"

var zbclient zbc.Client = GetClient(brokerAddr)
var ctx = context.Background()

func TestGetClient(t *testing.T) {

	zbclient := GetClient(brokerAddr)

	if zbclient == nil {
		t.Error("Fail to get a ZbClient.")
	}

}

func TestGetTopology(t *testing.T) {

	topology, err := zbclient.NewTopologyCommand().Send(ctx)

	if err != nil {
		t.Error(err)
	} else {
		for _, broker := range topology.Brokers {
			fmt.Println("Broker", broker.Host, ":", broker.Port)
			for _, partition := range broker.Partitions {
				fmt.Println("  Partition", partition.PartitionId, ":", RoleToString(partition.Role))
			}
		}

	}

}

func TestCreateInstance(t *testing.T) {

	formVars := make(map[string]interface{})
	formVars["orderId"] = "abc1111"

	msg := CreateInstance(formVars, BpmnPid)

	log.Println(msg)

}

func TestHandleTask(t *testing.T) {

	HandleTask(jobType)

}
