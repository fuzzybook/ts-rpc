#### Typescript: TSDeclaration

---

```typescript
// Typescript: TSDeclaration= Nullable<T> = T | null;
// Typescript: TSDeclaration= Record<K extends string | number | symbol, T> = { [P in K]: T; }
// Typescript: TSDeclaration= MySecialArray<T> = T[];
generate: export type Nullable<T> = T | null;
export type Record<K extends string | number | symbol, T> = { [P in K]: T };
export type MySecialArray<T> = T[];
```
