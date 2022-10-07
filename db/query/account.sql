-- generate the code using sqlc
-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES  (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
-- for update是给查询的那一行加锁，直到该tx commit前，其他tx无法update, delete 和select for update，会进入等待状态
-- 如果没有forupdate，两个goroutine在update account前都运行到get account1的话，那么查询出来的余额都是未改变的
-- 如果不加no key，因为transfer里有foreign key关联到accounts，在两个go routine查询accounts前都insert transfer的话，会造成死锁
-- 所以在其他表有外键关联到主键进行改动时，主键那张表的查询会等待。
-- 或者主键那张表的查询时，其他表有外键关联到主键进行改动时会等待。
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE; 

-- name: ListAccounts :many
SELECT * from accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
