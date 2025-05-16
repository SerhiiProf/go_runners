package services

import (
	"net/http"
	"runners/models"
	"runners/repositories"
	"time"
)

type ResultsService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewResultsService(runnersRepository *repositories.RunnersRepository, resultsRepository *repositories.ResultsRepository) *ResultsService {
	return &ResultsService{
		runnersRepository,
		resultsRepository,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid location",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid position",
			Status:  http.StatusBadRequest,
		}
	}

	currentYear := time.Now().Year()

	if result.Year < 0 || result.Year > currentYear {
		return nil, &models.ResponseError{
			Message: "Invalid year",
			Status:  http.StatusBadRequest,
		}
	}

	raceResult, err := parseRaceResult(result.RaceResult)

	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	// Кажется здесь создание результата и изменение runner'а не в транзакции, сл-но м.б. несогласованность...
	response, responseErr := rs.resultsRepository.CreateResult(result)

	if responseErr != nil {
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}

	// здесь runner не м.б. nil т.к. выше выполняется resultsRepository.CreateResult и она
	// теоретически обязана дать ошибку, если не будет найдет runner при ее выполнении...
	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	// update runner’s personal best
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse personal best",
				Status:  http.StatusInternalServerError,
			}
		}
		// на заметку: храним в виде строки, как она пришла, хотя имхо
		// было бы правильнее унифицировать входящие строки к типовому варианту
		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	// update runner’s season-best
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to parse season best",
					Status:  http.StatusInternalServerError,
				}
			}
			if raceResult < seasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}
	// почему UpdateRunnerResults, а не UpdateRunner? Ведь обновляем вроде runner, в котором хранятся поля лучших результатов.
	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return response, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}

	err := repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
	if err != nil {
		return &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}
	defer func() {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
	}()

	// по-хорошему, должно выполняться в транзакции с UpdateRunnerResults, который вызывается в конце метода... Так и есть, выше добавили
	result, responseErr := rs.resultsRepository.DeleteResult(resultId)
	if responseErr != nil {
		return responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)

	if responseErr != nil {
		return responseErr
	}

	// Checking if the deleted result is personal best for the runner
	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			return responseErr
		}
		runner.PersonalBest = personalBest
	}

	// Checking if the deleted result is season best for the runner
	// м.б. ситуация когда в этом году runner не участвовал в гонках и у него в SeasonBest время из прошлых лет...
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currentYear {
		seasonBest, responseErr := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if responseErr != nil {
			return responseErr
		}

		runner.SeasonBest = seasonBest
	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return responseErr
	}
	repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)
	return nil
}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(timeString[0:2] + "h" + timeString[3:5] + "m" + timeString[6:8] + "s")
}
