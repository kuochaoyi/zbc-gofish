package zbcgofish_test

import (
	"reflect"
	"testing"

	"github.com/kuochaoyi/zbc-gofish/zbcgofish"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

const (
	BrokerAddr    = "0.0.0.0:26500"
	BPMNProcessId = "order-process-4"
	// JobType	= "payment-service"
)

func TestGetClient(t *testing.T) {
	type args struct {
		brokerAddr string
	}
	tests := []struct {
		name string
		args args
		want zbc.Client
	}{
		// TODO: Add test cases.
		{BrokerAddr, args{BrokerAddr}, zbcgofish.GetClient(BrokerAddr)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zbcgofish.GetClient(tt.args.brokerAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateInstance(t *testing.T) {
	var variables = make(map[string]interface{})
	variables["orderId"] = "31243"

	type args struct {
		BpmnPID  string
		FormVars map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{BPMNProcessId, args{BPMNProcessId, variables}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zbcgofish.CreateInstance(tt.args.BpmnPID, tt.args.FormVars); got != tt.want {
				t.Errorf("CreateInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestGetTopology(t *testing.T) {

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

}*/
