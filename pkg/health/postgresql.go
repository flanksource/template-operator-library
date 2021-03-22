package health

import (
	"bytes"

	"github.com/flanksource/commons/console"
	"github.com/flanksource/kommons"
)

func PostgresqlCheck(c *kommons.Client, namespace string) (bytes.Buffer, error) {
	var testlog bytes.Buffer
	test := console.NewTestResults("postgresql", &testlog)
	client, err := c.GetClientset()
	if err != nil {
		test.Failf("postgresql", "Failed to get k8s client: %v", err)
		return testlog, err
	}
	kommons.TestNamespace(client, namespace, &test)

	return testlog, nil
}
