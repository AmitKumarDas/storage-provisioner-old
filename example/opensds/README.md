# OpensSDS storage provision
This document describes how to install an OpenSDS with Kubernetes CSI local cluster with storage-provisioner

`Hotpot: OpenSDS Controller Project. Latest release: capri`

- Prec-config:
 - Ubuntu 16.04 (All the installation work is tested on Ubuntu 16.04. )
 - Ansible 
 - Go
 - Docker
 - root user is REQUIRED 
 - Kubernetes version > v1.13
 - Machine should have name less than 14 characters to avoid below issue:
     https://github.com/opensds/opensds/issues/1023

- Download opensds-installer code
```
git clone https://github.com/opensds/opensds-installer.git
cd opensds-installer/ansible
git checkout stable/capri
```

- Follow these docs for the OpenSDS with CSI setup:
`https://github.com/opensds/opensds/wiki/OpenSDS-Cluster-Installation-through-Ansible
https://github.com/opensds/opensds/wiki/OpenSDS-Integration-with-Kubernetes-CSI`

## Following steps describe storage provisioner working with OpenSDS lvm volumes:
- Have a Kubernetes setup with minimum version of 1.13.7

```bash
storage-provisioner > kubectl get node
NAME        STATUS   ROLES    AGE   VERSION
127.0.0.1   Ready    <none>   19d   v1.14.0
```
Note: node name is 127.0.0.1 because we are using local kubernetes cluster for the setup

- Replace node name in example/opensds/storage.yaml

- Verify OpenSDS pool
```
storage-provisioner > osdsctl pool list
+--------------------------------------+-----------------+-------------+--------+---------------+--------------+
| Id                                   | Name            | Description | Status | TotalCapacity | FreeCapacity |
+--------------------------------------+-----------------+-------------+--------+---------------+--------------+
| 6db58762-4254-579b-a0d5-e4ca43856a9e | opensds-volumes |             |        | 10            | 10           |
+--------------------------------------+-----------------+-------------+--------+---------------+--------------+
```
A OpenSDS pool with 10 gb size created 

What actually happened is a loop-device of file 10 gb gets created 

```bash
storage-provisioner > sudo losetup -a
/dev/loop0: [2049]:28830 (/opt/opensds-hotpot-linux-amd64/volumegroups/opensds-volumes.img)
```

- Create a OpenSDS profile:
```
storage-provisioner > osdsctl profile create '{"name": "test", "description": "default policy", "storageType": "block"}'
+--------------------------+--------------------------------------+
| Property                 | Value                                |
+--------------------------+--------------------------------------+
| Id                       | 4a1b6140-693f-427e-a96b-e45bd8e1419c |
| CreatedAt                | 2019-10-15T00:17:15                  |
| Name                     | test                                 |
| Description              | default policy                       |
| StorageType              | block                                |
| ProvisioningProperties   | {                                    |
|                          |   "dataStorage": {                   |
|                          |     "isSpaceEfficient": false        |
|                          |   },                                 |
|                          |   "ioConnectivity": {}               |
|                          | }                                    |
|                          |                                      |
| ReplicationProperties    | {                                    |
|                          |   "dataProtection": {                |
|                          |     "isIsolated": false              |
|                          |   },                                 |
|                          |   "replicaInfos": {}                 |
|                          | }                                    |
|                          |                                      |
| SnapshotProperties       | {                                    |
|                          |   "schedule": {},                    |
|                          |   "retention": {},                   |
|                          |   "topology": {}                     |
|                          | }                                    |
|                          |                                      |
| DataProtectionProperties | {                                    |
|                          |   "dataProtection": {                |
|                          |     "isIsolated": false              |
|                          |   }                                  |
|                          | }                                    |
|                          |                                      |
| CustomProperties         | null                                 |
|                          |                                      |
+--------------------------+--------------------------------------+
```
- Note `profile-id` and replace it in example/opensds/opensds-sc.yaml

- Verify running of csi driver controller
```bash
storage-provisioner > kubectl get sts
NAME                            READY   AGE
csi-attacher-opensdsplugin      1/1     19d
csi-provisioner-opensdsplugin   1/1     19d
csi-snapshotter-opensdsplugin   1/1     19d

storage-provisioner > kubectl get sts -owide
NAME                            READY   AGE   CONTAINERS                                      IMAGES
csi-attacher-opensdsplugin      1/1     19d   csi-attacher,cluster-driver-registrar,opensds   quay.io/k8scsi/csi-attacher:v1.1.1,quay.io/k8scsi/csi-cluster-driver-registrar:v1.0.1,opensdsio/csiplugin:latest
csi-provisioner-opensdsplugin   1/1     19d   csi-provisioner,opensds                         quay.io/k8scsi/csi-provisioner:v1.1.0,opensdsio/csiplugin:latest
csi-snapshotter-opensdsplugin   1/1     19d   csi-snapshotter,opensds                         quay.io/k8scsi/csi-snapshotter:v1.1.0,opensdsio/csiplugin:latest
```

- Verify running of csi node daemon
```bash
storage-provisioner > kubectl get daemonset
NAME                           DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
csi-nodeplugin-opensdsplugin   1         1         1       1            1           <none>          19d

storage-provisioner > kubectl get daemonset -owide
NAME                           DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE   CONTAINERS                      IMAGES                                                                       SELECTOR
csi-nodeplugin-opensdsplugin   1         1         1       1            1           <none>          19d   node-driver-registrar,opensds   quay.io/k8scsi/csi-node-driver-registrar:v1.1.0,opensdsio/csiplugin:latest   app=csi-nodeplugin-opensdsplugin
```
- Verify all the pods to be running
```bash
storage-provisioner > kubectl get po
NAME                                 READY   STATUS    RESTARTS   AGE
csi-attacher-opensdsplugin-0         3/3     Running   0          19d
csi-nodeplugin-opensdsplugin-8tzqs   2/2     Running   0          19d
csi-provisioner-opensdsplugin-0      2/2     Running   0          19d
csi-snapshotter-opensdsplugin-0      2/2     Running   0          19d
```
- Verify csi driver specific StorageClass
```bash
storage-provisioner > kubectl get sc
NAME                   PROVISIONER               AGE
csi-sc-opensdsplugin   csi-opensdsplugin         78m
standard (default)     kubernetes.io/host-path   19d
```
```yaml
harshita_sharma_mayadata_io@hs-open:~$ kubectl get sc csi-sc-opensdsplugin -oyaml
allowedTopologies:
- matchLabelExpressions:
  - key: topology.csi-opensdsplugin/zone
    values:
    - default
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"allowedTopologies":[{"matchLabelExpressions":[{"key":"topology.csi-opensdsplugin/zone","values":["default"]}]}],"apiVersion":"storage.k8s.io/v1","kind":"StorageClass","metadata":{"annotations":{},"name":"csi-sc-opensdsplugin"},"parameters":{"attachMode":"rw","profile":"4a1b6140-693f-427e-a96b-e45bd8e1419c"},"provisioner":"csi-opensdsplugin"}
  creationTimestamp: "2019-10-15T00:17:54Z"
  name: csi-sc-opensdsplugin
  resourceVersion: "365186"
  selfLink: /apis/storage.k8s.io/v1/storageclasses/csi-sc-opensdsplugin
  uid: 3bc093e0-eee1-11e9-baa0-42010a80002d
parameters:
  attachMode: rw
  profile: 4a1b6140-693f-427e-a96b-e45bd8e1419c
provisioner: csi-opensdsplugin
reclaimPolicy: Delete
volumeBindingMode: Immediate
```
### CSI pods running:
```bash
storage-provisioner > kubectl get pods
NAME                                 READY   STATUS    RESTARTS   AGE
csi-attacher-opensdsplugin-0         3/3     Running   0          19d
csi-nodeplugin-opensdsplugin-8tzqs   2/2     Running   0          19d
csi-provisioner-opensdsplugin-0      2/2     Running   0          19d
csi-snapshotter-opensdsplugin-0      2/2     Running   0          19d
```

- Apply the yamls present in ./deploy/kubernetes folder
```bash
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/rbac.yaml
kubectl apply -f deploy/kubernetes/storage_crd.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml
```
- Apply opensds storage class
```bash
kubectl apply -f example/opensds/opensds-sc.yaml
```
- Apply a storage

```bash
kubectl apply -f example/opensds/storage.yaml
```
```yaml
#kubectl get stor magic-stor -oyaml
apiVersion: ddp.mayadata.io/v1alpha1
kind: Storage
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ddp.mayadata.io/v1alpha1","kind":"Storage","metadata":{"annotations":{"storageprovisioner.ddp.mayadata.io/csi-attacher-name":"csi-opensdsplugin","storageprovisioner.ddp.mayadata.io/storageclass-name":"csi-sc-opensdsplugin"},"name":"magic-stor","namespace":"default"},"spec":{"capacity":"1Gi","nodeName":"127.0.0.1"}}
    storageprovisioner.ddp.mayadata.io/csi-attacher-name: csi-opensdsplugin
    storageprovisioner.ddp.mayadata.io/storageclass-name: csi-sc-opensdsplugin
  creationTimestamp: "2019-10-15T01:16:22Z"
  generation: 1
  name: magic-stor
  namespace: default
  resourceVersion: "366080"
  selfLink: /apis/ddp.mayadata.io/v1alpha1/namespaces/default/storages/magic-stor
  uid: 664536a5-eee9-11e9-baa0-42010a80002d
spec:
  capacity: 1Gi
  nodeName: 127.0.0.1
```

- Verify if PVC gets created with below checks
  - It has Storage as its owner
  - It gets bound to a PV via CSI dynamic provisioner

```yaml
#kubectl get pvc default-magic-stor-v95b8 -oyaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    pv.kubernetes.io/bind-completed: "yes"
    pv.kubernetes.io/bound-by-controller: "yes"
    storageprovisioner.ddp.mayadata.io/csi-attacher-name: csi-opensdsplugin
    storageprovisioner.ddp.mayadata.io/node-name: 127.0.0.1
    volume.beta.kubernetes.io/storage-provisioner: csi-opensdsplugin
  creationTimestamp: "2019-10-15T01:16:22Z"
  finalizers:
  - kubernetes.io/pvc-protection
  generateName: default-magic-stor-
  name: default-magic-stor-v95b8
  namespace: default
  ownerReferences:
  - apiVersion: ddp.mayadata.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: Storage
    name: magic-stor
    uid: 664536a5-eee9-11e9-baa0-42010a80002d
  resourceVersion: "366094"
  selfLink: /api/v1/namespaces/default/persistentvolumeclaims/default-magic-stor-v95b8
  uid: 664888de-eee9-11e9-baa0-42010a80002d
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: csi-sc-opensdsplugin
  volumeMode: Filesystem
  volumeName: pvc-664888de-eee9-11e9-baa0-42010a80002d
status:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 1Gi
  phase: Bound
```

- Verify the pv

```yaml
kubectl get pv -oyaml
apiVersion: v1
items:
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    annotations:
      pv.kubernetes.io/provisioned-by: csi-opensdsplugin
    creationTimestamp: "2019-10-15T01:16:24Z"
    finalizers:
    - kubernetes.io/pv-protection
    - external-attacher/csi-opensdsplugin
    name: pvc-664888de-eee9-11e9-baa0-42010a80002d
    resourceVersion: "366095"
    selfLink: /api/v1/persistentvolumes/pvc-664888de-eee9-11e9-baa0-42010a80002d
    uid: 67daaaec-eee9-11e9-baa0-42010a80002d
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 1Gi
    claimRef:
      apiVersion: v1
      kind: PersistentVolumeClaim
      name: default-magic-stor-v95b8
      namespace: default
      resourceVersion: "366082"
      uid: 664888de-eee9-11e9-baa0-42010a80002d
    csi:
      driver: csi-opensdsplugin
      fsType: ext4
      volumeAttributes:
        attachMode: rw
        availabilityZone: default
        lvPath: /dev/opensds-volumes/volume-36bde598-da1e-4754-af07-dc7cb819a688
        name: pvc-664888de-eee9-11e9-baa0-42010a80002d
        poolId: 6db58762-4254-579b-a0d5-e4ca43856a9e
        profileId: 4a1b6140-693f-427e-a96b-e45bd8e1419c
        status: available
        storage.kubernetes.io/csiProvisionerIdentity: 1569445590237-8081-csi-opensdsplugin
      volumeHandle: 36bde598-da1e-4754-af07-dc7cb819a688
    nodeAffinity:
      required:
        nodeSelectorTerms:
        - matchExpressions:
          - key: topology.csi-opensdsplugin/zone
            operator: In
            values:
            - default
    persistentVolumeReclaimPolicy: Delete
    storageClassName: csi-sc-opensdsplugin
    volumeMode: Filesystem
  status:
    phase: Bound
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
```
- Verify VolumeAttachment
  - Check PVC as the owner reference

```yaml
apiVersion: storage.k8s.io/v1
kind: VolumeAttachment
metadata:
  annotations:
    csi.alpha.kubernetes.io/node-id: hs-open,iqn.1993-08.org.debian:01:1681b06d613a,nqn.ini.F17B25E2-DA94-65B1-C403-0F7F223728D2.hs-open,10.128.0.45
  creationTimestamp: "2019-10-15T01:16:24Z"
  finalizers:
  - external-attacher/csi-opensdsplugin
  name: default-magic-stor-v95b8
  ownerReferences:
  - apiVersion: v1
    blockOwnerDeletion: true
    controller: true
    kind: PersistentVolumeClaim
    name: default-magic-stor-v95b8
    uid: 664888de-eee9-11e9-baa0-42010a80002d
  resourceVersion: "366103"
  selfLink: /apis/storage.k8s.io/v1/volumeattachments/default-magic-stor-v95b8
  uid: 67e888b7-eee9-11e9-baa0-42010a80002d
spec:
  attacher: csi-opensdsplugin
  nodeName: 127.0.0.1
  source:
    persistentVolumeName: pvc-664888de-eee9-11e9-baa0-42010a80002d
status:
  attached: true
  attachmentMetadata:
    attachMode: rw
    attachmentId: d583f65f-3081-4cb9-9f9d-4d52abf04d8b
    attachmentStatus: creating
    hostIp: 10.128.0.45
    hostName: hs-open
```

```bash
storage-provisioner > lsblk
NAME                                                                MAJ:MIN RM SIZE RO TYPE MOUNTPOINT
sdb                                                                   8:16   0  35G  0 disk 
loop0                                                                 7:0    0  10G  0 loop 
└─opensds--volumes-volume--36bde598--da1e--4754--af07--dc7cb819a688 253:0    0   1G  0 lvm  
sda                                                                   8:0    0  30G  0 disk 
└─sda1                                                                8:1    0  30G  0 part /
```

```bash
storage-provisioner > ls /dev/disk/by-id/
dm-name-opensds--volumes-volume--36bde598--da1e--4754--af07--dc7cb819a688     lvm-pv-uuid-COxrZM-JNYf-IlEV-M77z-cIZk-ryec-K2DJ7l
dm-uuid-LVM-38YnkV4ZNgOM3U4rSJft4SpmpIVfQGzLEeDe3ot06z9xB3CYTmFht6P56EKB6VFs  scsi-0Google_PersistentDisk_hsopen
google-hsopen                                                                 scsi-0Google_PersistentDisk_hs-open
google-hs-open                                                                scsi-0Google_PersistentDisk_hs-open-part1
google-hs-open-part1
```

