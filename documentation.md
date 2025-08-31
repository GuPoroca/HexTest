 <!-- markdownlint-disable -->
Basic model for the JSON input:
   ```{
      "Name": "Project Name",
      "Url": "http://urlOfTheRESTAPI.com",
      "Auth": {
        "Auth_type": "Auth Type",
        "Token": "Auth Token(In the future will automatically auth for you)"
      },
      "Parallel": false, (if you want to run all project suites in parallel)
      "Suites": [
        {
          "Name": "Suite Name",
          "Parallel": false, (if you want to run all suite tests in parallel)
          "Tests": [
            {
              "Name": "Name of the Test 1",
              "Method": "HTTP Method (GET, POST, PUT, DELETE)",
              "Request_body": "{ \"RequestBodyExample\": \"The request body, if in Json care about the quotes, body1\" }",
              "Asserts": [
              {
                "Response Body": {
                  "==": "{\"message\": \"GET called\"}",
                  "notNull": "",
                  "containsKey": [
					"message",
					"test",
					"other"
					]
                },
                
                "Response Size": {
                  "<=": "1000"
                }
              }
              ]
              "Api_endpoint": "/apiEndpointOfTheTest1",
              "Request_Headers": {
                "Name_Of_Header_1": "content_of_header_1",
                "Name_Of_Header_2": "content_of_header_2"
              },
              "Comment": "Comment on this test, not obrigatory"
            },
            {
              "Name": "Name of the Test 2",
              "Method": "HTTP Method (GET, POST, PUT, DELETE)",
              "Request_body": "{ \"RequestBodyExample\": \"The request body, if in Json care about the quotes, body2\" }",
              "Asserts": [
                {
                  "Response Status": {
                    "==": "200 OK"
                  },
                  "Response Time": {
                    "<=": "900"
                  },
                }
              ]
              "Api_endpoint": "/apiEndpointOfTheTest2",
              "Request_Headers": {
                "Name_Of_Header_3": "content_of_header_3",
                "Name_Of_Header_4": "content_of_header_4"
              },
              "Comment": "Comment on this test, not obrigatory"
            }
          ]
        }
      ]
    }
```

Right now the assertable fields are: 

  Response Body
  Response Status
  Response Time
  Response Size

And the assertions that are possible are:

  "==": Returns true if the value equals the assertion,
  "!=": Returns true if the value is not equal to the assertion,
  ">=": Returns true if the value is greater or equal than the assertion,
  "<=": Returns true if the value is less or equal than the assertion,
  "isNull": Returns true if the value is null,
  "notNull": Returns true if the value is not null,
  "containsKey": Used for Response Body. Returns true if the Json contains the key.
  "containsKey -R": Used for Response Body. Returns true if the Json contains the key (Iterate through all the Json Recursevely).


*TIP* you may assert multiple values at the same time example:
    "containsKey":[
        "teste",
        "val 2",
        "val 3"
    ]
