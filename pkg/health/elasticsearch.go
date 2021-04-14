package health

import (
	"bytes"
	"encoding/base64"

	"github.com/flanksource/commons/console"
	"github.com/flanksource/kommons"
)

func ElasticsearchCheck(c *kommons.Client, namespace string) (bytes.Buffer, error) {
	var testlog bytes.Buffer
	test := console.NewTestResults("elasticsearch", &testlog)
	client, err := c.GetClientset()
	if err != nil {
		test.Failf("elasticsearch", "Failed to get k8s client: %v", err)
		return testlog, err
	}
	kommons.TestNamespace(client, namespace, &test)

	// clusterName := "estest"
	// userName := "elastic"

	// pod, err := c.GetFirstPodByLabelSelector(namespace, fmt.Sprintf("common.k8s.elastic.co/type=elasticsearch,elasticsearch.k8s.elastic.co/cluster-name=%s", clusterName))
	// if err != nil {
	// 	test.Failf("Elasticsearch", "Unable to find elastic pod")
	// 	return testlog, err
	// }

	// dialer, _ := c.GetProxyDialer(proxy.Proxy{
	// 	Namespace:    namespace,
	// 	Kind:         "pods",
	// 	ResourceName: pod.Name,
	// 	Port:         9200,
	// })
	// tr := &http.Transport{
	// 	DialContext:     dialer.DialContext,
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// httpClient := &http.Client{Transport: tr}

	// secret := c.GetSecret(namespace, fmt.Sprintf("%s-es-%s-user", clusterName, userName))
	// if secret == nil {
	// 	test.Failf("Elasticsearch", "Unable to get password for %s user %v", userName, err)
	// 	return testlog, errors.New("could not read elasticsearch secret")
	// }

	// req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s-es-http/_cluster/health", clusterName), nil)
	// req.Header.Add("Authorization", "Basic "+basicAuth(userName, string((*secret)[userName])))

	// resp, err := httpClient.Do(req)
	// if err != nil {
	// 	test.Failf("Elasticsearch", "Failed to get cluster health: %v", err)
	// 	return testlog, err
	// }
	// health := elasticsearch.Health{}
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// if err := json.Unmarshal(body, &health); err != nil {
	// 	test.Failf("Elasticsearch", "Failed to unmarshall :%v", err)
	// 	return testlog, err
	// } else if health.Status == elasticsearch.GreenHealth {
	// 	test.Passf("Elasticsearch", "elasticsearch cluster is: %s", health)
	// 	return testlog, nil
	// } else {
	// 	test.Failf("Elasticsearch", "elasticsearch cluster is: %s", health)
	// 	return testlog, errors.New("elasticsearch cluster not healthy")
	// }

	return testlog, nil
}

// nolint: deadcode
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
