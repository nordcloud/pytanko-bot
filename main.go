package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/machinebox/graphql"
)

// Response https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

//Build method created APIGatewayProxyResponse using provided status code and message or the body payload
func Respond(statusCode int, message string) Response {
	var buf bytes.Buffer

	var body []byte

	if message != "" {
		body, _ = json.Marshal(map[string]interface{}{
			"text": message,
		})
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: statusCode,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp
}

// ApiResponse defines basic AppSync return
type ApiResponse struct {
	Data   interface{} `json:"data"`
	Errors *string     `json:"errors"`
}

// Handler is the main entry for the Lambda function
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	params, err := url.ParseQuery(request.Body)

	if err != nil {
		log.Fatal("Could not parse the query")
	}

	if params.Get("text") == "" {
		log.Println("Question was not asked")
		return Respond(200, "Question was not provided ☹️"), nil
	}

	fmt.Printf("User %s asked: %s\n", params.Get("user_name"), params.Get("text"))
	client := graphql.NewClient(os.Getenv("API_URL"))

	now := time.Now()

	query := fmt.Sprintf(`
		mutation create {
			createQuestion(input:{content: "%s", date: "%s", username: "%s"}) {id}
		}
	`, params.Get("text"), now.Format(time.UnixDate), params.Get("user_name"))

	req := graphql.NewRequest(query)
	req.Header.Set("X-Api-Key", os.Getenv("API_KEY"))

	var respData ApiResponse
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.Errors != nil {
		log.Fatal("Received an error from the API", respData.Errors)
	}

	return Respond(200, "Question added 💪"), nil
}

func main() {
	lambda.Start(Handler)
}
