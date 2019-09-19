## Running on command line

For debugging, it's possible to run the storage-provisioner on command line:

```sh
# cd to root of this project
make

d-storprovisioner -kubeconfig ~/.kube/config -v 5
```

## Implementation details

The storage-provisioner follows [controller](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md) pattern and uses informers to watch for `Storage` and `PersistentVolumeClaim` create/update/delete events.

## Troubleshooting

### 1/ No kind is registered

```
Sync failed: Will re-queue storage "default:magic-stor": no kind is registered for the type v1alpha1.Storage in scheme "k8s.io/client-go/kubernetes/scheme/register.go:67"
```

- Try deleting storage provisioner pod
- Reduce the resync interval time

```bash
E0919 11:18:31.755351       1 controller.go:219] ddp-storage-provisioner: Sync failed: Will re-queue storage "default:magic-stor": no kind is registered for the type v1alpha1.Storage in scheme "k8s.io/client-go/kubernetes/scheme/register.go:67"
I0919 11:18:33.755665       1 controller.go:228] ddp-storage-provisioner: Sync started: Storage "default:magic-stor"
E0919 11:18:33.755723       1 controller.go:219] ddp-storage-provisioner: Sync failed: Will re-queue storage "default:magic-stor": no kind is registered for the type v1alpha1.Storage in scheme "k8s.io/client-go/kubernetes/scheme/register.go:67"
  
No resources found.
  storage-provisioner > kubectl get po -n ddp
NAME                                       READY   STATUS    RESTARTS   AGE
ddp-storage-provisioner-5668cdb69f-dfxrb   1/1     Running   0          63s
  storage-provisioner > 
  storage-provisioner > kubectl delete po -n ddp --all
pod "ddp-storage-provisioner-5668cdb69f-dfxrb" deleted
  storage-provisioner > 
  storage-provisioner > 
  storage-provisioner > kubectl get po -n ddp
NAME                                       READY   STATUS    RESTARTS   AGE
ddp-storage-provisioner-5668cdb69f-q4kb2   1/1     Running   0          12s
  storage-provisioner > 
  storage-provisioner > 
  storage-provisioner > kubectl -n ddp logs ddp-storage-provisioner-5668cdb69f-q4kb2
I0919 11:19:06.236421       1 main.go:105] Version: master-unreleased
I0919 11:19:06.240010       1 controller.go:148] Starting ddp-storage-provisioner
I0919 11:19:06.240224       1 reflector.go:120] Starting reflector *v1.PersistentVolumeClaim (10m0s) from k8s.io/client-go/informers/factory.go:134
I0919 11:19:06.240243       1 reflector.go:158] Listing and watching *v1.PersistentVolumeClaim from k8s.io/client-go/informers/factory.go:134
I0919 11:19:06.241812       1 reflector.go:120] Starting reflector *v1beta1.VolumeAttachment (10m0s) from k8s.io/client-go/informers/factory.go:134
I0919 11:19:06.241828       1 reflector.go:158] Listing and watching *v1beta1.VolumeAttachment from k8s.io/client-go/informers/factory.go:134
I0919 11:19:06.242175       1 reflector.go:120] Starting reflector *v1alpha1.Storage (10m0s) from github.com/AmitKumarDas/storage-provisioner/client/generated/informer/externalversions/factory.go:117
I0919 11:19:06.242193       1 reflector.go:158] Listing and watching *v1alpha1.Storage from github.com/AmitKumarDas/storage-provisioner/client/generated/informer/externalversions/factory.go:117
I0919 11:19:06.342781       1 shared_informer.go:227] caches populated
I0919 11:19:06.343228       1 controller.go:228] ddp-storage-provisioner: Sync started: Storage "default:magic-stor"
I0919 11:19:06.350465       1 controller.go:276] ddp-storage-provisioner: Sync started: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:06.350738       1 pvc.go:73] PVCReconciler default/default-magic-stor-7fs9p: Reconcile ignored: Volume not bound
I0919 11:19:06.350827       1 controller.go:292] ddp-storage-provisioner: Sync completed: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:06.354950       1 controller.go:244] ddp-storage-provisioner: Sync completed: Storage "default:magic-stor"
I0919 11:19:06.373925       1 controller.go:276] ddp-storage-provisioner: Sync started: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:06.374153       1 pvc.go:73] PVCReconciler default/default-magic-stor-7fs9p: Reconcile ignored: Volume not bound
I0919 11:19:06.374245       1 controller.go:292] ddp-storage-provisioner: Sync completed: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:10.422247       1 controller.go:276] ddp-storage-provisioner: Sync started: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:10.436168       1 controller.go:292] ddp-storage-provisioner: Sync completed: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:10.439872       1 controller.go:276] ddp-storage-provisioner: Sync started: PVC "default:default-magic-stor-7fs9p"
I0919 11:19:10.440144       1 pvc.go:105] PVCReconciler default/default-magic-stor-7fs9p: No change to desired state
I0919 11:19:10.440265       1 controller.go:292] ddp-storage-provisioner: Sync completed: PVC "default:default-magic-stor-7fs9p"
  storage-provisioner > 
```