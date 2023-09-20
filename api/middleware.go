package api

import (
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var googleOauthConfig = oauth2.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURL:  "your-redirect-url",
    Scopes:       []string{"profile", "email"},
    Endpoint:     google.Endpoint,
}

// UserInfo represents the structure of user information you want to retrieve.
type UserInfo struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Name     string `json:"name"`
    Picture  string `json:"picture"`
    // Add other fields as needed
}

func OAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract the OAuth token from the request (e.g., from cookies or session)
        token, err := c.Cookie("oauth_token") // Example: Store the token in a cookie

        if err != nil {
            // If the token is missing, return an unauthorized response
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort() // Abort the request processing
            return
        }

        // Validate the OAuth token using Google's OAuth2 package
        googleToken, err := googleOauthConfig.TokenSource(c, &oauth2.Token{AccessToken: token}).Token()
        if err != nil {
            // If the token is invalid, return an unauthorized response
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OAuth token"})
            c.Abort() // Abort the request processing
            return
        }

        // Retrieve user information using the OAuth token
        userInfo, err := FetchUserInfo(googleToken.AccessToken)
        if err != nil {
            // Handle the error when fetching user info (e.g., log it)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user info"})
            c.Abort() // Abort the request processing
            return
        }

        // Set the user information in the Gin context
        c.Set("user", userInfo)

        // Continue processing the request
        c.Next()
    }
}

// FetchUserInfo retrieves user information from the Google People API using the access token.
func FetchUserInfo(accessToken string) (*UserInfo, error) {
    // Make an HTTP GET request to the Google People API
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://people.googleapis.com/v1/people/me?personFields=emailAddresses,names,photos", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+accessToken)
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, err
    }

    // Parse the JSON response into the UserInfo struct
    var userInfo UserInfo
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&userInfo); err != nil {
        return nil, err
    }

    return &userInfo, nil
}
