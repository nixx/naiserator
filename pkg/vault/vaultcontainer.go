package vault

import (
	"fmt"
	nais "github.com/nais/naiserator/pkg/apis/nais.io/v1alpha1"
	config2 "github.com/nais/naiserator/pkg/naiserator/config"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	k8score "k8s.io/api/core/v1"
)

type config struct {
	vaultAddr          string
	initContainerImage string
	authPath           string
	app                nais.Application
}

//Creates vault init/sidecar containers
type Creator interface {
	AddVaultContainer(podSpec *k8score.PodSpec) (*k8score.PodSpec, error)
}

func defaultVaultTokenFileName() string {
	return filepath.Join(nais.DefaultVaultMountPath, "vault_token")
}

func (c config) validate() (bool, error) {

	var result = &multierror.Error{}

	if len(c.vaultAddr) == 0 {
		multierror.Append(result, fmt.Errorf("vault address not found in environment"))
	}

	if len(c.initContainerImage) == 0 {
		multierror.Append(result, fmt.Errorf("vault init container image not found in environment"))
	}

	if len(c.authPath) == 0 {
		multierror.Append(result, fmt.Errorf("vault auth path not found in environment"))
	}

	for _, p := range c.app.Spec.Vault.Mounts {
		if len(p.MountPath) == 0 {
			multierror.Append(result, fmt.Errorf("mount path not specified"))
			break
		}

		if len(p.KvPath) == 0 {
			multierror.Append(result, fmt.Errorf("vault kv path not specified"))
			break
		}
	}

	return result.ErrorOrNil() == nil, result.ErrorOrNil()

}

// Enabled checks if this Initializer is enabled
func Enabled() bool {
	return viper.GetBool(config2.FeaturesVault)
}

func defaultKVPath() string {
	return viper.GetString(config2.VaultKvPath)
}

func NewVaultContainerCreator(app nais.Application) (Creator, error) {
	config := config{
		vaultAddr:          viper.GetString(config2.VaultAddress),
		initContainerImage: viper.GetString(config2.VaultInitContainerImage),
		authPath:           viper.GetString(config2.VaultAuthPath),
		app:                app,
	}

	if ok, err := config.validate(); !ok {
		return nil, err
	}

	return config, nil
}

// Add init/sidecar container to pod spec.
func (c config) AddVaultContainer(podSpec *k8score.PodSpec) (*k8score.PodSpec, error) {
	if len(c.app.Spec.Vault.Mounts) == 0 {
		return c.addVaultContainer(podSpec, []nais.SecretPath{c.defaultSecretPath()})
	} else {
		return c.addVaultContainer(podSpec, c.app.Spec.Vault.Mounts)
	}
}

func (c config) addVaultContainer(spec *k8score.PodSpec, paths []nais.SecretPath) (*k8score.PodSpec, error) {
	m := make(map[string]string, len(paths))
	for _, s := range paths {
		if old, exists := m[s.MountPath]; exists {
			return nil, fmt.Errorf("illegal to mount multiple Vault secrets: %s and %s to the same  path: %s", s.KvPath, old, s.MountPath)
		}
		m[s.MountPath] = s.KvPath
	}

	volumeMounts := createVolumeMounts(m)
	args := c.createVaultContainerArgs(m)

	container := k8score.Container{
		Name:         "vks-0",
		VolumeMounts: volumeMounts,
		Args:         args,
		Image:        c.initContainerImage,
		Env: []k8score.EnvVar{
			{
				Name:  "VAULT_AUTH_METHOD",
				Value: "kubernetes",
			},
			{
				Name:  "VAULT_SIDEKICK_ROLE",
				Value: c.app.Name,
			},
			{
				Name:  "VAULT_K8S_LOGIN_PATH",
				Value: c.authPath,
			},
		},
	}

	if c.app.Spec.Vault.Sidecar {
		spec.Containers = append(spec.Containers, container)
	} else {
		spec.InitContainers = append(spec.InitContainers, container)
	}

	spec.Volumes = append(spec.Volumes, k8score.Volume{
		Name: "vault-volume",
		VolumeSource: k8score.VolumeSource{
			EmptyDir: &k8score.EmptyDirVolumeSource{
				Medium: k8score.StorageMediumMemory,
			},
		},
	})

	for i := range spec.Containers {
		if spec.Containers[i].Name == c.app.Name {
			spec.Containers[i].VolumeMounts = append(spec.Containers[i].VolumeMounts, volumeMounts...)
		}
	}
	return spec, nil
}

func (c config) createVaultContainerArgs(paths map[string]string) []string {

	args := []string{
		"-v=10",
		"-logtostderr",
		"-renew-token",
		fmt.Sprintf("-vault=%s", c.vaultAddr),
		fmt.Sprintf("-save-token=%s", defaultVaultTokenFileName()),
	}

	for mountPath, secretPath := range paths {
		args = append(args, fmt.Sprintf("-cn=secret:%s:dir=%s,fmt=flatten", secretPath, mountPath))
	}
	if !c.app.Spec.Vault.Sidecar {
		args = append(args, "-one-shot")
	}

	return args
}

func createVolumeMounts(paths map[string]string) []k8score.VolumeMount {

	volumeMounts := make([]k8score.VolumeMount, 0, len(paths))
	for mountPath := range paths {
		volumeMounts = append(volumeMounts, k8score.VolumeMount{
			Name:      "vault-volume",
			MountPath: mountPath,
			SubPath:   filepath.Join("vault", mountPath), //Just to make sure subpath does not start with "/"
		})
	}

	//Adding default vault mount if it does not exists
	if _, exists := paths[nais.DefaultVaultMountPath]; !exists {
		volumeMounts = append(volumeMounts, k8score.VolumeMount{
			Name:      "vault-volume",
			MountPath: nais.DefaultVaultMountPath,
			SubPath:   filepath.Join("vault", nais.DefaultVaultMountPath), //Just to make sure subpath does not start with "/"
		})
	}
	return volumeMounts
}

func (in config) defaultSecretPath() nais.SecretPath {
	return nais.SecretPath{
		MountPath: nais.DefaultVaultMountPath,
		KvPath:    fmt.Sprintf("%s/%s/%s", defaultKVPath(), in.app.Name, in.app.Namespace),
	}
}
