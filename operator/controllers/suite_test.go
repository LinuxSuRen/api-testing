/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package controllers

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	corev1alpha1 "github.com/linuxsuren/api-testing/operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.TODO())
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = corev1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&ATestReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Deploy", func() {
	It("Normal", func() {
		ctx := context.Background()
		err := k8sClient.Create(ctx, &corev1alpha1.ATest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: "default",
			},
			Spec: corev1alpha1.ATestSpec{
				Image:   "busybox",
				Version: "master",
				Persistent: &corev1alpha1.Persistent{
					Enabled: true,
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() bool {
			deploy := &appsv1.Deployment{}
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      "sample",
				Namespace: "default",
			}, deploy)
			return err == nil
		}).WithTimeout(time.Second * 3).WithPolling(time.Second).Should(BeTrue())

		atest := &corev1alpha1.ATest{}
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      "sample",
			Namespace: "default",
		}, atest)
		Expect(err).NotTo(HaveOccurred())

		atest.Spec.Version = "v1.0.0"
		err = k8sClient.Update(ctx, atest)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() bool {
			deploy := &appsv1.Deployment{}
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      "sample",
				Namespace: "default",
			}, deploy)

			return err == nil && deploy.Spec.Template.Spec.Containers[0].Image == "busybox:v1.0.0"
		}).WithTimeout(time.Second * 3).WithPolling(time.Second).Should(BeTrue())
	})
})

var _ = Describe("CombineImageTag", func() {
	It("normal", func() {
		result := CombineImageTag("busybox", "v1.0.0")
		Expect(result).To(Equal("busybox:v1.0.0"))
	})

	It("no tag", func() {
		result := CombineImageTag("busybox", "")
		Expect(result).To(Equal("busybox"))
	})

	It("with tag has", func() {
		result := CombineImageTag("busybox", "v1.0.0:xxxxxxx")
		Expect(result).To(Equal("busybox@v1.0.0:xxxxxxx"))
	})
})
