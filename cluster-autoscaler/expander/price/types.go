package price

import (
 apiv1 "k8s.io/api/core/v1"
)

type PreferredNodeProvider interface{ Node() (*apiv1.Node, error) }
type NodeUnfitness func(preferredNode, evaluatedNode *apiv1.Node) float64
