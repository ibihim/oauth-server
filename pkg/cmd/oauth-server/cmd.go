package oauth_server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/klog/v2"

	configv1 "github.com/openshift/api/config/v1"
	osinv1 "github.com/openshift/api/osin/v1"
	"github.com/openshift/library-go/pkg/serviceability"
)

type OsinServerOptions struct {
	ConfigFile string
	Audit      *options.AuditOptions
}

func NewOsinServerCommand(out, errout io.Writer, stopCh <-chan struct{}) *cobra.Command {
	options := &OsinServerOptions{
		Audit: options.NewAuditOptions(),
	}

	cmd := &cobra.Command{
		Use:   "osinserver",
		Short: "Launch OpenShift osin server",
		Run: func(c *cobra.Command, args []string) {
			if err := options.Validate(); err != nil {
				klog.Fatal(err)
			}

			serviceability.StartProfiler()

			if err := options.RunOsinServer(stopCh); err != nil {
				if kerrors.IsInvalid(err) {
					if details := err.(*kerrors.StatusError).ErrStatus.Details; details != nil {
						fmt.Fprintf(errout, "Invalid %s %s\n", details.Kind, details.Name)
						for _, cause := range details.Causes {
							fmt.Fprintf(errout, "  %s: %s\n", cause.Field, cause.Message)
						}
						os.Exit(255)
					}
				}
				klog.Fatal(err)
			}
		},
	}

	// Handle Flags
	flags := cmd.Flags()
	options.Audit.AddFlags(flags)

	flags.StringVar(&options.ConfigFile, "config", "", "Location of the osin configuration file to run from.")
	cmd.MarkFlagFilename("config", "yaml", "yml")
	cmd.MarkFlagRequired("config")

	return cmd
}

func (o *OsinServerOptions) Validate() error {
	var errs []error

	if len(o.ConfigFile) == 0 {
		errs = append(errs, errors.New("--config is required for this command"))
	}

	if err := o.Audit.Validate(); err != nil {
		errs = append(errs, err...)
	}

	return utilerrors.NewAggregate(errs)
}

func (o *OsinServerOptions) RunOsinServer(stopCh <-chan struct{}) error {
	configContent, err := ioutil.ReadFile(o.ConfigFile)
	if err != nil {
		return err
	}

	// TODO this probably needs to be updated to a container inside openshift/api/osin/v1
	scheme := runtime.NewScheme()
	utilruntime.Must(osinv1.Install(scheme))
	codecs := serializer.NewCodecFactory(scheme)
	obj, err := runtime.Decode(codecs.UniversalDecoder(osinv1.GroupVersion, configv1.GroupVersion), configContent)
	if err != nil {
		return err
	}

	config, ok := obj.(*osinv1.OsinServerConfig)
	if !ok {
		return fmt.Errorf("expected OsinServerConfig, got %T", config)
	}

	return RunOsinServer(config, o.Audit, stopCh)
}
