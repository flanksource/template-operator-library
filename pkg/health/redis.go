package health

import (
	"bytes"

	"github.com/flanksource/commons/console"
	"github.com/flanksource/kommons"
)

func RedisCheck(c *kommons.Client, namespace string) (bytes.Buffer, error) {
	var testlog bytes.Buffer
	test := console.NewTestResults("redis", &testlog)
	client, err := c.GetClientset()
	if err != nil {
		test.Failf("redis", "Failed to get k8s client: %v", err)
		return testlog, err
	}
	kommons.TestNamespace(client, namespace, &test)

	return testlog, nil
}
