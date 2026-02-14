package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type TableClient struct {
	serviceClient *aztables.ServiceClient
}

func NewTableClient(connectionString string) (*TableClient, error) {
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service client: %w", err)
	}

	return &TableClient{
		serviceClient: serviceClient,
	}, nil
}

func (tc *TableClient) GetTableClient(tableName string) *aztables.Client {
	return tc.serviceClient.NewClient(tableName)
}

func (tc *TableClient) EnsureTable(ctx context.Context, tableName string) error {
	client := tc.GetTableClient(tableName)
	_, err := client.CreateTable(ctx, nil)
	if err != nil {
		// TODO - check errors
		return nil
	}
	return nil
}

func (tc *TableClient) DeleteTable(ctx context.Context, tableName string) error {
	client := tc.GetTableClient(tableName)
	_, err := client.Delete(ctx, nil)
	return err
}
