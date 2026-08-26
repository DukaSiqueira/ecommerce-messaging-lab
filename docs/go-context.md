# Go Context
## Visão geral
Abaixo uma descrição sobre `context` um pacote da biblioteca padrão do Go. `context` oferece uma interface para consultar sinais relacionados ao ciclo de vida das operações.

- - -

## Qual problema o Context resolve?
Pense no cenário de um e-commerce:
```text
Cliente
  → API de Pedidos
    → Caso de uso
      → Banco de dados
      → Serviço externo
      → Pub/Sub
```

Imagine que: 
1. O cliente inicia uma compra.
2. A API começa a trabalhar.
3. O banco fica lento.
4. O cliente fecha a conexão ou o servidor ultrapassa seu prazo.
5. As operações internas continuam consumindo conexões, memória e CPU.

O `context` oferece um sinal comum que pode percorrer toda essa cadeia.

Quando a operação principal perde a utilidade, o cancelamento pode ser propagado para as operações filhas:
```
Requisição cancelada
  → caso de uso recebe o sinal
    → repositório recebe o sinal
    → cliente HTTP recebe o sinal
    → chamada Pub/Sub recebe o sinal
```
Essa é uma técnica para liberar recursos que não estão mais sendo utilizados. O Go recomenda propagar o contexto entre as funções envolvidas na mesma operação.

> [!info] Cancelar um Context não desfazer uma operação que já foi concluída.

- - -

## O que Context é?
`Context` é um valor comum de Go, passado entre funções, com outros parâmetros.
```Go
func ProcessarPedido(ctx context.Context, pedido Pedido) error
```
* `pedido` diz o que processar.
* `ctx` informa as condições de vida desse processamento.

`context` é um sinal propagável, o chamador passa o contexto para função seguinte como parâmetro:
```
Handler recebe ctx
   ↓
Caso de uso recebe ctx
   ↓
Repositório recebe ctx
```
A regra geral é que a função não deve descartar o contexto recebido e criar `context.Background()` no meio da cadeia. Isso quebraria a propagação do cancelamento.

`context` funciona como uma árvore, podendo gerar contextos filhos:
```
Contexto da requisição
├── Contexto da consulta ao banco
└── Contexto da publicação no Pub/Sub
```
Se o contexto pai for cancelado, os filhos também são cancelados. Um filho pode ter um prazo mais curto que o pai, mas não consegue prolongar a vida além do pai.

`Contexto` é um mecanismo cooperativo, ele não força o encerramento de uma função, ele apenas comunica que uma operação foi cancelada e a responsabilidade de respeitar esse sinal é da função ou biblioteca. Bibliotecas como banco de dados, HTTP e Pub/Sub normalmente já recebem `Context`.

- - - 

## O que Context não é?
Context não é:
- O conteúdo do pedido.
- O evento `OrderPlaced`.
- A definição do tópico.
- O cliente Pub/Sub.
- A conexão de rede.
- Uma transação de banco.
- Um mecanismo automático de rollback.
- Um lugar genérico para colocar qualquer informação.
- Um objeto que interrompe qualquer código à força.
- Um substituto para parâmetros de função.
- Um container de injeção de dependências.

Uma comparação útil:

```
Pedido          → o que deve ser processado
Cliente Pub/Sub → onde executar a operação
Context         → até quando a operação ainda deve continuar
Resultado       → o que aconteceu
```

## Aplicação ao código atual

Na chamada:

```
CreateTopic(operationCtx, topicRequest)
```

Você pode interpretar assim:

```
topicRequest
→ “Quero criar este tópico.”

operationCtx
→ “Esta tentativa pode trabalhar por até cinco segundos.”

CreateTopic
→ “Envie a solicitação ao Pub/Sub.”

createdTopic ou err
→ “Este foi o resultado.”
```

O Context não sabe que existe um tópico. Ele poderia ser usado da mesma maneira para:
- Consultar um pedido no banco.
- Chamar uma API de pagamento.
- Publicar uma mensagem.
- Aguardar uma tarefa concorrente.
- Criar um tópico

- - - 

## Por que chamar cancel quando já existe um timeout?
`context.WithTimeout` cria, entre outras coisas, um temporizador. Se a operação terminar em um segundo, não precisamos manter esse recurso aguardando os outros quatro segundos.

```
defer cancel()
```

não executa `cancel()` imediatamente. Ele registra:

> “Quando esta função terminar, execute `cancel()`.”

Isso garante a liberação mesmo se a função sair antecipadamente por algum `return`.

Uma precisão importante: o `defer` executa no final da **função atual**, não exatamente no final de `CreateTopic`. No seu `main`, essas duas coisas estão próximas. Em uma função muito longa, talvez fosse necessário delimitar melhor o contexto.

E `cancel()`:
- Libera os recursos associados ao contexto.
- Avisa operações ainda em andamento.
- Não cancela o contexto pai.
- Não transforma uma operação já concluída em falha.
- Não desfaz efeitos já realizados.