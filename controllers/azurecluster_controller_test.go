/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-azure/internal/test/mock_log"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ = Describe("AzureClusterReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile an AzureCluster", func() {
		It("should not error and not requeue the request with insufficient set up", func() {

			ctx := context.Background()

			reconciler := &AzureClusterReconciler{
				Client: k8sClient,
				Log:    log.Log,
			}

			instance := &infrav1.AzureCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}

			// Create the AzureCluster object and expect the Reconcile and Deployment to be created
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			result, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				},
			})
			Expect(err).To(BeNil())
			Expect(result.RequeueAfter).To(BeZero())
		})

		It("should fail with context timeout error if context expires", func() {
			log := mock_log.NewMockLogger(gomock.NewController(GinkgoT()))
			log.EXPECT().WithValues(gomock.Any()).DoAndReturn(func(args ...interface{}) logr.Logger {
				time.Sleep(1 * time.Second)
				return log
			})

			reconciler := &AzureClusterReconciler{
				Client:           k8sClient,
				Log:              log,
				ReconcileTimeout: 1 * time.Second,
			}

			instance := &infrav1.AzureCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
			_, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Or(Equal("context deadline exceeded"), Equal("rate: Wait(n=1) would exceed context deadline")))
		})
	})
})
