// Package main implements a Cosmo Router plugin for user management.
// It provides both local mock users and integration with external user data from JSONPlaceholder.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	service "github.com/wundergraph/cosmo/plugin/generated"

	routerplugin "github.com/wundergraph/cosmo/router-plugin"
	"github.com/wundergraph/cosmo/router-plugin/httpclient"
	"google.golang.org/grpc"
)

// httpClient is used to make external API requests to JSONPlaceholder
var httpClient *httpclient.Client

// main initializes and starts the router plugin service
func main() {
	pl, err := routerplugin.NewRouterPlugin(func(s *grpc.Server) {
		s.RegisterService(&service.UsersService_ServiceDesc, &UsersService{})
	})

	if err != nil {
		log.Fatalf("failed to create router plugin: %v", err)
	}

	// Initialize HTTP client for external API calls
	httpClient = httpclient.New(
		httpclient.WithBaseURL("https://jsonplaceholder.typicode.com"),
		httpclient.WithTimeout(5*time.Second),
	)

	pl.Serve()
}

// Geo represents geographic coordinates in the JSONPlaceholder API
type Geo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// Address represents an address in JSONPlaceholder API
type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     Geo    `json:"geo"`
}

// Company represents a company in JSONPlaceholder API
type Company struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
	Bs          string `json:"bs"`
}

// ExternalUser represents a user from the JSONPlaceholder API
type ExternalUser struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Address  Address `json:"address"`
	Company  Company `json:"company"`
}

// UsersService implements the gRPC service for user management
type UsersService struct {
	service.UnimplementedUsersServiceServer
}

// Mock user data store for demonstration purposes
// In a production environment, this would be replaced with a proper database
var mockUsers = map[string]*service.User{
	"1": {Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN},
	"2": {Id: "2", Name: "Bob Smith", Email: "bob@example.com", Role: service.UserRole_USER_ROLE_USER},
	"3": {Id: "3", Name: "Charlie Brown", Email: "charlie@example.com", Role: service.UserRole_USER_ROLE_USER},
	"4": {Id: "4", Name: "Dana Lee", Email: "dana@example.com", Role: service.UserRole_USER_ROLE_GUEST},
}

// LookupUserById implements the batch lookup functionality.
// It receives an array of user IDs and returns the corresponding user objects.
// This method is optimized for DataLoader patterns in GraphQL resolvers.
func (s *UsersService) LookupUserById(ctx context.Context, req *service.LookupUserByIdRequest) (*service.LookupUserByIdResponse, error) {
	response := &service.LookupUserByIdResponse{
		Result: make([]*service.User, 0, len(req.Keys)),
	}

	// Process each key in the batch request
	for _, key := range req.Keys {
		if user, found := mockUsers[key.Id]; found {
			response.Result = append(response.Result, user)
		} else {
			// Return nil or empty user for keys that don't exist
			response.Result = append(response.Result, &service.User{Id: key.Id})
		}
	}

	return response, nil
}

// QueryUser looks up a single user by ID.
// Returns the user if found, otherwise returns an empty response.
func (s *UsersService) QueryUser(ctx context.Context, req *service.QueryUserRequest) (*service.QueryUserResponse, error) {
	response := &service.QueryUserResponse{}

	if user, found := mockUsers[req.Id]; found {
		response.User = user
	}

	return response, nil
}

// QueryUsers returns all users from the mock data store.
// This method doesn't support pagination or filtering in this implementation.
func (s *UsersService) QueryUsers(ctx context.Context, req *service.QueryUsersRequest) (*service.QueryUsersResponse, error) {
	response := &service.QueryUsersResponse{
		Users: make([]*service.User, 0, len(mockUsers)),
	}

	for _, user := range mockUsers {
		response.Users = append(response.Users, user)
	}

	return response, nil
}

// MutationUpdateUser updates a user's information.
// Only updates fields that are provided in the input.
// Returns the updated user if found, otherwise returns an empty response.
func (s *UsersService) MutationUpdateUser(ctx context.Context, req *service.MutationUpdateUserRequest) (*service.MutationUpdateUserResponse, error) {
	response := &service.MutationUpdateUserResponse{}

	// Check if user exists
	user, found := mockUsers[req.Input.Id]
	if !found {
		return response, nil
	}

	// Update user fields if provided in the input
	if req.Input.Name != "" {
		user.Name = req.Input.Name
	}

	if req.Input.Email != "" {
		user.Email = req.Input.Email
	}

	// Update role if provided
	if req.Input.Role != service.UserRole_USER_ROLE_UNSPECIFIED {
		user.Role = req.Input.Role
	}

	// Update the user in our mock database
	mockUsers[req.Input.Id] = user

	// Return the updated user
	response.UpdateUser = user

	return response, nil
}

// MutationUpdateUsers updates multiple users' information in a batch.
// Only updates fields that are provided in each input.
// Returns the updated users that were found, skipping any that don't exist.
func (s *UsersService) MutationUpdateUsers(ctx context.Context, req *service.MutationUpdateUsersRequest) (*service.MutationUpdateUsersResponse, error) {
	response := &service.MutationUpdateUsersResponse{
		UpdateUsers: make([]*service.User, 0, len(req.Input)),
	}

	// Process each user update in the batch request
	for _, input := range req.Input {
		// Skip if no ID is provided
		if input.Id == "" {
			continue
		}

		// Check if user exists
		user, found := mockUsers[input.Id]
		if !found {
			continue
		}

		// Update user fields if provided in the input
		if input.Name != "" {
			user.Name = input.Name
		}

		if input.Email != "" {
			user.Email = input.Email
		}

		// Update role if provided
		if input.Role != service.UserRole_USER_ROLE_UNSPECIFIED {
			user.Role = input.Role
		}

		// Update the user in our mock database
		mockUsers[input.Id] = user

		// Add the updated user to the response
		response.UpdateUsers = append(response.UpdateUsers, user)
	}

	return response, nil
}

// QueryExternalUsers fetches users from the JSONPlaceholder API.
// It demonstrates integration with an external REST API.
func (s *UsersService) QueryExternalUsers(ctx context.Context, req *service.QueryExternalUsersRequest) (*service.QueryExternalUsersResponse, error) {
	response := &service.QueryExternalUsersResponse{}

	// Use the httpClient to make the request
	resp, err := httpClient.Get(ctx, "/users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external users: %w", err)
	}

	// Unmarshal the JSON response into our data structure
	externalUsers, err := httpclient.UnmarshalTo[[]ExternalUser](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal external users: %w", err)
	}

	// Convert to service.ExternalUser objects
	serviceExternalUsers := make([]*service.ExternalUser, 0, len(externalUsers))
	for _, user := range externalUsers {
		serviceExternalUser := &service.ExternalUser{
			Id:       fmt.Sprintf("%d", user.ID),
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Phone:    user.Phone,
			Website:  user.Website,
			Company: &service.Company{
				Name:        user.Company.Name,
				CatchPhrase: user.Company.CatchPhrase,
				Bs:          user.Company.Bs,
			},
			Address: &service.Address{
				Street:  user.Address.Street,
				Suite:   user.Address.Suite,
				City:    user.Address.City,
				Zipcode: user.Address.Zipcode,
				Geo: &service.Geo{
					Lat: user.Address.Geo.Lat,
					Lng: user.Address.Geo.Lng,
				},
			},
		}
		serviceExternalUsers = append(serviceExternalUsers, serviceExternalUser)
	}

	// Set the external users in the response
	response.ExternalUsers = serviceExternalUsers

	return response, nil
}

// QueryExternalUser fetches a single user by ID from the JSONPlaceholder API.
// It demonstrates how to fetch a specific resource from an external REST API.
func (s *UsersService) QueryExternalUser(ctx context.Context, req *service.QueryExternalUserRequest) (*service.QueryExternalUserResponse, error) {
	response := &service.QueryExternalUserResponse{}

	// Use the httpClient to make the request for a specific user
	resp, err := httpClient.Get(ctx, fmt.Sprintf("/users/%s", req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external user: %w", err)
	}

	// Unmarshal the JSON response into our data structure
	user, err := httpclient.UnmarshalTo[ExternalUser](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal external user: %w", err)
	}

	// Convert to service.ExternalUser
	serviceExternalUser := &service.ExternalUser{
		Id:       fmt.Sprintf("%d", user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
		Website:  user.Website,
		Company: &service.Company{
			Name:        user.Company.Name,
			CatchPhrase: user.Company.CatchPhrase,
			Bs:          user.Company.Bs,
		},
		Address: &service.Address{
			Street:  user.Address.Street,
			Suite:   user.Address.Suite,
			City:    user.Address.City,
			Zipcode: user.Address.Zipcode,
			Geo: &service.Geo{
				Lat: user.Address.Geo.Lat,
				Lng: user.Address.Geo.Lng,
			},
		},
	}

	// Set the external user in the response
	response.ExternalUser = serviceExternalUser

	return response, nil
}
