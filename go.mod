module flanksource/template-operator-dbs

go 1.16

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	golang.org/x/net v0.0.0-20210226101413-39120d07d75e // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.20.2 // indirect
	k8s.io/apimachinery v0.20.4
	k8s.io/klog/v2 v2.5.0 // indirect
	sigs.k8s.io/controller-runtime v0.8.2
	sigs.k8s.io/structured-merge-diff/v4 v4.0.3 // indirect
)

replace (
	gopkg.in/hairyhenderson/yaml.v2 => github.com/maxaudron/yaml v0.0.0-20190411130442-27c13492fe3c
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
