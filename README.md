Como rodar:

Como rodar os testes:

Exemplo das requests:
Você pode copiar os seguintes curls e colar em uma ferramenta como insomnia ou postman para ter as requests já pré montadas.
Create Order: `localhost:8080/orders`
```
curl --request POST \
  --url http://localhost:8080/orders \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.3.0' \
  --data '{
	"owner_order_id": "aab4d348-0c67-4796-b977-9e779b29499c",
	"price_order_brl": 1500,
	"price_order_bt": 8,
	"type_order": 8,
	"status": 1
}'
```

Update Status Order: `localhost:8080/orders/:id/status/:newStatus`
```
curl --request PATCH \
  --url http://localhost:8080/orders/eff91ed6-9a78-433e-aa80-d34a7507cc6d/status/4 \
  --header 'User-Agent: insomnia/11.3.0'
```

List Order:
```
curl --request GET \
  --url http://localhost:8080/orders \
  --header 'User-Agent: insomnia/11.3.0'
```

Get client:
```
curl --request GET \
  --url http://localhost:8080/client/aab4d348-0c67-4796-b977-9e779b29499c \
  --header 'User-Agent: insomnia/11.3.0'
```

Para realização de testes usando essas rotas já deixamos 