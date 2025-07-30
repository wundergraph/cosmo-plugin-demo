package main

import (
	service "github.com/wundergraph/cosmo/plugin/generated"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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

// Mock posts data
var mockPosts = map[string]*service.Post{
	"1": {Id: "1", Title: "Getting Started with GraphQL", AuthorId: "1"},
	"2": {Id: "2", Title: "Advanced Federation Patterns", AuthorId: "1"},
	"3": {Id: "3", Title: "Building Scalable APIs", AuthorId: "2"},
	"4": {Id: "4", Title: "TypeScript Best Practices", AuthorId: "3"},
}

// Mock comments data
var mockComments = map[string]*service.Comment{
	"1": {Id: "1", Content: "Great post! Very helpful.", AuthorId: "2"},
	"2": {Id: "2", Content: "Thanks for sharing this.", AuthorId: "3"},
	"3": {Id: "3", Content: "Looking forward to more content.", AuthorId: "4"},
	"4": {Id: "4", Content: "Excellent examples provided.", AuthorId: "1"},
}

// Mock user data store for demonstration purposes
// In a production environment, this would be replaced with a proper database
var mockUsers = map[string]*service.User{
	"1": {
		Id:          "1",
		Name:        "Alice Johnson",
		Email:       "alice@example.com",
		Role:        service.UserRole_USER_ROLE_ADMIN,
		Permissions: []string{"read", "write"},
		Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"admin", "user"}}},
		SkillCategories: &service.ListOfListOfString{
			List: &service.ListOfListOfString_List{
				Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"JavaScript", "TypeScript"}}},
					{List: &service.ListOfString_List{Items: []string{"React", "Vue", "Angular"}}},
					{List: &service.ListOfString_List{Items: []string{"Node.js", "Express"}}},
				},
			},
		},
		RecentActivity: []*service.ActivityItem{
			{Value: &service.ActivityItem_Post{Post: mockPosts["1"]}},
			{Value: &service.ActivityItem_Post{Post: mockPosts["2"]}},
			{Value: &service.ActivityItem_Comment{Comment: mockComments["4"]}},
		},
		Profile: &service.Profile{
			DisplayName: &wrapperspb.StringValue{Value: "Alice J."},
			Timezone:    &wrapperspb.StringValue{Value: "America/New_York"},
			Theme:       service.Theme_THEME_DARK,
		},
		Bio: &wrapperspb.StringValue{Value: "Full-stack developer with 5+ years of experience"},
		Age: &wrapperspb.Int32Value{Value: 28},
	},
	"2": {
		Id:          "2",
		Name:        "Bob Smith",
		Email:       "bob@example.com",
		Role:        service.UserRole_USER_ROLE_USER,
		Permissions: []string{"read"},
		Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"user"}}},
		SkillCategories: &service.ListOfListOfString{
			List: &service.ListOfListOfString_List{
				Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"Python", "Java"}}},
					{List: &service.ListOfString_List{Items: []string{"Django", "Spring"}}},
				},
			},
		},
		RecentActivity: []*service.ActivityItem{
			{Value: &service.ActivityItem_Post{Post: mockPosts["3"]}},
			{Value: &service.ActivityItem_Comment{Comment: mockComments["1"]}},
		},
		Profile: &service.Profile{
			DisplayName: &wrapperspb.StringValue{Value: "Bob"},
			Timezone:    &wrapperspb.StringValue{Value: "Europe/London"},
			Theme:       service.Theme_THEME_LIGHT,
		},
		Bio: &wrapperspb.StringValue{Value: "Backend developer passionate about clean code"},
		Age: &wrapperspb.Int32Value{Value: 32},
	},
	"3": {
		Id:          "3",
		Name:        "Charlie Brown",
		Email:       "charlie@example.com",
		Role:        service.UserRole_USER_ROLE_USER,
		Permissions: []string{"read"},
		Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"user"}}},
		SkillCategories: &service.ListOfListOfString{
			List: &service.ListOfListOfString_List{
				Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"Go", "Rust"}}},
					{List: &service.ListOfString_List{Items: []string{"Docker", "Kubernetes"}}},
				},
			},
		},
		RecentActivity: []*service.ActivityItem{
			{Value: &service.ActivityItem_Post{Post: mockPosts["4"]}},
			{Value: &service.ActivityItem_Comment{Comment: mockComments["2"]}},
		},
		Profile: &service.Profile{
			Timezone: &wrapperspb.StringValue{Value: "Asia/Tokyo"},
			Theme:    service.Theme_THEME_AUTO,
		},
		Age: &wrapperspb.Int32Value{Value: 29},
	},
	"4": {
		Id:          "4",
		Name:        "Dana Lee",
		Email:       "dana@example.com",
		Role:        service.UserRole_USER_ROLE_GUEST,
		Permissions: []string{"read"},
		Tags:        &service.ListOfString{List: &service.ListOfString_List{Items: []string{"guest"}}},
		SkillCategories: &service.ListOfListOfString{
			List: &service.ListOfListOfString_List{
				Items: []*service.ListOfString{
					{List: &service.ListOfString_List{Items: []string{"HTML", "CSS"}}},
				},
			},
		},
		RecentActivity: []*service.ActivityItem{
			{Value: &service.ActivityItem_Comment{Comment: mockComments["3"]}},
		},
		Profile: &service.Profile{
			DisplayName: &wrapperspb.StringValue{Value: "Dana"},
			Theme:       service.Theme_THEME_LIGHT,
		},
		Bio: &wrapperspb.StringValue{Value: "Learning web development"},
		Age: &wrapperspb.Int32Value{Value: 24},
	},
}

// User activity mappings for efficient lookup
var userActivityMap = map[string][]*service.ActivityItem{
	"1": {
		{Value: &service.ActivityItem_Post{Post: mockPosts["1"]}},
		{Value: &service.ActivityItem_Post{Post: mockPosts["2"]}},
		{Value: &service.ActivityItem_Comment{Comment: mockComments["4"]}},
	},
	"2": {
		{Value: &service.ActivityItem_Post{Post: mockPosts["3"]}},
		{Value: &service.ActivityItem_Comment{Comment: mockComments["1"]}},
	},
	"3": {
		{Value: &service.ActivityItem_Post{Post: mockPosts["4"]}},
		{Value: &service.ActivityItem_Comment{Comment: mockComments["2"]}},
	},
	"4": {
		{Value: &service.ActivityItem_Comment{Comment: mockComments["3"]}},
	},
}
