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
	"google.golang.org/protobuf/types/known/wrapperspb"
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

// Interface guard to ensure that UsersService implements the UsersServiceServer interface
var _ service.UsersServiceServer = (*UsersService)(nil)

// UsersService implements the gRPC service for user management
type UsersService struct {
	service.UnimplementedUsersServiceServer
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
	if req.Input.Name.GetValue() != "" {
		user.Name = req.Input.Name.GetValue()
	}

	if req.Input.Email.GetValue() != "" {
		user.Email = req.Input.Email.GetValue()
	}

	// Update role if provided
	if req.Input.Role != service.UserRole_USER_ROLE_UNSPECIFIED {
		user.Role = req.Input.Role
	}

	if len(req.Input.Permissions.GetList().GetItems()) > 0 {
		user.Permissions = req.Input.Permissions.GetList().GetItems()
	}

	if len(req.Input.Tags.GetList().GetItems()) > 0 {
		user.Tags = &service.ListOfString{List: &service.ListOfString_List{Items: req.Input.Tags.GetList().GetItems()}}
	}

	// Update skill categories if provided
	if len(req.Input.SkillCategories.GetList().GetItems()) > 0 {
		user.SkillCategories = req.Input.SkillCategories
	}

	// Update bio if provided
	if req.Input.Bio.GetValue() != "" {
		user.Bio = req.Input.Bio
	}

	// Update age if provided
	if req.Input.Age != nil {
		user.Age = req.Input.Age
	}

	// Update profile if provided
	if req.Input.Profile != nil {
		if user.Profile == nil {
			user.Profile = &service.Profile{}
		}
		if req.Input.Profile.DisplayName.GetValue() != "" {
			user.Profile.DisplayName = req.Input.Profile.DisplayName
		}
		if req.Input.Profile.Timezone.GetValue() != "" {
			user.Profile.Timezone = req.Input.Profile.Timezone
		}
		if req.Input.Profile.Theme != service.Theme_THEME_UNSPECIFIED {
			user.Profile.Theme = req.Input.Profile.Theme
		}
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
		if input.Name.GetValue() != "" {
			user.Name = input.Name.GetValue()
		}

		if input.Email.GetValue() != "" {
			user.Email = input.Email.GetValue()
		}

		// Update role if provided
		if input.Role != service.UserRole_USER_ROLE_UNSPECIFIED {
			user.Role = input.Role
		}

		if len(input.Permissions.GetList().GetItems()) > 0 {
			user.Permissions = input.Permissions.GetList().GetItems()
		}

		if len(input.Tags.GetList().GetItems()) > 0 {
			user.Tags = &service.ListOfString{List: &service.ListOfString_List{Items: input.Tags.GetList().GetItems()}}
		}

		// Update skill categories if provided
		if len(input.SkillCategories.GetList().GetItems()) > 0 {
			user.SkillCategories = input.SkillCategories
		}

		// Update bio if provided
		if input.Bio.GetValue() != "" {
			user.Bio = input.Bio
		}

		// Update age if provided
		if input.Age != nil {
			user.Age = input.Age
		}

		// Update profile if provided
		if input.Profile != nil {
			if user.Profile == nil {
				user.Profile = &service.Profile{}
			}
			if input.Profile.DisplayName.GetValue() != "" {
				user.Profile.DisplayName = input.Profile.DisplayName
			}
			if input.Profile.Timezone.GetValue() != "" {
				user.Profile.Timezone = input.Profile.Timezone
			}
			if input.Profile.Theme != service.Theme_THEME_UNSPECIFIED {
				user.Profile.Theme = input.Profile.Theme
			}
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
			Phone:    &wrapperspb.StringValue{Value: user.Phone},
			Website:  &wrapperspb.StringValue{Value: user.Website},
			Company: &service.Company{
				Name:        user.Company.Name,
				CatchPhrase: &wrapperspb.StringValue{Value: user.Company.CatchPhrase},
				Bs:          &wrapperspb.StringValue{Value: user.Company.Bs},
			},
			Address: &service.Address{
				Street:  &wrapperspb.StringValue{Value: user.Address.Street},
				Suite:   &wrapperspb.StringValue{Value: user.Address.Suite},
				City:    &wrapperspb.StringValue{Value: user.Address.City},
				Zipcode: &wrapperspb.StringValue{Value: user.Address.Zipcode},
				Geo: &service.Geo{
					Lat: &wrapperspb.StringValue{Value: user.Address.Geo.Lat},
					Lng: &wrapperspb.StringValue{Value: user.Address.Geo.Lng},
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
		Phone:    &wrapperspb.StringValue{Value: user.Phone},
		Website:  &wrapperspb.StringValue{Value: user.Website},
		Company: &service.Company{
			Name:        user.Company.Name,
			CatchPhrase: &wrapperspb.StringValue{Value: user.Company.CatchPhrase},
			Bs:          &wrapperspb.StringValue{Value: user.Company.Bs},
		},
		Address: &service.Address{
			Street:  &wrapperspb.StringValue{Value: user.Address.Street},
			Suite:   &wrapperspb.StringValue{Value: user.Address.Suite},
			City:    &wrapperspb.StringValue{Value: user.Address.City},
			Zipcode: &wrapperspb.StringValue{Value: user.Address.Zipcode},
			Geo: &service.Geo{
				Lat: &wrapperspb.StringValue{Value: user.Address.Geo.Lat},
				Lng: &wrapperspb.StringValue{Value: user.Address.Geo.Lng},
			},
		},
	}

	// Set the external user in the response
	response.ExternalUser = serviceExternalUser

	return response, nil
}

// QueryUserActivity returns recent activity items for a user
func (s *UsersService) QueryUserActivity(ctx context.Context, req *service.QueryUserActivityRequest) (*service.QueryUserActivityResponse, error) {
	response := &service.QueryUserActivityResponse{}

	// Get activities for the user from our mock data
	activities, found := userActivityMap[req.UserId]
	if !found {
		// Return empty list if user not found
		response.UserActivity = []*service.ActivityItem{}
		return response, nil
	}

	// Apply limit if specified
	limit := int(req.Limit.GetValue())
	if limit == 0 || limit > len(activities) {
		limit = len(activities)
	}

	// Return the requested activities
	response.UserActivity = activities[:limit]
	return response, nil
}

// MutationCreatePost creates a new post and associates it with the author
func (s *UsersService) MutationCreatePost(ctx context.Context, req *service.MutationCreatePostRequest) (*service.MutationCreatePostResponse, error) {
	response := &service.MutationCreatePostResponse{}

	// Check if the author exists
	author, found := mockUsers[req.Input.AuthorId]
	if !found {
		return nil, fmt.Errorf("author with ID %s not found", req.Input.AuthorId)
	}

	// Generate a simple ID (in production, this would be from a database)
	newID := fmt.Sprintf("%d", len(mockPosts)+1)

	// Create the new post
	newPost := &service.Post{
		Id:       newID,
		Title:    req.Input.Title,
		AuthorId: req.Input.AuthorId,
	}

	// Add to our mock data
	mockPosts[newID] = newPost

	// Create an activity item for the new post
	newActivity := &service.ActivityItem{
		Value: &service.ActivityItem_Post{Post: newPost},
	}

	// Add to the author's recent activity (prepend to show most recent first)
	author.RecentActivity = append([]*service.ActivityItem{newActivity}, author.RecentActivity...)

	// Update the userActivityMap as well
	userActivityMap[req.Input.AuthorId] = author.RecentActivity

	// Update the user in our mock database
	mockUsers[req.Input.AuthorId] = author

	// Return the created post
	response.CreatePost = newPost
	return response, nil
}
