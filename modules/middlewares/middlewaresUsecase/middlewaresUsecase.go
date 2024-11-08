package middlewaresusecase

import middlewaresrepository "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresRepository"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresrepository.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresrepository.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}
