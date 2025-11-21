package cmd

import (
	"fmt"
	"sync"
	"time"
)

// Possiveis estados do Circuit Breaker
type State int

const (
	Closed   State = iota // Fechado
	Open                  // Aberto
	HalfOpen              // Meio-aberto
)

type CircuitBreaker struct {
	failureThreshold int
	timeout          time.Duration
	// Estado
	mu           sync.Mutex
	currentState State
	failureCount int
	lastFailure  time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		timeout:          timeout,
		currentState:     Closed,
	}
}

// Função que envolve a chamada ao serviço externo
func (cb *CircuitBreaker) Execute(requestFunc func() (interface{}, error)) (interface{}, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.currentState {
	case Open:
		// O circuito está aberto: falhar imediatamente (fail-fast)
		if time.Since(cb.lastFailure) > cb.timeout {
			// Tentativa de transição para HalfOpen
			cb.currentState = HalfOpen
			// FIX: Sera?????
			return nil, fmt.Errorf("circuit still open, transitioning to half-open for a retry")
		} else {
			return nil, fmt.Errorf("circuit is open: external service is down")
		}

	case HalfOpen:
		// Permite apenas uma chamada para testar
		result, err := requestFunc()
		if err != nil {
			// Falhou no teste, volta para Open
			cb.currentState = Open
			cb.lastFailure = time.Now()
			cb.failureCount = 0 // Reinicia contagem se falhar aqui, ou define para o threshold
			return nil, err
		}
		// Sucesso no teste, volta para Closed
		cb.currentState = Closed
		cb.failureCount = 0
		return result, nil

	case Closed:
		// Executa a chamada normalmente
		result, err := requestFunc()
		if err != nil {
			// Contabiliza a falha e verifica o threshold
			cb.failureCount++
			if cb.failureCount >= cb.failureThreshold {
				// Abre o circuito
				cb.currentState = Open
				cb.lastFailure = time.Now()
			}
			return nil, err
		}
		// Chamada bem-sucedida, reseta a contagem
		cb.failureCount = 0
		return result, nil
	}

	// Case para HalfOpen que não foi capturado no switch (se for o caso)
	// Se a chamada de teste falhar, o código acima já trata.
	return nil, fmt.Errorf("unexpected circuit breaker state logic")
}
