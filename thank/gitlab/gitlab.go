package gitlab

import (
	"github.com/adamliesko/go-thanks/discover"
	"github.com/pkg/errors"
)

type Thanker struct {
	apiToken string
}

func New(token string) Thanker {
	return Thanker{apiToken: token}
}

func (t Thanker) Auth() error {
	return errors.New("not implemented")
}

func (t Thanker) CanThank(discover.Repository) (bool, error) {
	return false, nil
}

func (t Thanker) Thank(discover.Repository) error {
	return nil
}
