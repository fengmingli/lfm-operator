package status

import (
	v1 "lfm-operator/api/v1"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type errorOption struct{}

func (e errorOption) ApplyOption(_ *v1.VolumeGroup) {}

func (e errorOption) GetResult() (reconcile.Result, error) {
	return reconcile.Result{}, errors.Errorf("error")
}

type successOption struct{}

func (s successOption) ApplyOption(_ *v1.VolumeGroup) {}

func (s successOption) GetResult() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

type retryOption struct{}

func (r retryOption) ApplyOption(_ *v1.VolumeGroup) {}

func (r retryOption) GetResult() (reconcile.Result, error) {
	return reconcile.Result{Requeue: true}, nil
}

func TestDetermineReconciliationResult(t *testing.T) {

	t.Run("A single error option should result in an error return", func(t *testing.T) {
		opts := []Option{
			errorOption{},
			successOption{},
			successOption{},
		}

		res, err := determineReconciliationResult(opts)
		assert.NotNil(t, err)
		assert.Equal(t, false, res.Requeue)
		assert.Equal(t, time.Duration(0), res.RequeueAfter)
	})

	t.Run("An error option takes precedence over a retry", func(t *testing.T) {
		opts := []Option{
			errorOption{},
			retryOption{},
			successOption{},
		}
		res, err := determineReconciliationResult(opts)
		assert.NotNil(t, err)
		assert.Equal(t, false, res.Requeue)
		assert.Equal(t, time.Duration(0), res.RequeueAfter)
	})

	t.Run("No errors will result in a successful reconciliation", func(t *testing.T) {
		opts := []Option{
			successOption{},
			successOption{},
			successOption{},
		}
		res, err := determineReconciliationResult(opts)
		assert.Nil(t, err)
		assert.Equal(t, false, res.Requeue)
		assert.Equal(t, time.Duration(0), res.RequeueAfter)
	})

	t.Run("A retry will take precedence over success", func(t *testing.T) {
		opts := []Option{
			successOption{},
			successOption{},
			retryOption{},
		}
		res, err := determineReconciliationResult(opts)
		assert.Nil(t, err)
		assert.Equal(t, true, res.Requeue)
	})

}
