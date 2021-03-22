package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/flanksource/kommons"
	"github.com/flanksource/template-operator-library/pkg/health"
	"github.com/mitchellh/go-homedir"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

type TestInstance struct {
	CRD                  string
	Template             string
	Fixture              string
	Test                 func(client *kommons.Client, namespace string) (bytes.Buffer, error)
	ResourceReadyTimeout time.Duration
}

var TestSpace []TestInstance

func TestMain(m *testing.M) {
	TestSpace = []TestInstance{
		{
			CRD:                  "../config/crd/db/db.flanksource.com_elasticsearchdbs.yaml",
			Template:             "../config/templates/elasticsearch-db.yaml",
			Fixture:              "fixtures/elasticsearch-db.yaml",
			Test:                 health.ElasticsearchCheck,
			ResourceReadyTimeout: 8 * time.Minute,
		},
		{
			CRD:                  "../config/crd/db/db.flanksource.com_redisdbs.yaml",
			Template:             "../config/templates/redis-db.yaml",
			Fixture:              "fixtures/redis-db.yaml",
			Test:                 health.RedisCheck,
			ResourceReadyTimeout: 5 * time.Minute,
		},
		{
			CRD:                  "../config/crd/db/db.flanksource.com_postgresqldbs.yaml",
			Template:             "../config/templates/postgresql-db.yaml",
			Fixture:              "fixtures/postgresql-db.yaml",
			Test:                 health.PostgresqlCheck,
			ResourceReadyTimeout: 5 * time.Minute,
		},
	}
	code := m.Run()
	os.Exit(code)
}

func TestRunChecks(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatalf("Coulf not determine home path: %v", err)
	}
	kubefile := path.Join(home, ".kube/config")
	config, err := ioutil.ReadFile(kubefile)
	if err != nil {
		t.Fatalf("Could not open kube config: %v", err)
	}
	client, err := kommons.NewClientFromBytes(config)
	if err != nil {
		t.Fatalf("Could not create kubernetes client: %v", err)
	}
	client.GetRESTConfig = client.GetRESTConfigFromKubeconfig
	client.GetKustomizePatches = func() ([]string, error) {
		return []string{}, nil
	}
	wg := sync.WaitGroup{}
	for _, fixture := range TestSpace {
		wg.Add(1)
		_fixture := fixture
		go func() {
			runFixture(t, _fixture, client)
			wg.Done()
		}()
	}
	wg.Wait()
}

func ApplyObject(path string, client *kommons.Client) error {
	var obj *unstructured.Unstructured
	var err error
	if err != nil {
		return err
	}
	crd, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(crd)), 1024)

	for {
		err := decoder.Decode(&obj)
		if err != nil {
			break
		}
		if obj == nil {
			continue
		}
		err = client.ApplyUnstructured(obj.GetNamespace(), obj)
		if err != nil {
			return fmt.Errorf("error decoding %s: %s", crd, err)
		}
	}
	return nil
}

func runFixture(t *testing.T, fixture TestInstance, client *kommons.Client) {
	t.Run(fixture.Fixture, func(t *testing.T) {
		t.Log("Applying CRD")
		if err := ApplyObject(fixture.CRD, client); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 10)
		t.Log("Applying Template")
		if err := ApplyObject(fixture.Template, client); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 10)
		t.Log("Applying test instance")
		if err := ApplyObject(fixture.Fixture, client); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 10)

		testfixture, err := ioutil.ReadFile(fixture.Fixture)
		if err != nil {
			t.Fatal(err)
		}
		var obj unstructured.Unstructured
		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(testfixture)), 1024)
		if err := decoder.Decode(&obj); err != nil {
			t.Fatalf("Could not unmarshal test instance: %v", err)
		}

		timeout := fixture.ResourceReadyTimeout
		if timeout == 0 {
			timeout = 1 * time.Minute
		}
		t.Logf("Waiting for %s %s/%s to be ready", obj.GetKind(), obj.GetNamespace(), obj.GetName())
		_, err = client.WaitForCRD(obj.GetKind(), obj.GetNamespace(), obj.GetName(), timeout, client.IsReadyWithConditions)
		if err != nil {
			t.Fatalf("Resource %s %s/%s was not ready in %v: %v", obj.GetKind(), obj.GetNamespace(), obj.GetName(), timeout, err)
		}

		// js, _ := json.Marshal(item.Object["status"].(map[string]interface{})["conditions"])
		// t.Logf("Item conditions %s\n", js)

		// clientset, _ := client.GetClientset()
		// deployments, _ := clientset.AppsV1().Deployments(obj.GetNamespace()).List(context.Background(), metav1.ListOptions{})
		// for _, depl := range deployments.Items {
		// 	fmt.Printf("Deployment %s %d/%d\n", depl.Name, depl.Status.ReadyReplicas, depl.Status.Replicas)
		// }

		// sts, _ := clientset.AppsV1().StatefulSets(obj.GetNamespace()).List(context.Background(), metav1.ListOptions{})
		// for _, s := range sts.Items {
		// 	fmt.Printf("Deployment %s %d/%d\n", s.Name, s.Status.ReadyReplicas, s.Status.Replicas)
		// }

		t.Log("Checking test instance output")
		testlog, err := fixture.Test(client, obj.GetNamespace())
		t.Log(testlog.String())
		if err != nil {
			t.Fatalf("test failed: %v", err)
		}
	})
}
