package integration_test

// // import "context"

// // // create a client (safe to share across requests)
// // client := graphql.NewClient("https://machinebox.io/graphql")

// // // make a request
// // req := graphql.NewRequest(`
// //     query ($key: String!) {
// //         items (id:$key) {
// //             field1
// //             field2
// //             field3
// //         }
// //     }
// // `)

// // // set any variables
// // req.Var("key", "value")

// // // set header fields
// // req.Header.Set("Cache-Control", "no-cache")

// // // define a Context for the request
// // ctx := context.Background()

// // // run it and capture the response
// // var respData ResponseStruct
// // if err := client.Run(ctx, req, &respData); err != nil {
// //     log.Fatal(err)
// // }

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws/session"
// 	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
// )

// func TestAutoPass(t *testing.T) {
// 	client := new(http.Client)
// 	query := map[string]string{
// 		"query": `
//             {
//                 people {
// 					id,
//                     name,
//                     age,
//                 }
//             }
//         `,
// 	}
// 	b, err := json.Marshal(&query)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// construct the request object
// 	req, err := http.NewRequest("POST", "https://63yrawximfckvc3l6qjcpjgmte.appsync-api.us-east-1.amazonaws.com/graphql", bytes.NewReader(b))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// get aws credential
// 	sess := session.Must(session.NewSession())

// 	//sign the request
// 	signer := v4.NewSigner(sess.Config.Credentials)
// 	signer.Sign(req, bytes.NewReader(b), "appsync", "us-east-1", time.Now())

// 	//FIRE!!
// 	response, err := client.Do(req)
// 	fmt.Println(err)

// 	//print the response
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(response.Body)
// 	newStr := buf.String()

// 	fmt.Println(newStr)
// }
