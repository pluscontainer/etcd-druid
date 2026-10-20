// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package statefulset

import (
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
)

const (
	// inMemoryEtcdName is the name of the etcd whose data volume is memory-backed.
	inMemoryEtcdName = "etcd-main"
	// inMemoryVolumeName is the name of the memory-backed data volume.
	inMemoryVolumeName = "etcd-in-memory"
	// inMemoryVolumeOverhead is the space on top of twice the quota for WAL and snapshot files.
	inMemoryVolumeOverhead int64 = 1024 * 1024 * 1024
	// backupRestorePreStopSleepSeconds is the duration of the backup-restore preStop sleep.
	backupRestorePreStopSleepSeconds = 15
)

// getInMemoryVolumes returns the memory-backed data volume of an in-memory etcd, sized to twice the quota plus overhead.
func (b *stsBuilder) getInMemoryVolumes() []corev1.Volume {
	if b.etcd.Name != inMemoryEtcdName {
		return nil
	}
	quota := defaultQuota
	if b.etcd.Spec.Etcd.Quota != nil {
		quota = b.etcd.Spec.Etcd.Quota.Value()
	}
	return []corev1.Volume{
		{
			Name: inMemoryVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium:    corev1.StorageMediumMemory,
					SizeLimit: apiresource.NewQuantity(2*quota+inMemoryVolumeOverhead, apiresource.BinarySI),
				},
			},
		},
	}
}

// getBackupRestoreLifecycle returns the lifecycle of the backup-restore container with a preStop sleep.
func getBackupRestoreLifecycle() *corev1.Lifecycle {
	return &corev1.Lifecycle{
		PreStop: &corev1.LifecycleHandler{
			Sleep: &corev1.SleepAction{Seconds: backupRestorePreStopSleepSeconds},
		},
	}
}
