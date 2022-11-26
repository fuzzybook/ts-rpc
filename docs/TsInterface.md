#### Typescript: interface

---

```typescript
// Typescript: interface
type HTTPResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}
generate:
export interface HTTPResponse {
	data: unknown;
	error: unknown;
}

// Typescript: interface
type FormRequest struct {
	Req   string `json:"req"`
	Count int    `json:"count"`
}
generate:
export interface FormRequest {
	req: string;
	count: number;
}

// Typescript: interface
type FormResponse struct {
	Test string `json:"test"`
}
generate:
export interface FormResponse {
	test: string;
}

// Typescript: interface
type FormResponse2 struct {
	Test string `json:"test"`
	User  string `json:"user,omitempty"`
}
generate:
export interface FormResponse2 {
	test: string;
	user?: string
}


// Typescript: interface
type FormResponse3 struct {
	FormRequest 	`ts:"expand"`
	Test string 	`json:"test"`
	User string 	`json:"user,omitempty"`
	Time time.Time 	`json:"time" ts:"type=Date"`
}
generate:
export interface FormResponse3 {
	req: string;
	count: number;
	test: string;
	user?: string;
	time: Date;
}

```
