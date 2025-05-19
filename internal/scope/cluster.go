package scope

type Cluster struct {
	Name                     string `json:"name"`
	ControlplaneMachineCount int64  `json:"controlplaneMachineCount"`
	WorkerMachineCount       int64  `json:"workerMachineCount"`
}
