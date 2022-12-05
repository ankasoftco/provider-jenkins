package clients

import (
	"context"

	"github.com/pkg/errors"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-jenkins/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	jenkins "github.com/bndr/gojenkins"
)

// Config for Jenkins Client authentication struct
type Config struct {
	BaseURL  string
	Username string
	Password string
}

// NewClient creates new Jenkins Client with provided Jenkins Configurations.
func NewClient(c Config) *jenkins.Jenkins {
	jenkins := jenkins.CreateJenkins(nil, c.BaseURL, c.Username, c.Password)

	_, err := jenkins.Init(context.Background())
	if err != nil {
		panic("Something Went Wrong")
		return nil
	}
	return jenkins
}

// GetConfig constructs a Config that can be used to authenticate to Jenkins
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	switch {
	case mg.GetProviderConfigReference() != nil:
		return UseProviderConfig(ctx, c, mg)
	default:
		return nil, errors.New("providerConfigRef is not given")
	}
}

// UseProviderConfig to produce a config that can be used to authenticate to AWS.
func UseProviderConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	pc := &v1alpha1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced Provider")
	}

	t := resource.NewProviderConfigUsageTracker(c, &v1alpha1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	switch s := pc.Spec.Credentials.Source; s { //nolint:exhaustive
	case xpv1.CredentialsSourceSecret:
		csr := pc.Spec.Credentials.SecretRef
		if csr == nil {
			return nil, errors.New("no credentials secret referenced")
		}
		s := &corev1.Secret{}
		if err := c.Get(ctx, types.NamespacedName{Namespace: csr.Namespace, Name: csr.Name}, s); err != nil {
			return nil, errors.Wrap(err, "cannot get credentials secret")
		}
		return &Config{BaseURL: pc.Spec.BaseURL, Username: pc.Spec.Username, Password: string(s.Data[csr.Key])}, nil
	default:
		return nil, errors.Errorf("credentials source %s is not currently supported", s)
	}
}
