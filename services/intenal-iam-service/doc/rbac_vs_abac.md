## RBAC vs ABAC: A Simple Comparison

| **Criteria**      | **RBAC**                                | **ABAC**                                          |
|-------------------|-----------------------------------------|--------------------------------------------------|
| **Access Control** | Based on predefined roles               | Based on attributes (user, resource, context)      |
| **Flexibility**    | Limited flexibility                     | Highly flexible with dynamic policies              |
| **Granularity**    | Coarse-grained access                   | Fine-grained, contextual access                    |
| **Complexity**     | Simpler to implement but grows complex with role proliferation | More complex to set up but scales better in large systems |
| **Examples**       | Admin, Manager, Employee roles          | User from HR department accessing specific data during work hours |