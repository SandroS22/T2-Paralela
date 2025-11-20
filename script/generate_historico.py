import json
import random
from datetime import datetime, timedelta

# --- Configurações ---
NUM_TRANSACTIONS = 150
PRODUCT_IDS = ["APL_2025", "GOOG_2025", "MSFT_2025", "TSLA_2025", "AMZN_2025"]
START_DATE = datetime(2025, 11, 1, 9, 0, 0) # Começa em 1 de Novembro
OUTPUT_FILENAME = "../data/historical_data.json"

def generate_transactions(num):
    transactions = []
    current_time = START_DATE
    
    # Preços base para os ativos
    base_prices = {
        "APL_2025": 170.00,
        "GOOG_2025": 2850.00,
        "MSFT_2025": 300.00,
        "TSLA_2025": 900.00,
        "AMZN_2025": 3200.00
    }

    for i in range(1, num + 1):
        product_id = random.choice(PRODUCT_IDS)
        
        # Simula o avanço do tempo (entre 1 e 10 minutos por transação)
        current_time += timedelta(minutes=random.randint(1, 10), seconds=random.randint(0, 59))
        
        # Simula a variação do preço base (máximo 1% de variação)
        price_variation = base_prices[product_id] * random.uniform(-0.01, 0.01)
        price = round(base_prices[product_id] + price_variation, 2)
        
        transaction = {
            "transaction_id": f"T{i:06d}",
            "product_id": product_id,
            "timestamp": current_time.isoformat() + 'Z', # Formato ISO 8601 com 'Z' (UTC)
            "price": price,
            "volume": random.choice([10, 50, 100, 200, 500]),
            "transaction_type": random.choice(["BUY", "SELL"])
        }
        transactions.append(transaction)
        
        # Atualiza o preço base para a próxima transação, simulando movimento de mercado
        base_prices[product_id] = price

    return transactions

# --- Execução ---
historical_data = generate_transactions(NUM_TRANSACTIONS)

with open(OUTPUT_FILENAME, 'w') as f:
    json.dump(historical_data, f, indent=2)

print(f"Arquivo '{OUTPUT_FILENAME}' gerado com sucesso!")
print(f"Total de {len(historical_data)} transações criadas.")
