// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package statefulset

import (
	"testing"

	testutils "github.com/gardener/etcd-druid/test/utils"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

func TestGetInMemoryVolumes(t *testing.T) {
	tests := []struct {
		name          string
		etcdName      string
		quota         *apiresource.Quantity
		wantSizeLimit *apiresource.Quantity
	}{
		{"no volume for other etcds", "etcd-events", nil, nil},
		{"default quota", inMemoryEtcdName, nil, ptr.To(apiresource.MustParse("17Gi"))},
		{"custom quota", inMemoryEtcdName, ptr.To(apiresource.MustParse("2Gi")), ptr.To(apiresource.MustParse("5Gi"))},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			etcd := testutils.EtcdBuilderWithDefaults(tc.etcdName, testutils.TestNamespace).Build()
			etcd.Spec.Etcd.Quota = tc.quota

			volumes := (&stsBuilder{etcd: etcd}).getInMemoryVolumes()

			if tc.wantSizeLimit == nil {
				g.Expect(volumes).To(BeEmpty())
				return
			}
			g.Expect(volumes).To(HaveLen(1))
			g.Expect(volumes[0].Name).To(Equal(inMemoryVolumeName))
			g.Expect(volumes[0].EmptyDir.Medium).To(Equal(corev1.StorageMediumMemory))
			g.Expect(volumes[0].EmptyDir.SizeLimit.Cmp(*tc.wantSizeLimit)).To(BeZero())
		})
	}
}

func TestGetBackupRestoreLifecycle(t *testing.T) {
	g := NewWithT(t)

	lifecycle := getBackupRestoreLifecycle()

	g.Expect(lifecycle.PreStop.Exec).To(BeNil())
	g.Expect(lifecycle.PreStop.Sleep).To(Equal(&corev1.SleepAction{Seconds: backupRestorePreStopSleepSeconds}))
}
