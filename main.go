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

// Response alias
type Response events.APIGatewayProxyResponse

type ApiResponse struct {
	Data   interface{} `json:"data"`
	Errors *string     `json:"errors"`
}

// Handler is the main entry for the Lambda function
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	var buf bytes.Buffer

	params, err := url.ParseQuery(request.Body)

	if err != nil {
		log.Fatal("Could not parse the query")
	}

	fmt.Println(params.Get("text"))
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
		log.Fatal("Wyjebane", respData.Errors)
	}

	body, _ := json.Marshal(map[string]string{
		"text": "Question added 💪",
	})

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: 200,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
