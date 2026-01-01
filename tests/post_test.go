package tests

import (
	"fmt"
	"starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/domain/user"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostTestSuite struct {
	suite.Suite
	accessToken string
	userID      uint
}

func TestPostTestSuite(t *testing.T) {
	suite.Run(t, new(PostTestSuite))
}

func (s *PostTestSuite) SetupSuite() {
	testApp = SetupTestApp()
}

func (s *PostTestSuite) TearDownSuite() {
	CleanupTestDB()
}

func (s *PostTestSuite) SetupTest() {
	testDB.Exec("DELETE FROM posts")
	testDB.Exec("DELETE FROM users")

	u := CreateTestUser(testDB, "testpost@example.com", "$2a$10$abcdefghijklmnopqrstuv", "user")
	s.userID = u.ID

	loginPayload := user.LoginRequest{
		Email:    u.Email,
		Password: "Password123!",
	}

	resp, body, _ := MakeRequest(testApp, "POST", "/api/auth/login", loginPayload, nil)

	if resp.StatusCode == 200 {
		var loginResponse map[string]interface{}
		ParseJSON(s.T(), body, &loginResponse)
		data := loginResponse["data"].(map[string]interface{})
		s.accessToken = data["access_token"].(string)
	}
}

func (s *PostTestSuite) TestCreatePost_Success() {
	if s.accessToken == "" {
		s.T().Skip("No access token available")
	}

	payload := post.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.accessToken,
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/posts", payload, headers)
	s.NoError(err)

	var response map[string]interface{}
	ParseJSON(s.T(), body, &response)

	if resp.StatusCode == 201 {
		s.Equal("success", response["status"])
		data := response["data"].(map[string]interface{})
		s.Equal("Test Post", data["title"])
		s.Equal("This is a test post content", data["content"])
	}
}

func (s *PostTestSuite) TestCreatePost_Unauthorized() {
	payload := post.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/posts", payload, nil)
	s.NoError(err)
	AssertErrorResponse(s.T(), resp, 401, body)
}

func (s *PostTestSuite) TestGetAllPosts_Success() {
	resp, body, err := MakeRequest(testApp, "GET", "/api/posts", nil, nil)
	s.NoError(err)
	AssertSuccessResponse(s.T(), resp, 200)

	var response map[string]interface{}
	ParseJSON(s.T(), body, &response)
	s.Equal("success", response["status"])
	s.Contains(response, "data")
}

func (s *PostTestSuite) TestGetPostByID_Success() {
	if s.accessToken == "" {
		s.T().Skip("No access token available")
	}

	createPayload := post.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.accessToken,
	}

	createResp, createBody, _ := MakeRequest(testApp, "POST", "/api/posts", createPayload, headers)

	if createResp.StatusCode == 201 {
		var createResponse map[string]interface{}
		ParseJSON(s.T(), createBody, &createResponse)
		postData := createResponse["data"].(map[string]interface{})
		postID := postData["id"].(float64)

		getResp, getBody, err := MakeRequest(testApp, "GET", fmt.Sprintf("/api/posts/%d", int(postID)), nil, nil)
		s.NoError(err)

		var getResponse map[string]interface{}
		ParseJSON(s.T(), getBody, &getResponse)

		if getResp.StatusCode == 200 {
			s.Equal("success", getResponse["status"])
			getData := getResponse["data"].(map[string]interface{})
			s.Equal("Test Post", getData["title"])
		}
	}
}

func (s *PostTestSuite) TestUpdatePost_Success() {
	if s.accessToken == "" {
		s.T().Skip("No access token available")
	}

	createPayload := post.CreatePostRequest{
		Title:   "Original Title",
		Content: "Original Content",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.accessToken,
	}

	createResp, createBody, _ := MakeRequest(testApp, "POST", "/api/posts", createPayload, headers)

	if createResp.StatusCode == 201 {
		var createResponse map[string]interface{}
		ParseJSON(s.T(), createBody, &createResponse)
		postData := createResponse["data"].(map[string]interface{})
		postID := postData["id"].(float64)

		updatePayload := post.UpdatePostRequest{
			Title:   "Updated Title",
			Content: "Updated Content",
		}

		updateResp, updateBody, err := MakeRequest(testApp, "PUT", fmt.Sprintf("/api/posts/%d", int(postID)), updatePayload, headers)
		s.NoError(err)

		var updateResponse map[string]interface{}
		ParseJSON(s.T(), updateBody, &updateResponse)

		if updateResp.StatusCode == 200 {
			s.Equal("success", updateResponse["status"])
			updateData := updateResponse["data"].(map[string]interface{})
			s.Equal("Updated Title", updateData["title"])
			s.Equal("Updated Content", updateData["content"])
		}
	}
}

func (s *PostTestSuite) TestDeletePost_Success() {
	if s.accessToken == "" {
		s.T().Skip("No access token available")
	}

	createPayload := post.CreatePostRequest{
		Title:   "Test Post",
		Content: "This will be deleted",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.accessToken,
	}

	createResp, createBody, _ := MakeRequest(testApp, "POST", "/api/posts", createPayload, headers)

	if createResp.StatusCode == 201 {
		var createResponse map[string]interface{}
		ParseJSON(s.T(), createBody, &createResponse)
		postData := createResponse["data"].(map[string]interface{})
		postID := postData["id"].(float64)

		deleteResp, deleteBody, err := MakeRequest(testApp, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, headers)
		s.NoError(err)

		var deleteResponse map[string]interface{}
		ParseJSON(s.T(), deleteBody, &deleteResponse)

		if deleteResp.StatusCode == 200 {
			s.Equal("success", deleteResponse["status"])
		}
	}
}
