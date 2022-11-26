#### Typescript: TStype

---

```typescript
// Typescript: TStype=  MyType = number
generate:
export type  MyType = number

// examples
// Typescript: type
type TestType []int
generate:
export type  TestType = number[]


// Typescript: type
type TestTypeMap map[string]map[int]string
generate:
export type  TestTypeMap = Record<string, <Record<number, string>>>


// Typescript: type=Date
type TestTypeTime time.Time
generate:
export type  TestTypeTime = Date

// typescript: type
type TestNullable *string
generate:
export type TestNullable = Nullable<string>
```
