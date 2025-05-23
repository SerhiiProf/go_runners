package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	ctx := context.Background()
	transaction, err := resultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	runnersRepository.transaction = transaction
	resultsRepository.transaction = transaction

	return nil
}

func RollbackTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction
	// откат выполняем в defer, поэтому перед этим м.б. выполнен коммит, который обнулит транзакцию
	if transaction == nil {
		return nil
	}

	runnersRepository.transaction = nil
	resultsRepository.transaction = nil

	return transaction.Rollback()
}

func CommitTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.transaction

	// на всякий случай
	if transaction == nil {
		return nil
	}

	runnersRepository.transaction = nil
	resultsRepository.transaction = nil

	return transaction.Commit()
}
