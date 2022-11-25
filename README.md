# ts-rpc

#### typescript RPC

---

#####The ghost bridge from front-end to backend

this project is a POC - Proof of concept
the goal is to generate typescript code from golang code using AST
using intellisense in the development IDE

```typescript
const input = <server.FormRequest>{ req: "some request", count: 456 };
const { data, error } = await server.postTest(input);
if (!error) {
  console.log(data);
  result = data.test;
}
```

golang declaration

```golang
// Typescript: TSDeclaration= Nullable<T> = T | null;
// Typescript: TSDeclaration= Record<K extends string | number | symbol, T> = { [P in K]: T; }

// Typescript: interface
type HTTPResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

// Typescript: interface
type FormRequest struct {
	Req   string `json:"req"`
	Count int    `json:"count"`
}

// Typescript: interface
type FormResponse struct {
	Test string `json:"test"`
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
```

generated code

```typescript
//
// namespace server
//

export namespace server {
  export interface FormRequest {
    req: string;
    count: number;
  }

  export interface FormResponse {
    test: string;
  }

  export interface HTTPResponse {
    data: unknown;
    error: unknown;
  }

  // Typescript: TSEndpoint= path=/postTest;  name=postTest; method=POST; request=FormRequest; response=FormResponse
  // server/server.go Line: 73
  export const postTest = async (data: FormRequest): Promise<{ data: FormResponse; error: Nullable<string> }> => {
    return (await api.POST("/postTest", data)) as { data: FormResponse; error: Nullable<string> };
  };
}
```

[fetch code](TsFetch.md)
