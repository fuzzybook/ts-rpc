#### Struct field tags

---

```typescript

tag `json:"fieldName"` this tag is manatory
tag `json:"-"` the field is not exported
tag `json:"fieldName,omitempty"` generate fieldName?: ("the props can be undefined")

tag `ts:"expand"` generate the golang composition in the interface
tag `ts:"type=Date"` force to use the type Date
tag `ts:"type=WhatEverYouWant"` force to use the type WhatEverYouWant
```
