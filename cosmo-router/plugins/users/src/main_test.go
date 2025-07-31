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

// verifyActivityContent is a helper function to verify that actual activities match expected activities
func verifyActivityContent(t *testing.T, expectedActivities []*service.ActivityItem, actualActivities []*service.ActivityItem, maxItems int) {
	t.Helper()
	if maxItems == 0 || maxItems > len(expectedActivities) {
		maxItems = len(expectedActivities)
	}

	assert.Equal(t, maxItems, len(actualActivities))

	for j, expectedActivity := range expectedActivities {
		if j >= maxItems {
			break
		}
		if j < len(actualActivities) {
			actualActivity := actualActivities[j]
			if expectedActivity.GetPost() != nil {
				assert.NotNil(t, actualActivity.GetPost())
				assert.Nil(t, actualActivity.GetComment())
				assert.Equal(t, expectedActivity.GetPost().Id, actualActivity.GetPost().Id)
				assert.Equal(t, expectedActivity.GetPost().Title, actualActivity.GetPost().Title)
			}
			if expectedActivity.GetComment() != nil {
				assert.NotNil(t, actualActivity.GetComment())
				assert.Nil(t, actualActivity.GetPost())
				assert.Equal(t, expectedActivity.GetComment().Id, actualActivity.GetComment().Id)
				assert.Equal(t, expectedActivity.GetComment().Content, actualActivity.GetComment().Content)
			}
		}
	}
}

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
			t.Errorf("failed to serve: %v", err)
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
				{
					Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN,
					Permissions: []string{"read", "write"}, Tags: &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
						{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
						{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
						Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
						Theme:       service.Theme_THEME_DARK,
					},
					Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
					Age: &wrapperspb.Int32Value{Value: 28},
				},
				{
					Id: "2", Name: "Bob Smith", Email: "bob@example.com", Role: service.UserRole_USER_ROLE_USER,
					Permissions: []string{"read"}, Tags: &service.ListOfString{List: &service.ListOfString_List{Items: []string{"user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"Python", "Java"}}},
						{List: &service.ListOfString_List{Items: []string{"Django", "Spring"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						DisplayName: &wrapperspb.StringValue{Value: "Bob"},
						Timezone:    &wrapperspb.StringValue{Value: "Europe/London"},
						Theme:       service.Theme_THEME_LIGHT,
					},
					Bio: &wrapperspb.StringValue{Value: "Backend developer passionate about clean code"},
					Age: &wrapperspb.Int32Value{Value: 32},
				},
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
				{
					Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN,
					Permissions: []string{"read", "write"}, Tags: &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
						{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
						{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
						Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
						Theme:       service.Theme_THEME_DARK,
					},
					Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
					Age: &wrapperspb.Int32Value{Value: 28},
				},
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
					assert.Equal(t, want.Permissions, resp.Result[i].Permissions)
					if want.Tags != nil {
						assert.Equal(t, want.Tags.GetList().GetItems(), resp.Result[i].Tags.GetList().GetItems())
					}
					if want.SkillCategories != nil {
						assert.Equal(t, len(want.SkillCategories.GetList().GetItems()), len(resp.Result[i].SkillCategories.GetList().GetItems()))
					}
					// Check RecentActivity against actual mock data
					expectedActivities := userActivityMap[want.Id]
					verifyActivityContent(t, expectedActivities, resp.Result[i].RecentActivity, 0)
					if want.Profile != nil {
						assert.Equal(t, want.Profile.GetDisplayName(), resp.Result[i].Profile.GetDisplayName())
						assert.Equal(t, want.Profile.GetTimezone(), resp.Result[i].Profile.GetTimezone())
						assert.Equal(t, want.Profile.Theme, resp.Result[i].Profile.Theme)
					}
					assert.Equal(t, want.Bio.GetValue(), resp.Result[i].Bio.GetValue())
					assert.Equal(t, want.Age.GetValue(), resp.Result[i].Age.GetValue())
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
			name: "valid user",
			id:   "1",
			want: &service.User{
				Id: "1", Name: "Alice Johnson", Email: "alice@example.com", Role: service.UserRole_USER_ROLE_ADMIN,
				Permissions: []string{"read", "write"}, Tags: &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
				SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				}}},
				RecentActivity: []*service.ActivityItem{},
				Profile: &service.Profile{
					DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
					Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
					Theme:       service.Theme_THEME_DARK,
				},
				Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
				Age: &wrapperspb.Int32Value{Value: 28},
			},
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
				assert.Equal(t, tt.want.Permissions, resp.User.Permissions)
				if tt.want.Tags != nil {
					assert.Equal(t, tt.want.Tags.GetList().GetItems(), resp.User.Tags.GetList().GetItems())
				}
				if tt.want.SkillCategories != nil {
					assert.Equal(t, len(tt.want.SkillCategories.GetList().GetItems()), len(resp.User.SkillCategories.GetList().GetItems()))
				}
				// Check RecentActivity against actual mock data
				expectedActivities := userActivityMap[tt.want.Id]
				verifyActivityContent(t, expectedActivities, resp.User.RecentActivity, 0)
				if tt.want.Profile != nil {
					assert.Equal(t, tt.want.Profile.GetDisplayName(), resp.User.Profile.GetDisplayName())
					assert.Equal(t, tt.want.Profile.GetTimezone(), resp.User.Profile.GetTimezone())
					assert.Equal(t, tt.want.Profile.Theme, resp.User.Profile.Theme)
				}
				assert.Equal(t, tt.want.Bio.GetValue(), resp.User.Bio.GetValue())
				assert.Equal(t, tt.want.Age.GetValue(), resp.User.Age.GetValue())
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
		assert.Equal(t, mockUser.Permissions, respUser.Permissions)
		if mockUser.Tags != nil {
			assert.Equal(t, mockUser.Tags.GetList().GetItems(), respUser.Tags.GetList().GetItems())
		}
		if mockUser.SkillCategories != nil {
			assert.Equal(t, len(mockUser.SkillCategories.GetList().GetItems()), len(respUser.SkillCategories.GetList().GetItems()))
		}
		// Check RecentActivity against actual mock data
		expectedActivities := userActivityMap[mockUser.Id]
		verifyActivityContent(t, expectedActivities, respUser.RecentActivity, 0)
		if mockUser.Profile != nil {
			assert.Equal(t, mockUser.Profile.DisplayName.GetValue(), respUser.Profile.DisplayName.GetValue())
			assert.Equal(t, mockUser.Profile.Timezone.GetValue(), respUser.Profile.Timezone.GetValue())
			assert.Equal(t, mockUser.Profile.Theme, respUser.Profile.Theme)
		}
		assert.Equal(t, mockUser.Bio.GetValue(), respUser.Bio.GetValue())
		assert.Equal(t, mockUser.Age.GetValue(), respUser.Age.GetValue())
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

	mu := mockUsers["1"]

	// Store original user data to restore after test
	origUser := &service.User{
		Id:              mu.Id,
		Name:            mu.Name,
		Email:           mu.Email,
		Role:            mu.Role,
		Permissions:     mu.Permissions,
		Tags:            mu.Tags,
		SkillCategories: mu.SkillCategories,
		RecentActivity:  mu.RecentActivity,
		Profile:         mu.Profile,
		Bio:             mu.Bio,
		Age:             mu.Age,
	}

	defer func() {
		mockUsers["1"] = origUser
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
				Id:          "1",
				Name:        "Alice Updated",
				Email:       "alice@example.com",
				Role:        service.UserRole_USER_ROLE_ADMIN,
				Permissions: []string{"read", "write"},
				Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
				SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				}}},
				RecentActivity: []*service.ActivityItem{},
				Profile: &service.Profile{
					DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
					Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
					Theme:       service.Theme_THEME_DARK,
				},
				Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
				Age: &wrapperspb.Int32Value{Value: 28},
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
				Id:          "1",
				Name:        "Alice Updated",
				Email:       "alice.updated@example.com",
				Role:        service.UserRole_USER_ROLE_ADMIN,
				Permissions: []string{"read", "write"},
				Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
				SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				}}},
				RecentActivity: []*service.ActivityItem{},
				Profile: &service.Profile{
					DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
					Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
					Theme:       service.Theme_THEME_DARK,
				},
				Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
				Age: &wrapperspb.Int32Value{Value: 28},
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
				Id:          "1",
				Name:        "Alice Updated",
				Email:       "alice.updated@example.com",
				Role:        service.UserRole_USER_ROLE_USER,
				Permissions: []string{"read", "write"},
				Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
				SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				}}},
				RecentActivity: []*service.ActivityItem{},
				Profile: &service.Profile{
					DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
					Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
					Theme:       service.Theme_THEME_DARK,
				},
				Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
				Age: &wrapperspb.Int32Value{Value: 28},
			},
			wantErr: false,
		},
		{
			name: "valid update - permissions and tags",
			input: &service.UserInput{
				Id:          "1",
				Permissions: &service.ListOfString{List: &service.ListOfString_List{Items: []string{"read", "write", "delete"}}},
				Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"super-admin", "developer"}}},
			},
			want: &service.User{
				Id:          "1",
				Name:        "Alice Updated",
				Email:       "alice.updated@example.com",
				Role:        service.UserRole_USER_ROLE_USER,
				Permissions: []string{"read", "write", "delete"},
				Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"super-admin", "developer"}}},
				SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				}}},
				RecentActivity: []*service.ActivityItem{},
				Profile: &service.Profile{
					DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
					Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
					Theme:       service.Theme_THEME_DARK,
				},
				Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
				Age: &wrapperspb.Int32Value{Value: 28},
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
				assert.Equal(t, tt.want.Permissions, resp.UpdateUser.Permissions)
				if tt.want.Tags != nil {
					assert.Equal(t, tt.want.Tags.GetList().GetItems(), resp.UpdateUser.Tags.GetList().GetItems())
				}
				if tt.want.SkillCategories != nil {
					assert.Equal(t, len(tt.want.SkillCategories.GetList().GetItems()), len(resp.UpdateUser.SkillCategories.GetList().GetItems()))
				}
				// Check RecentActivity against actual mock data
				expectedActivities := userActivityMap[tt.want.Id]
				verifyActivityContent(t, expectedActivities, resp.UpdateUser.RecentActivity, 0)
				if tt.want.Profile != nil {
					assert.Equal(t, tt.want.Profile.GetDisplayName(), resp.UpdateUser.Profile.GetDisplayName())
					assert.Equal(t, tt.want.Profile.GetTimezone(), resp.UpdateUser.Profile.GetTimezone())
					assert.Equal(t, tt.want.Profile.Theme, resp.UpdateUser.Profile.Theme)
				}
				assert.Equal(t, tt.want.Bio.GetValue(), resp.UpdateUser.Bio.GetValue())
				assert.Equal(t, tt.want.Age.GetValue(), resp.UpdateUser.Age.GetValue())
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
		userCopy := &service.User{
			Id:              user.Id,
			Name:            user.Name,
			Email:           user.Email,
			Role:            user.Role,
			Permissions:     user.Permissions,
			Tags:            user.Tags,
			SkillCategories: user.SkillCategories,
			RecentActivity:  user.RecentActivity,
			Profile:         user.Profile,
			Bio:             user.Bio,
			Age:             user.Age,
		}
		origUsers[id] = userCopy
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
					Id:          "1",
					Name:        "Alice Batch Updated",
					Email:       "alice.batch@example.com",
					Role:        service.UserRole_USER_ROLE_ADMIN,
					Permissions: []string{"read", "write"}, // Original permissions from mockUsers
					Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
						{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
						{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
						Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
						Theme:       service.Theme_THEME_DARK,
					},
					Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
					Age: &wrapperspb.Int32Value{Value: 28},
				},
				{
					Id:          "2",
					Name:        "Bob Batch Updated",
					Email:       "bob@example.com",
					Role:        service.UserRole_USER_ROLE_ADMIN,
					Permissions: []string{"read"},
					Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"Python", "Java"}}},
						{List: &service.ListOfString_List{Items: []string{"Django", "Spring"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						DisplayName: &wrapperspb.StringValue{Value: "Bob"},
						Timezone:    &wrapperspb.StringValue{Value: "Europe/London"},
						Theme:       service.Theme_THEME_LIGHT,
					},
					Bio: &wrapperspb.StringValue{Value: "Backend developer passionate about clean code"},
					Age: &wrapperspb.Int32Value{Value: 32},
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
					Id:          "3",
					Name:        "Charlie Batch Updated",
					Email:       "charlie@example.com",
					Role:        service.UserRole_USER_ROLE_USER,
					Permissions: []string{"read"},
					Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"user"}}},
					SkillCategories: &service.ListOfListOfString{List: &service.ListOfListOfString_List{Items: []*service.ListOfString{
						{List: &service.ListOfString_List{Items: []string{"Go", "Rust"}}},
						{List: &service.ListOfString_List{Items: []string{"Docker", "Kubernetes"}}},
					}}},
					RecentActivity: []*service.ActivityItem{},
					Profile: &service.Profile{
						Timezone: &wrapperspb.StringValue{Value: "Asia/Tokyo"},
						Theme:    service.Theme_THEME_AUTO,
					},
					Age: &wrapperspb.Int32Value{Value: 29},
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
				assert.Equal(t, expected.Permissions, updatedUser.Permissions)
				if expected.Tags != nil {
					assert.Equal(t, expected.Tags.GetList().GetItems(), updatedUser.Tags.GetList().GetItems())
				}
				if expected.SkillCategories != nil {
					assert.Equal(t, len(expected.SkillCategories.GetList().GetItems()), len(updatedUser.SkillCategories.GetList().GetItems()))
				}
				// Check RecentActivity against actual mock data
				expectedActivities := userActivityMap[expected.Id]
				verifyActivityContent(t, expectedActivities, updatedUser.RecentActivity, 0)
				if expected.Profile != nil {
					assert.Equal(t, expected.Profile.GetDisplayName(), updatedUser.Profile.GetDisplayName())
					assert.Equal(t, expected.Profile.GetTimezone(), updatedUser.Profile.GetTimezone())
					assert.Equal(t, expected.Profile.Theme, updatedUser.Profile.Theme)
				}
				assert.Equal(t, expected.Bio.GetValue(), updatedUser.Bio.GetValue())
				assert.Equal(t, expected.Age.GetValue(), updatedUser.Age.GetValue())
			}
		})
	}
}

func TestQueryUserActivity(t *testing.T) {
	// Setup basic service
	svc := setupTestService(t)
	defer svc.cleanup()

	tests := []struct {
		name    string
		userId  string
		limit   int32
		wantErr bool
	}{
		{
			name:    "valid user with default limit",
			userId:  "1",
			limit:   0, // default limit
			wantErr: false,
		},
		{
			name:    "valid user with specific limit",
			userId:  "1",
			limit:   2,
			wantErr: false,
		},
		{
			name:    "user with fewer activities",
			userId:  "4",
			limit:   0, // default limit
			wantErr: false,
		},
		{
			name:    "nonexistent user",
			userId:  "999",
			limit:   10,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.QueryUserActivityRequest{
				UserId: tt.userId,
				Limit:  &wrapperspb.Int32Value{Value: tt.limit},
			}

			resp, err := svc.usersClient.QueryUserActivity(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Get expected activities from mock data and verify content
			expectedActivities := userActivityMap[tt.userId]
			expectedLen := len(expectedActivities)

			// If limit is specified and > 0, use the minimum of limit and expected length
			if tt.limit > 0 && int(tt.limit) < expectedLen {
				expectedLen = int(tt.limit)
			}

			verifyActivityContent(t, expectedActivities, resp.UserActivity, expectedLen)
		})
	}
}

func TestMutationCreatePost(t *testing.T) {
	// Setup basic service
	svc := setupTestService(t)
	defer svc.cleanup()

	// Store original user activity count to verify the post was added
	originalActivityCount := len(userActivityMap["1"])

	tests := []struct {
		name    string
		input   *service.PostInput
		wantErr bool
	}{
		{
			name: "valid post creation",
			input: &service.PostInput{
				Title:    "Test Post",
				AuthorId: "1",
			},
			wantErr: false,
		},
		{
			name: "invalid author",
			input: &service.PostInput{
				Title:    "Post by nonexistent user",
				AuthorId: "999",
			},
			wantErr: true,
		},
		{
			name: "another valid post",
			input: &service.PostInput{
				Title:    "Another Test Post",
				AuthorId: "2",
			},
			wantErr: false,
		},
	}

	validPostsCreated := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &service.MutationCreatePostRequest{
				Input: tt.input,
			}

			resp, err := svc.usersClient.MutationCreatePost(context.Background(), req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp.CreatePost)
			assert.Equal(t, tt.input.Title, resp.CreatePost.Title)
			assert.Equal(t, tt.input.AuthorId, resp.CreatePost.AuthorId)
			assert.NotEmpty(t, resp.CreatePost.Id)

			// Verify the post was added to the author's recent activity
			userResp, err := svc.usersClient.QueryUser(context.Background(), &service.QueryUserRequest{Id: tt.input.AuthorId})
			assert.NoError(t, err)
			assert.NotNil(t, userResp.User)

			// The new post should be the first item in recent activity (most recent first)
			assert.True(t, len(userResp.User.RecentActivity) > 0)
			firstActivity := userResp.User.RecentActivity[0]
			assert.NotNil(t, firstActivity.GetPost())
			assert.Equal(t, resp.CreatePost.Id, firstActivity.GetPost().Id)
			assert.Equal(t, resp.CreatePost.Title, firstActivity.GetPost().Title)

			validPostsCreated++
		})
	}

	// Verify user 1's activity count increased by 1 (only 1 valid post for user 1)
	newActivityCount := len(userActivityMap["1"])
	assert.Equal(t, originalActivityCount+1, newActivityCount)
}

func TestCreatePostIntegration(t *testing.T) {
	// Setup basic service
	svc := setupTestService(t)
	defer svc.cleanup()

	userID := "2" // Use Bob as our test user
	postTitle := "Integration Test Post"

	// Step 1: Query user to get initial state
	initialUserResp, err := svc.usersClient.QueryUser(context.Background(), &service.QueryUserRequest{Id: userID})
	assert.NoError(t, err)
	assert.NotNil(t, initialUserResp.User)

	initialActivityCount := len(initialUserResp.User.RecentActivity)
	t.Logf("User %s initial activity count: %d", userID, initialActivityCount)

	// Step 2: Create a new post for this user
	createPostReq := &service.MutationCreatePostRequest{
		Input: &service.PostInput{
			Title:    postTitle,
			AuthorId: userID,
		},
	}

	createPostResp, err := svc.usersClient.MutationCreatePost(context.Background(), createPostReq)
	assert.NoError(t, err)
	assert.NotNil(t, createPostResp.CreatePost)
	assert.Equal(t, postTitle, createPostResp.CreatePost.Title)
	assert.Equal(t, userID, createPostResp.CreatePost.AuthorId)

	createdPostID := createPostResp.CreatePost.Id
	t.Logf("Created post with ID: %s, Title: %s", createdPostID, postTitle)

	// Step 3: Query user again to verify the post was added to their activity
	updatedUserResp, err := svc.usersClient.QueryUser(context.Background(), &service.QueryUserRequest{Id: userID})
	assert.NoError(t, err)
	assert.NotNil(t, updatedUserResp.User)

	// Verify activity count increased by 1
	updatedActivityCount := len(updatedUserResp.User.RecentActivity)
	assert.Equal(t, initialActivityCount+1, updatedActivityCount)
	t.Logf("User %s updated activity count: %d", userID, updatedActivityCount)

	// Verify the new post is the first item in recent activity (most recent first)
	assert.True(t, len(updatedUserResp.User.RecentActivity) > 0)
	firstActivity := updatedUserResp.User.RecentActivity[0]

	// Verify it's a Post activity (not a Comment)
	assert.NotNil(t, firstActivity.GetPost())
	assert.Nil(t, firstActivity.GetComment())

	// Verify the post content matches what we created
	assert.Equal(t, createdPostID, firstActivity.GetPost().Id)
	assert.Equal(t, postTitle, firstActivity.GetPost().Title)
	assert.Equal(t, userID, firstActivity.GetPost().AuthorId)

	t.Logf("Verified new post appears first in user's recent activity")

	// Step 4: Query user activity directly to double-check
	activityResp, err := svc.usersClient.QueryUserActivity(context.Background(), &service.QueryUserActivityRequest{
		UserId: userID,
		Limit:  &wrapperspb.Int32Value{Value: 1}, // Just get the most recent activity
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(activityResp.UserActivity))

	// Verify this also shows our new post
	mostRecentActivity := activityResp.UserActivity[0]
	assert.NotNil(t, mostRecentActivity.GetPost())
	assert.Equal(t, createdPostID, mostRecentActivity.GetPost().Id)
	assert.Equal(t, postTitle, mostRecentActivity.GetPost().Title)

	t.Logf("Verified new post also appears via QueryUserActivity endpoint")
}
