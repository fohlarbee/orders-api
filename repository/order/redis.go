package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fohlarbee/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert (ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	} 

	key := orderIDKey(order.OrderID)
	
	txn := r.Client.TxPipeline()

	res := txn.SetNX (ctx, key, string(data), 0)

	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add order to set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil { 
		 return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("order does not exist");

func (r * RedisRepo) GetByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	}else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return order, nil
}



func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	txn := r.Client.TxPipeline()
	err := txn.Del(ctx, key).Err()

	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove order from set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil

}	

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil){
		return ErrNotExist

	}else if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}


type FindAllPage struct {
	Size uint
	Offset uint
}

type FindResult struct {
	Orders []model.Order
	Cursor uint
}
func (r *RedisRepo) FinAll(ctx context.Context, page FindAllPage) (FindResult, error){
	 
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size) )
	keys, cursor, err := res.Result()	
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to scan orders: %w", err)
	}
	if len(keys) == 0 {
		return FindResult{
			Orders:  []model.Order{},

		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, 0, len(xs))

	for i, x := range xs {
		x := x.(string)  
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to unmarshal order at index %d: %w", i, err)	

		}
		orders[i] = order 


	}

	return FindResult{
		Orders:  orders,
		Cursor: uint(cursor),
	}, nil
}