module openebs.io/storage-provisioner

go 1.12

require (
	github.com/container-storage-interface/spec v1.1.0
	github.com/golang/mock v1.3.1
	github.com/kubernetes-csi/csi-lib-utils v0.6.1
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/kubernetes-csi/external-provisioner v1.3.0
	github.com/kubernetes-csi/external-snapshotter v1.2.1
	github.com/miekg/dns v1.1.16 // indirect
	github.com/pkg/errors v0.8.0
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/grpc v1.23.0
	k8s.io/api v0.0.0-20190826194732-9f642ccb7a30
	k8s.io/apimachinery v0.0.0-20190826234335-9a93b3ad8769
	k8s.io/apiserver v0.0.0-20190826120651-3cf73271e6df
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190823013255-e3d4ac5c99fb
	k8s.io/csi-translation-lib v0.0.0-20190823054420-59bd3cbb3c27
	k8s.io/klog v0.4.0
	sigs.k8s.io/sig-storage-lib-external-provisioner v4.0.0+incompatible
)
