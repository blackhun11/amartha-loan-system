package loan_test

import (
	"context"
	"loan_system/internal/model"
	"loan_system/internal/repository/loan"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	repo := loan.NewRepository()

	t.Run("Save and FindByID", func(t *testing.T) {
		l := &model.Loan{Principal: 1000}
		err := repo.Save(context.TODO(), l)
		assert.NoError(t, err)
		assert.NotZero(t, l.ID)

		found, err := repo.FindByID(context.TODO(), l.ID)
		assert.NoError(t, err)
		assert.Equal(t, l.Principal, found.Principal)
	})

	t.Run("Update existing loan", func(t *testing.T) {
		l := &model.Loan{Principal: 2000}
		err := repo.Save(context.TODO(), l)
		assert.NoError(t, err)
		assert.NotZero(t, l.ID)

		l.Principal = 3000
		err = repo.Update(context.TODO(), l)
		assert.NoError(t, err)

		updated, err := repo.FindByID(context.TODO(), l.ID)
		assert.NoError(t, err)
		assert.Equal(t, float64(3000), updated.Principal)
	})

	t.Run("Concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				l := &model.Loan{Principal: 500}
				err := repo.Save(context.TODO(), l)
				assert.NoError(t, err)
				assert.NotZero(t, l.ID)

				_, err = repo.FindByID(context.TODO(), l.ID)
				assert.NoError(t, err)
			}()
		}
		wg.Wait()
	})

	t.Run("Find non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(context.TODO(), 999999)
		assert.Error(t, err)
	})

	t.Run("Update non-existent loan", func(t *testing.T) {
		err := repo.Update(context.TODO(), &model.Loan{ID: 999999})
		assert.Error(t, err)
	})
}
