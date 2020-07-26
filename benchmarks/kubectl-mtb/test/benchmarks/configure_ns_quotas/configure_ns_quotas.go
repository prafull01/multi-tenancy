package configurensquotas

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/bundle/box"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/pkg/benchmark"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/test"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/test/utils"
)

var b = &benchmark.Benchmark{
	// Check if user can list nodes
	PreRun: func(tenantNamespace string, kclient, tclient *kubernetes.Clientset) error {
		resources := []utils.GroupResource{
			{
				APIGroup: "",
				APIResource: metav1.APIResource{
					Name: "resourcequotas",
				},
			},
		}
		verb := "list"
		for _, resource := range resources {

			access, msg, err := utils.RunAccessCheck(tclient, tenantNamespace, resource, verb)
			if err != nil {
				fmt.Println(err.Error())
			}
			if !access {
				return fmt.Errorf(msg)
			}
		}
		return nil
	},
	Run: func(tenantNamespace string, kclient, tclient *kubernetes.Clientset) error {

		resourceNameList := [3]string{"cpu", "ephemeral-storage", "memory"}
		tenantResourcequotas := utils.GetTenantResoureQuotas(tenantNamespace, tclient)
		expectedVal := strings.Join(tenantResourcequotas, " ")
		for _, r := range resourceNameList {
			if !strings.Contains(expectedVal, r) {
				return fmt.Errorf("%s must be configured in tenant %s namespace resource quotas", r, tenantNamespace)
			}
		}

		return nil
	},
}

func init() {
	// Get the []byte representation of a file, or an error if it doesn't exist:
	err := b.ReadConfig(box.Get("configure_ns_quotas/config.yaml"))
	if err != nil {
		fmt.Println(err)
	}

	test.BenchmarkSuite.Add(b)
}
