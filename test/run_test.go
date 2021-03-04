package test

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/flanksource/kommons"
)
type TestInstance struct {
	CRD string
	Template string
	Fixture string
}

var TestSpace []TestInstance

func TestMain(m *testing.M) {
	TestSpace = []TestInstance{
		{
			CRD:      "../config/crd/db/db.flanksource.com_elasticsearchdbs.yaml",
			Template: "../config/templates/elasticsearch-db.yaml",
			Fixture:  "fixtures/elasticsearch-db.yaml",
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
		if err != nil { break }
		if obj == nil { continue }
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
		time.Sleep(time.Second*10)
		t.Log("Applying Template")
		if err := ApplyObject(fixture.Template, client); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second*10)
		t.Log("Applying test instance")
		if err := ApplyObject(fixture.Fixture, client); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second*10)

		fixture, err := ioutil.ReadFile(fixture.Fixture)
		if err != nil {
			t.Fatal(err)
		}
		var obj unstructured.Unstructured
		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(fixture)), 1024)
		if err := decoder.Decode(&obj); err != nil {
			t.Fatalf("Could not unmarshal test instance: %v",err)
		}

		client, err := client.GetClientset()
		if err != nil {
			t.Fatal("Could not retrieve client")
		}
		start := time.Now()

		t.Log("Checking test instance output")
		kommons.WaitForNamespace(client, obj.GetNamespace(), time.Second*300)
		if start.Add(time.Second*300).Before(time.Now()) {
			t.Log("Namespace not ready within 5 min")
			t.Fail()
		}
	})
}
