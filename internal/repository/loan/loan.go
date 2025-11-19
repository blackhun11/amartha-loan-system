package loan

import (
	"context"
	"errors"
	"fmt"
	"loan_system/internal/model"
	"sync"

	"github.com/bwmarrin/snowflake"
)

//go:generate mockgen -source=loan.go -destination=mock/loan_mock.go -package=mock
type Repository interface {
	FindAll(ctx context.Context) ([]*model.Loan, error)
	Save(ctx context.Context, loan *model.Loan) error
	FindByID(ctx context.Context, id int64) (*model.Loan, error)
	Update(ctx context.Context, loan *model.Loan) error
}

type repository struct {
	mu            sync.RWMutex
	snowflakeNode *snowflake.Node
	loans         map[int64]*model.Loan
}

func NewRepository() Repository {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &repository{
		snowflakeNode: node,
		loans:         make(map[int64]*model.Loan),
	}
}

func (r *repository) FindAll(ctx context.Context) ([]*model.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loans := make([]*model.Loan, 0, len(r.loans))
	for _, loan := range r.loans {
		loans = append(loans, loan)
	}

	return loans, nil
}

func (r *repository) Save(ctx context.Context, loan *model.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if loan.ID == 0 {
		loan.ID = r.snowflakeNode.Generate().Int64()
	}

	if _, exists := r.loans[loan.ID]; exists {
		return errors.New("loan already exists")
	}

	r.loans[loan.ID] = loan
	return nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*model.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loan, exists := r.loans[id]
	if !exists {
		return nil, errors.New("loan not found")
	}

	return loan, nil
}

func (r *repository) Update(ctx context.Context, loan *model.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.loans[loan.ID]; !exists {
		return errors.New("loan not found")
	}

	r.loans[loan.ID] = loan
	return nil
}
