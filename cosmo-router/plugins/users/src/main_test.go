package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	service "github.com/wundergraph/cosmo/plugin/generated"
	"github.com/wundergraph/cosmo/router-plugin/httpclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const bufSize = 1024 * 1024

// testService is a wrapper that holds the gRPC test components
type testService struct {
	grpcConn    *grpc.ClientConn
	usersClient service.UsersServiceClient
	cleanup     func()
}

// setupTestService creates a local gRPC server for testing
func setupTestService(t *testing.T) *testService {
	// Create a buffer for gRPC connections
	lis := bufconn.Listen(bufSize)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register our service
	service.RegisterUsersServiceServer(grpcServer, &UsersService{})

	// Start the server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create a client connection
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	// Create the service client
	client := service.NewUsersServiceClient(conn)

	// Return cleanup function
	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
	}

	return &testService{
		grpcConn:    conn,
		usersClient: client,
		cleanup:     cleanup,
	}
}

// setupExternalTestService creates a test service with a mock HTTP server for external API tests
func setupExternalTestService(t *testing.T) (*testService, *httptest.Server) {
	// Set up the basic service
	svc := setupTestService(t)

	// Create a mock HTTP server for external API tests
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[
				{
					"id": 1,
					"name": "Leanne Graham",
					"username": "Bret",
					"email": "Sincere@april.biz",
					"phone": "1-770-736-8031 x56442",
					"website": "hildegard.org",
					"address": {
						"street": "Kulas Light",
						"suite": "Apt. 556",
						"city": "Gwenborough",
						"zipcode": "92998-3874",
						"geo": {
							"lat": "-37.3159",
							"lng": "81.1496"
						}
					},
					"company": {
						"name": "Romaguera-Crona",
						"catchPhrase": "Multi-layered client-server neural-net",
						"bs": "harness real-time e-markets"
					}
				}
			]`))
		case "/users/1":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"id": 1,
				"name": "Leanne Graham",
				"username": "Bret",
				"email": "Sincere@april.biz",
				"phone": "1-770-736-8031 x56442",
				"website": "hildegard.org",
				"address": {
					"street": "Kulas Light",
					"suite": "Apt. 556",
					"city": "Gwenborough",
					"zipcode": "92998-3874",
					"geo": {
						"lat": "-37.3159",
						"lng": "81.1496"
					}
				},
				"company": {
					"name": "Romaguera-Crona",
					"catchPhrase": "Multi-layered client-server neural-net",
					"bs": "harness real-time e-markets"
				}
			}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	// Save the original HTTP client and replace with the mock
	originalClient := httpClient
	httpClient = httpclient.New(
		httpclient.WithBaseURL(mockServer.URL),
		httpclient.WithTimeout(5*time.Second),
	)

	// Update the cleanup function to also clean up the HTTP mock
	oldCleanup := svc.cleanup
	svc.cleanup = func() {
		oldCleanup()
		mockServer.Close()
		httpClient = originalClient // Restore the original client
	}

	return svc, mockServer
}

func TestLookupUserById(t *testing.T) {
	// Setup basic service - no need for HTTP mocks
	svc := setupTestService(t)
	defer svc.cleanup()

	tests := []struct {
		name    string
		ids     []string
		want    []*service.User
		wantErr bool
	}{
		{
			name: "valid users",
			ids:  []string{"1", "2"},
			want: []*service.User{
				{Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN},
				{Id: "2", Name: "Bob Smith", Email: "bob@example.com", Role: service.UserRole_USER_ROLE_USER},
			},
			wantErr: false,
		},
		{
			name:    "nonexistent user",
			ids:     []string{"999"},
			want:    []*service.User{{Id: "999"}},
			wantErr: false,
		},
		{
			name: "mixed valid and invalid",
			ids:  []string{"1", "999"},
			want: []*service.User{
				{Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN},
				{Id: "999"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.LookupUserByIdRequest{
				Keys: make([]*service.LookupUserByIdRequestKey, 0, len(tt.ids)),
			}
			for _, id := range tt.ids {
				req.Keys = append(req.Keys, &service.LookupUserByIdRequestKey{Id: id})
			}

			resp, err := svc.usersClient.LookupUserById(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(resp.Result))

			for i, want := range tt.want {
				if want.Name != "" { // Only check non-empty fields
					assert.Equal(t, want.Id, resp.Result[i].Id)
					assert.Equal(t, want.Name, resp.Result[i].Name)
					assert.Equal(t, want.Email, resp.Result[i].Email)
					assert.Equal(t, want.Role, resp.Result[i].Role)
				} else {
					// For non-existent users, just check ID
					assert.Equal(t, want.Id, resp.Result[i].Id)
				}
			}
		})
	}
}

func TestQueryUser(t *testing.T) {
	// Setup basic service - no need for HTTP mocks
	svc := setupTestService(t)
	defer svc.cleanup()

	tests := []struct {
		name    string
		id      string
		want    *service.User
		wantErr bool
	}{
		{
			name:    "valid user",
			id:      "1",
			want:    &service.User{Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN},
			wantErr: false,
		},
		{
			name:    "nonexistent user",
			id:      "999",
			want:    &service.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.QueryUserRequest{Id: tt.id}

			resp, err := svc.usersClient.QueryUser(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.want.Name != "" { // Only check non-empty fields for valid users
				assert.Equal(t, tt.want.Id, resp.User.Id)
				assert.Equal(t, tt.want.Name, resp.User.Name)
				assert.Equal(t, tt.want.Email, resp.User.Email)
				assert.Equal(t, tt.want.Role, resp.User.Role)
			} else {
				// For non-existent users, user should be nil
				assert.Nil(t, resp.User)
			}
		})
	}
}

func TestQueryUsers(t *testing.T) {
	// Setup basic service - no need for HTTP mocks
	svc := setupTestService(t)
	defer svc.cleanup()

	req := &service.QueryUsersRequest{}
	resp, err := svc.usersClient.QueryUsers(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, len(mockUsers), len(resp.Users))

	// Create a map to check if each user is in the response
	responseMap := make(map[string]*service.User)
	for _, user := range resp.Users {
		responseMap[user.Id] = user
	}

	// Verify each mock user is in the response
	for id, mockUser := range mockUsers {
		respUser, ok := responseMap[id]
		assert.True(t, ok, "User %s should be in the response", id)
		assert.Equal(t, mockUser.Id, respUser.Id)
		assert.Equal(t, mockUser.Name, respUser.Name)
		assert.Equal(t, mockUser.Email, respUser.Email)
		assert.Equal(t, mockUser.Role, respUser.Role)
	}
}

func TestQueryExternalUsers(t *testing.T) {
	// Setup service with HTTP mocks for external API
	svc, _ := setupExternalTestService(t)
	defer svc.cleanup()

	req := &service.QueryExternalUsersRequest{}
	resp, err := svc.usersClient.QueryExternalUsers(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, resp.ExternalUsers, 1)

	user := resp.ExternalUsers[0]
	assert.Equal(t, "1", user.Id)
	assert.Equal(t, "Leanne Graham", user.Name)
	assert.Equal(t, "Bret", user.Username)
	assert.Equal(t, "Sincere@april.biz", user.Email)

	// Check address
	assert.Equal(t, "Kulas Light", user.Address.Street.GetValue())
	assert.Equal(t, "Apt. 556", user.Address.Suite.GetValue())
	assert.Equal(t, "Gwenborough", user.Address.City.GetValue())
	assert.Equal(t, "92998-3874", user.Address.Zipcode.GetValue())
	assert.Equal(t, "-37.3159", user.Address.Geo.Lat.GetValue())
	assert.Equal(t, "81.1496", user.Address.Geo.Lng.GetValue())

	// Check company
	assert.Equal(t, "Romaguera-Crona", user.Company.Name)
	assert.Equal(t, "Multi-layered client-server neural-net", user.Company.CatchPhrase.GetValue())
	assert.Equal(t, "harness real-time e-markets", user.Company.Bs.GetValue())
}

func TestQueryExternalUser(t *testing.T) {
	// Setup service with HTTP mocks for external API
	svc, _ := setupExternalTestService(t)
	defer svc.cleanup()

	tests := []struct {
		name    string
		id      string
		want    *service.ExternalUser
		wantErr bool
	}{
		{
			name: "valid external user",
			id:   "1",
			want: &service.ExternalUser{
				Id:       "1",
				Name:     "Leanne Graham",
				Username: "Bret",
				Email:    "Sincere@april.biz",
			},
			wantErr: false,
		},
		{
			name:    "nonexistent external user",
			id:      "999",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.QueryExternalUserRequest{Id: tt.id}

			resp, err := svc.usersClient.QueryExternalUser(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp.ExternalUser)
			assert.Equal(t, tt.want.Id, resp.ExternalUser.Id)
			assert.Equal(t, tt.want.Name, resp.ExternalUser.Name)
			assert.Equal(t, tt.want.Username, resp.ExternalUser.Username)
			assert.Equal(t, tt.want.Email, resp.ExternalUser.Email)
		})
	}
}

func TestMutationUpdateUser(t *testing.T) {
	// Setup basic service - no need for HTTP mocks
	svc := setupTestService(t)
	defer svc.cleanup()

	// Store original user data to restore after test
	origUser := *mockUsers["1"]
	defer func() {
		mockUsers["1"] = &origUser
	}()

	tests := []struct {
		name    string
		input   *service.UserInput
		want    *service.User
		wantErr bool
	}{
		{
			name: "valid update - name only",
			input: &service.UserInput{
				Id:   "1",
				Name: &wrapperspb.StringValue{Value: "Alice Updated"},
			},
			want: &service.User{
				Id:    "1",
				Name:  "Alice Updated",
				Email: "alice@example.com",
				Role:  service.UserRole_USER_ROLE_ADMIN,
			},
			wantErr: false,
		},
		{
			name: "valid update - email only",
			input: &service.UserInput{
				Id:    "1",
				Email: &wrapperspb.StringValue{Value: "alice.updated@example.com"},
			},
			want: &service.User{
				Id:    "1",
				Name:  "Alice Updated", // Name from previous test
				Email: "alice.updated@example.com",
				Role:  service.UserRole_USER_ROLE_ADMIN,
			},
			wantErr: false,
		},
		{
			name: "valid update - role only",
			input: &service.UserInput{
				Id:   "1",
				Role: service.UserRole_USER_ROLE_USER,
			},
			want: &service.User{
				Id:    "1",
				Name:  "Alice Updated",             // Name from previous test
				Email: "alice.updated@example.com", // Email from previous test
				Role:  service.UserRole_USER_ROLE_USER,
			},
			wantErr: false,
		},
		{
			name: "nonexistent user",
			input: &service.UserInput{
				Id:   "999",
				Name: &wrapperspb.StringValue{Value: "Nonexistent User"},
			},
			want:    &service.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.MutationUpdateUserRequest{
				Input: tt.input,
			}

			resp, err := svc.usersClient.MutationUpdateUser(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.want.Name != "" { // Only check when we expect a valid response
				assert.NotNil(t, resp.UpdateUser)
				assert.Equal(t, tt.want.Id, resp.UpdateUser.Id)
				assert.Equal(t, tt.want.Name, resp.UpdateUser.Name)
				assert.Equal(t, tt.want.Email, resp.UpdateUser.Email)
				assert.Equal(t, tt.want.Role, resp.UpdateUser.Role)
			} else {
				// For nonexistent users, expect empty response
				assert.Nil(t, resp.UpdateUser)
			}
		})
	}
}

func TestMutationUpdateUsers(t *testing.T) {
	// Setup basic service - no need for HTTP mocks
	svc := setupTestService(t)
	defer svc.cleanup()

	// Store original user data to restore after test
	origUsers := make(map[string]*service.User)
	for id, user := range mockUsers {
		userCopy := *user
		origUsers[id] = &userCopy
	}
	defer func() {
		for id, user := range origUsers {
			mockUsers[id] = user
		}
	}()

	tests := []struct {
		name    string
		inputs  []*service.UserInput
		want    []*service.User
		wantErr bool
	}{
		{
			name: "update multiple valid users",
			inputs: []*service.UserInput{
				{
					Id:    "1",
					Name:  &wrapperspb.StringValue{Value: "Alice Batch Updated"},
					Email: &wrapperspb.StringValue{Value: "alice.batch@example.com"},
				},
				{
					Id:   "2",
					Name: &wrapperspb.StringValue{Value: "Bob Batch Updated"},
					Role: service.UserRole_USER_ROLE_ADMIN,
				},
			},
			want: []*service.User{
				{
					Id:    "1",
					Name:  "Alice Batch Updated",
					Email: "alice.batch@example.com",
					Role:  service.UserRole_USER_ROLE_ADMIN,
				},
				{
					Id:    "2",
					Name:  "Bob Batch Updated",
					Email: "bob@example.com",
					Role:  service.UserRole_USER_ROLE_ADMIN,
				},
			},
			wantErr: false,
		},
		{
			name: "mix of valid and invalid users",
			inputs: []*service.UserInput{
				{
					Id:   "3",
					Name: &wrapperspb.StringValue{Value: "Charlie Batch Updated"},
				},
				{
					Id:   "999", // Nonexistent user
					Name: &wrapperspb.StringValue{Value: "Nonexistent User"},
				},
			},
			want: []*service.User{
				{
					Id:    "3",
					Name:  "Charlie Batch Updated",
					Email: "charlie@example.com",
					Role:  service.UserRole_USER_ROLE_USER,
				},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			inputs:  []*service.UserInput{},
			want:    []*service.User{},
			wantErr: false,
		},
		{
			name: "missing ID",
			inputs: []*service.UserInput{
				{
					Name: &wrapperspb.StringValue{Value: "Missing ID User"},
				},
			},
			want:    []*service.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.MutationUpdateUsersRequest{
				Input: tt.inputs,
			}

			resp, err := svc.usersClient.MutationUpdateUsers(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(resp.UpdateUsers))

			// Create a map of expected users by ID for easier comparison
			expectedUsers := make(map[string]*service.User)
			for _, user := range tt.want {
				expectedUsers[user.Id] = user
			}

			// Check each returned user against expected values
			for _, updatedUser := range resp.UpdateUsers {
				expected, ok := expectedUsers[updatedUser.Id]
				assert.True(t, ok, "User %s should be in the expected results", updatedUser.Id)
				assert.Equal(t, expected.Id, updatedUser.Id)
				assert.Equal(t, expected.Name, updatedUser.Name)
				assert.Equal(t, expected.Email, updatedUser.Email)
				assert.Equal(t, expected.Role, updatedUser.Role)
			}
		})
	}
}
