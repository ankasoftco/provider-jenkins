/*
Copyright 2022 The Crossplane Authors.

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

package jenkinsnode

import (
	"context"
	"fmt"
	jenkins "github.com/bndr/gojenkins"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-jenkins/apis/dashboard/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-jenkins/apis/v1alpha1"
	"github.com/crossplane/provider-jenkins/internal/controller/features"

	clients "github.com/crossplane/provider-jenkins/internal/clients"
)

const (
	errNotJenkinsNode = "managed resource is not a JenkinsNode custom resource"
	errTrackPCUsage   = "cannot track ProviderConfig usage"
	errGetPC          = "cannot get ProviderConfig"
	errGetCreds       = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// Setup adds a controller that reconciles JenkinsNode managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.JenkinsNodeGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.JenkinsNodeGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: clients.NewClient}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.JenkinsNode{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(c clients.Config) *jenkins.Jenkins
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.JenkinsNode)
	if !ok {
		return nil, errors.New(errNotJenkinsNode)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	cfg, err := clients.GetConfig(ctx, c.kube, cr)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\nConnect Completed")
	return &external{kube: c.kube, service: c.newServiceFn(*cfg)}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	kube    client.Client
	service *jenkins.Jenkins
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.JenkinsNode)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotJenkinsNode)
	}

	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", cr)

	forProvider := &cr.Spec.ForProvider
	node, err := c.service.GetNode(ctx, forProvider.Name)
	if err != nil && err.Error() == "No node found" {
		fmt.Println("404 Node Cannot Found: " + forProvider.Name)
		return managed.ExternalObservation{ResourceExists: false}, nil // trigger Create
	} else if err != nil {
		fmt.Println("\nGet Node Error: " + err.Error())
	} else {
		nodeIsOnline, _ := node.IsOnline(ctx)
		if nodeIsOnline {
			fmt.Println("Node is Online")
		} else {
			fmt.Println("Node is Offline")
		}
		fmt.Print("\nNode Exist: " + node.GetName() + " Everything OK\n\n")
	}

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.JenkinsNode)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotJenkinsNode)
	}

	fmt.Printf("Creating: %+v", cr)

	forProvider := &cr.Spec.ForProvider
	node, err := c.service.CreateNode(ctx, forProvider.Name, int(forProvider.NumExecutors), forProvider.Description, forProvider.RemoteFS, forProvider.Label)

	if err != nil {
		fmt.Println("\nCreating Node Error " + err.Error())
	} else {
		fmt.Println("\nNode Successfully Created:  ", node.GetName())
	}

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.JenkinsNode)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotJenkinsNode)
	}

	fmt.Printf("Updating: %+v", cr)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.JenkinsNode)
	if !ok {
		return errors.New(errNotJenkinsNode)
	}

	fmt.Printf("Deleting: %+v", cr)

	forProvider := &cr.Spec.ForProvider
	node, err := c.service.GetNode(ctx, forProvider.Name)
	if err != nil && err.Error() == "No node found" || node == nil {
		fmt.Println("404 Node Cannot Found: " + forProvider.Name)
	} else {
		result, err := c.service.DeleteNode(ctx, forProvider.Name)
		if err != nil {
			fmt.Print("Deleting Node Error: ", err.Error())
		} else {
			if result != true {
				fmt.Print("Node Successfully Deleted")
			} else {
				fmt.Print("Failed to Delete Node")
			}
		}
	}

	return nil
}
