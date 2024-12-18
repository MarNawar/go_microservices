package main

import "context"

type queryResolver struct{
	server *Server
}

// Accounts


func (r *queryResolver)Accounts(ctx context.Context, pagination *PaginationInput, id *string)([]*Account, error){

}

// Products

func (r *queryResolver)Products(ctx context.Context, pagination *PaginationInput, query *string)([]*Product, error){
	
}