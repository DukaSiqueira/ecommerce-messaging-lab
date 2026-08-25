# Evento VS Comando

Um comando expressa uma intenção:
> "Reserve o estoque do pedido 001."

Características:
* Geralmente possui um destinatário específico.
* Pode ser aceito ou rejeitado.
* Costuma usar o verbo no imperativo
* Exemplos: `ReserveStock`, `CancelOrder`, `SendNotification`.

- - -

Um evento comunica um fato que já aconteceu:
> "O pedido 123 foi realizado."

Características:
* Não exige que o publicador conheça os interessados.
* Pode ser consumido por zero, um ou vários serviços.
* É normalmente nomeado no passado.
* Exemplos: `OrderPlaced`, `StockReserved`, `PaymentApproved`.

No estudo, o fluxo conceitual será:
```
Serviço de Pedidos
    publica OrderPlaced
              |
              +--> Estoque decide reservar os produtos
              |
              +--> Notificação decide enviar a confirmação
```

O serviço de Pedidos não manda Estoque ou Notificação fazerem algo diretamente. Ele apenas informa um fato. Cada consumidor decide como reagir.

Dentro do Estoque, o evento `OrderPlaced` pode originar uma ação interna equivalente a `ReserveStock`. Essa separação reduz o acoplamento.

## Consistência eventual
Quando o pedido é criado, Estoque e Notificação não terminam necessariamente no mesmo instante.

Por alguns segundos, o sistema pode estar assim:
```
Pedido: criado
Estoque: reserva pendente
Notificação: ainda não enviada
```

Isso é consistência eventual: diferentes partes do sistema chegam ao estado esperado em momento distintos.

Em um e-commerce real, a interface não deveria afirmar imediatamente "estoque confirmado" se a reserva ainda está pendente. Poderia mostrar algo como "pedido recebido, aguardando confirmação".

