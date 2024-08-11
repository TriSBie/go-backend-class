# How to prevent deadlock in golang 
-> What is deadlock caused by ?
Supposed 2 transactions occurs concurrently in same phrase 
=> Transaction A is trigger update queries, same time transaction B is also triggers a transaction updates -> Transaction B will waiting the update results from transaction A in periods of waiting time till the operation response.
<-> Solutions :
- Using update query by order

# SELECT ... FOR UPDATE (1)
- Locking mechanisms used within transaction
- Prevent race conditions and ensure data consistency.
# SELECT ... FOR NO KEY UPDATE (2)
- equivalence as (1) query -> but weaker
- selected rows are not modified in ways that affect the non-key columns
============================================================================
# Why context in golang ?
- Carry deadlines, Cancellation signals, request-scoped data.
- Helps to propagating cancellation signals throughout the application.
- Initialize context by using:
  + context.WithCancel(parent) -> cancel by using cancel()
  + context.WithDeadline(parent, deadline) -> automatically canceled when `deadline` is reached.
  + context.WithTimeOut(parent, timeout) -> similarly to withDeadline
  + context.WithValue(parent, key, value) -> [often use with the request-scoped]
============================================================================
# Dives deeper into Transaction
Command showing the current isolation levels :
    - `show transaction isolation level`
Command options for setting isolation levels :
    - `set session transaction isolation level serializable`
1. # Read Phenomena
 + `Dirty Read`: transaction reads data written by other concurrent uncommitted transaction
 + `Non-repeated read`: reads the same row twice and sees difference value
 + `Phantom read` : transaction re-executes a query to find rows satisfy the condition and sees the difference set of rows
 + `Serialization anomaly` : 
2. # 4 standards isolation levels
 + Read uncommitted -> read committed -> Repeatable read -> Serializable read