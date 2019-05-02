package openshiftmachineapi

type scalableResource interface {
 ID() string
 MaxSize() int
 MinSize() int
 Name() string
 Namespace() string
 Nodes() ([]string, error)
 SetSize(nreplicas int32) error
 Replicas() int32
}
