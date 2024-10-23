package middlewaresusecase

import middlewaresrepository "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresRepository"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresrepository.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresrepository.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}
