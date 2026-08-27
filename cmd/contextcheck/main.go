package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Contexto e propagação de erros")
	pedido := "Pedido A-001"

	ctx := context.Background()
	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// encerra o contexto quando a função terminar
	defer cancel()

	err := processarPedido(requestCtx, pedido)
	if err != nil {
		log.Printf("erro ao processar pedido: %v", err)
		log.Printf("requestCtx: %v", requestCtx.Err())
		return
	}

}

func processarPedido(ctx context.Context, pedido string) error {
	err := salvarPedido(ctx, pedido)
	if err != nil {
		return err
	}

	err = publicarEvento(ctx, pedido)
	if err != nil {
		return err
	}
	fmt.Println("pedido processado com sucesso")
	return nil
}

func salvarPedido(ctx context.Context, pedido string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		// sucesso
		fmt.Println("pedido salvo com sucesso", pedido)
		return nil
	case <-queryCtx.Done():
		// retornar ctxFilho.Err()
		return queryCtx.Err()
	}
}

func publicarEvento(ctx context.Context, pedido string) error {
	pubCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	select {
	case <-time.After(4 * time.Second):
		// sucesso
		fmt.Println("evento publicado com sucesso", pedido)
		return nil
	case <-pubCtx.Done():
		// retornar ctxFilho.Err()
		return pubCtx.Err()
	}
}
