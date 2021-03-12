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
	CRD         string
	Template    string
	Fixture     string
	Test        func(client *kommons.Client, namespace string) (bytes.Buffer, error)
	ReadyChecks []ReadyCheck
}

type ReadyCheck struct {
	Name                 string
	Resource             string
	ResourceReadyTimeout time.Duration
	WaitFN               kommons.WaitFN
}

var TestSpace []TestInstance

func TestMain(m *testing.M) {
	TestSpace = []TestInstance{
		{
			CRD:      "../config/crd/db/db.flanksource.com_elasticsearchdbs.yaml",
			Template: "../config/templates/elasticsearch-db.yaml",
			Fixture:  "fixtures/elasticsearch-db.yaml",
			Test:     health.ElasticsearchCheck,
			ReadyChecks: []ReadyCheck{
				{
					Resource:             "Elasticsearch",
					ResourceReadyTimeout: 9 * time.Minute,
					WaitFN:               health.IsElasticReady,
				},
			},
		},
		{
			CRD:      "../config/crd/db/db.flanksource.com_redisdbs.yaml",
			Template: "../config/templates/redis-db.yaml",
			Fixture:  "fixtures/redis-db.yaml",
			Test:     health.RedisCheck,
			ReadyChecks: []ReadyCheck{
				{
					Name:                 "rfr-redisdb-e2e",
					Resource:             "StatefulSet",
					ResourceReadyTimeout: 5 * time.Minute,
				},
				{
					Name:                 "rfs-redisdb-e2e",
					Resource:             "Deployment",
					ResourceReadyTimeout: 5 * time.Minute,
				},
			},
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

		for _, check := range fixture.ReadyChecks {
			timeout := check.ResourceReadyTimeout
			if timeout == 0 {
				timeout = 1 * time.Minute
			}
			t.Logf("Waiting for %s %s/%s to be ready", check.Resource, obj.GetNamespace(), obj.GetName())
			waitFN := check.WaitFN
			if waitFN == nil {
				waitFN = client.IsReady
			}
			name := check.Name
			if name == "" {
				name = obj.GetName()
			}
			if _, err := client.WaitForCRD(check.Resource, obj.GetNamespace(), name, timeout, waitFN); err != nil {
				t.Fatalf("Resource %s %s/%s was not ready in %v: %v", check.Resource, obj.GetNamespace(), obj.GetName(), timeout, err)
			}
		}

		t.Log("Checking test instance output")
		testlog, err := fixture.Test(client, obj.GetNamespace())
		t.Log(testlog.String())
		if err != nil {
			t.Fatalf("Elasticsearch test failed: %v", err)
		}
	})
}
