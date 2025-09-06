package worker

import (
	"context"
	"math/big"
	"order/internal/domain"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type PaymentWatcher struct {
	orderRepo domain.OrderRepository
	ethClient *ethclient.Client
	logger    *zap.Logger
	ctx       context.Context
}

func (w *PaymentWatcher) Start() {
	w.logger.Info("Payment watched started...")
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			w.logger.Info("Checking for pending payment...")
			w.processPendingOrders()
		case <-w.ctx.Done():
			w.logger.Info("Payment watcher shutting down.")
			return
		}
	}
}

func (w *PaymentWatcher) processPendingOrders() {
	limit := 100
	offset := 0

	for {
		orders, err := w.orderRepo.FindAwaitingPayment(w.ctx, limit, offset)
		if err != nil {
			w.logger.Error("Failed to fetch awaiting payment orders", zap.Error(err))
			return
		}

		if len(orders) == 0 {
			break
		}

		for _, order := range orders {
			w.checkOrderPayment(order)
		}

		offset += limit
	}
}

func (w *PaymentWatcher) checkOrderPayment(order *domain.Order) {
	if order.PaymentAddress == nil {
		return
	}

	address := common.HexToAddress(*order.PaymentAddress)
	balance, err := w.ethClient.BalanceAt(w.ctx, address, nil)
	if err != nil {
		w.logger.Error("Failed to get balance for address",
			zap.String("address", *order.PaymentAddress),
			zap.Error(err),
		)
		return
	}

	if balance.Cmp(big.NewInt(0)) > 0 {
		w.logger.Info("Payment detected!",
			zap.String("order_id", string(order.ID)),
			zap.String("address", *order.PaymentAddress),
			zap.String("balance", balance.String()),
		)

		// TODO: Find the actual transaction hash and save it. This is a more advanced step.

		order.Status = domain.StatusPaid
		err := w.orderRepo.Update(w.ctx, order)
		if err != nil {
			w.logger.Error("Failed to update order status to PAID", zap.String("order_id", string(order.ID)), zap.Error(err))
		} else {
			w.logger.Info("Order status updated to PAID", zap.String("order_id", string(order.ID)))
		}
	}
}

func (w *PaymentWatcher) Close() {
	w.ethClient.Close()
}
