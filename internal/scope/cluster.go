package scope

type Cluster struct {
	Name                     string `json:"name"`
	ControlplaneMachineCount int    `json:"controlplaneMachineCount"`
	WorkerMachineCount       int    `json:"workerMachineCount"`
}
