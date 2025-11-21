package main

import (
	"circuit-breaker/cmd"
	"fmt"
	"time"
)

// SuccessfulRequest simula uma chamada bem-sucedida ao serviço externo.
// Retorna o preço de cotação e nil (sem erro).
func SuccessfulRequest() (interface{}, error) {
	// Simula uma pequena latência normal
	time.Sleep(50 * time.Millisecond)

	// Retorna um valor de cotação simulado
	return 185.75, nil
}

// FailingRequest simula uma falha total no serviço externo.
// Retorna nil e um erro para o Circuit Breaker capturar.
func FailingRequest() (interface{}, error) {
	// Simula o tempo que levaria até o timeout ou a falha
	time.Sleep(10 * time.Second)

	// Retorna nil e o erro
	return nil, fmt.Errorf("external service failed: 503 Service Unavailable")
}

func main() {
	for {
		cb := cmd.NewCircuitBreaker(5, time.Second*5)
		price, err := cb.Execute(SuccessfulRequest)
		if err == nil {
			// SUCESSO: Circuito Fechado ou HalfOpen com sucesso
			fmt.Printf("[%s] SUCESSO. Preço: %.2f. Publicando no Pub/Sub...\n", "Fechado", price.(float64))
			// Lógica de Pub/Sub: Publisher.Publish("prices", price)
		} else {
			// FALHA: Circuito Aberto ou Falha no Fechado/HalfOpen
			fmt.Printf("[%s] FALHA. Erro: %v\n", "Aberto", err)

			// Se falhou, o Circuit Breaker já atualizou seu estado interno.
			// O próximo ciclo do loop verificará o novo estado.
		}

		// 2. Controlar a taxa de tentativas/processamento
		// Espera um tempo curto antes de tentar novamente.
		// Se o circuito estiver OPEN, o Execute fará a checagem do timeout.
		time.Sleep(1 * time.Second)
	}
}
