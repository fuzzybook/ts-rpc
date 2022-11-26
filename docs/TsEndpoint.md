#### Typescript: TSEndpoint

---

```typescript

// Typescript: TSEndpoint= path=/ping; name=ping; method=GET; response=string
r.GET("/ping", func(c *gin.Context) {
    response := HTTPResponse{Data: "pong", Error: nil}
    c.JSON(http.StatusOK, response)
})
generate:
// Typescript: TSEndpoint= path=/ping; name=ping; method=GET; response=string
// server/server.go Line: 53
export const ping = async ():Promise<{ data:string; error: Nullable<string> }> => {
	return await api.GET("/ping") as { data: string; error: Nullable<string> };
}

// Typescript: TSEndpoint= path=/postTest;  name=postTest; method=POST; request=FormRequest; response=FormResponse
r.POST("/postTest", func(c *gin.Context) {
    var requestBody FormRequest
    if err := c.BindJSON(&requestBody); err != nil {
        response := HTTPResponse{Data: nil, Error: "wrongData"}
        c.JSON(http.StatusOK, response)
    }
    response := HTTPResponse{Data: FormResponse{Test: fmt.Sprintf("%d", requestBody.Count)}, Error: nil}
    c.JSON(http.StatusOK, response)
})

generate:

// Typescript: TSEndpoint= path=/postTest;  name=postTest; method=POST; request=FormRequest; response=FormResponse
// server/server.go Line: 73
export const postTest = async (data: FormRequest):Promise<{ data:FormResponse; error: Nullable<string> }> => {
	return await api.POST("/postTest", data) as { data: FormResponse; error: Nullable<string> };
}
```
