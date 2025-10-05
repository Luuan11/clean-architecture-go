package graph

import (
	"github.com/graphql-go/graphql"
	"github.com/luuan11/clean-architecture/internal/usecase"
)

type OrderResolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderResolver(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
) *OrderResolver {
	return &OrderResolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"tax": &graphql.Field{
			Type: graphql.Float,
		},
		"final_price": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

func (r *OrderResolver) CreateMutation() *graphql.Field {
	return &graphql.Field{
		Type: orderType,
		Args: graphql.FieldConfigArgument{
			"price": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"tax": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Float),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			price := p.Args["price"].(float64)
			tax := p.Args["tax"].(float64)

			input := usecase.CreateOrderInputDTO{
				Price: price,
				Tax:   tax,
			}

			output, err := r.CreateOrderUseCase.Execute(input)
			if err != nil {
				return nil, err
			}

			return output, nil
		},
	}
}

func (r *OrderResolver) ListQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(orderType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			orders, err := r.ListOrdersUseCase.Execute()
			if err != nil {
				return nil, err
			}
			return orders, nil
		},
	}
}

func (r *OrderResolver) BuildSchema() (graphql.Schema, error) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listOrders": r.ListQuery(),
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createOrder": r.CreateMutation(),
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}
