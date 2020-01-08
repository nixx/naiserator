package resourcecreator

import (
	"fmt"

	"github.com/imdario/mergo"
	nais "github.com/nais/naiserator/pkg/apis/nais.io/v1alpha1"
	google_sql_crd "github.com/nais/naiserator/pkg/apis/sql.cnrm.cloud.google.com/v1alpha3"
	k8s_meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AvailabilityTypeRegional     = "REGIONAL"
	AvailabilityTypeZonal        = "ZONAL"
	DefaultSqlInstanceDiskType   = nais.CloudSqlInstanceDiskTypeSSD
	DefaultSqlInstanceAutoBackup = "02:00"
	DefaultSqlInstanceTier       = "db-f1-micro"
	DefaultSqlInstanceDiskSize   = 10
)

func GoogleSqlInstance(app *nais.Application, instance nais.CloudSqlInstance) (*google_sql_crd.SQLInstance, error) {
	objectMeta := app.CreateObjectMeta()
	objectMeta.Name = instance.Name

	if !instance.CascadingDelete {
		ApplyAbandonDeletionPolicy(&objectMeta)
	}

	return &google_sql_crd.SQLInstance{
		TypeMeta: k8s_meta.TypeMeta{
			Kind:       "SQLInstance",
			APIVersion: "sql.cnrm.cloud.google.com/v1alpha3",
		},
		ObjectMeta: objectMeta,
		Spec: google_sql_crd.SQLInstanceSpec{
			DatabaseVersion: string(instance.Type),
			Region:          GoogleRegion,
			Settings: google_sql_crd.SQLInstanceSettings{
				AvailabilityType:    availabilityType(instance.HighAvailability),
				BackupConfiguration: google_sql_crd.SQLInstanceBackupConfiguration{},
				DiskAutoResize:      instance.DiskAutoResize,
				DiskSize:            instance.DiskSize,
				DiskType:            instance.DiskType.GoogleType(),
				Tier:                instance.Tier,
			},
		},
	}, nil
}

func CloudSqlInstanceWithDefaults(instance nais.CloudSqlInstance, appName string) (nais.CloudSqlInstance, error) {
	var err error

	defaultInstance := nais.CloudSqlInstance{
		Name:       appName,
		Tier:       DefaultSqlInstanceTier,
		DiskType:   DefaultSqlInstanceDiskType,
		DiskSize:   DefaultSqlInstanceDiskSize,
		AutoBackup: DefaultSqlInstanceAutoBackup,
		Databases:  []nais.CloudSqlDatabase{{Name: appName}},
	}

	if err = mergo.Merge(&instance, defaultInstance); err != nil {
		return nais.CloudSqlInstance{}, fmt.Errorf("unable to merge default sqlinstance values: %s", err)
	}

	return instance, err
}

func availabilityType(highAvailability bool) string {
	if highAvailability {
		return AvailabilityTypeRegional
	} else {
		return AvailabilityTypeZonal
	}
}
