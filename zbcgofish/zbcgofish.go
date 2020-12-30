package zbcgofish

/*
	Lib: https://github.com/zeebe-io/zeebe/tree/develop/clients/go
	Example: https://github.com/zeebe-io/zeebe-get-started-go-client
 */

import (
	"context"
	"log"

	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/pb"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

const brokerAddr string = "0.0.0.0:26500"
const jobType string = "payment-service"

var zbclent zbc.Client
var ctx = context.Background()

func GetClient(brokerAddr string) zbc.Client {

	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         brokerAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	return client
}

// Get Roles for Zbc topology.
func GetTopology() {

	zbclient := GetClient(brokerAddr)

	topology, err := zbclient.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}

	for _, broker := range topology.Brokers {
		log.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			log.Println("  Partition", partition.PartitionId,
				":", RoleToString(partition.Role))
		}
	}
}

func RoleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}

func CreateInstance(FormVars map[string]interface{}, BpmnPid string) string {

	zbclient := GetClient(brokerAddr)

	request, err := zbclient.NewCreateInstanceCommand().
		BPMNProcessId(BpmnPid).LatestVersion().
		VariablesFromMap(FormVars)
	if err != nil {
		panic(err)
	}

	// ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	// fmt.Println(msg.String())
	/*
		2020/08/06 14:39:14 workflowKey:2251799813685373 bpmnProcessId:"order-process" version:6 workflowInstanceKey:2251799813685588
	*/
	return msg.String()
}

// Tasks

func HandleTask(jobType string) {

	zbclient := GetClient(brokerAddr)

	jobWorker := zbclient.NewJobWorker().JobType(jobType).Handler(handleJob).Open()
	defer jobWorker.Close()

	jobWorker.AwaitClose()

}

func handleJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	headers, err := job.GetCustomHeadersAsMap()
	if err != nil {
		// failed to handle job as we require the custom job headers
		failJob(client, job)
		return
	}

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	// variables["totalPrice"] = 46.50
	variables["remarks"] = "Nothing."
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["orderId"])
	log.Println("Collect money using payment method:", headers["method"])

	// ctx := context.Background()

	request.Send(ctx)
}

func failJob(client worker.JobClient, job entities.Job) {

	log.Println("Failed to complete job", job.GetKey())

	// ctx := context.Background()

	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
}
