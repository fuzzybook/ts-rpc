#### Typescript: enum

---

```typescript
type Direction int
// Typescript: enum=Direction
const (
	North Direction = iota
	East
	South
	West
)
func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}
generate:
export const EnumDirection = {
North: 0,
East: 1,
South: 2,
West: 3,
} as const
export type Direction = typeof EnumDirection[keyof typeof EnumDirection]

type Season string
// Typescript: enum=Season
const (
	Summer Season = "summer"
	Autumn        = "autumn"
	Winter        = "winter"
	Spring        = "spring"
)
generate:
export const EnumSeason = {
Summer: "summer",
Autumn: "autumn",
Winter: "winter",
Spring: "spring",
} as const
export type Season = typeof EnumSeason[keyof typeof EnumSeason]

// Typescript: enum=Test
const (
	A int = iota
	B
	C
	D
)
generate:
export const EnumTest = {
A: 0,
B: 1,
C: 2,
D: 3,
}
export type Test = typeof EnumTest[keyof typeof EnumTest]
```
