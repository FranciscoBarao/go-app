package sql

// OrderBy is the ORDER BY clause fragment.
const OrderBy = ` ORDER BY %s %s`

// WhereLike is the WHERE ILIKE clause fragment for partial string matching.
const WhereLike = ` WHERE %s ILIKE $1`

// WhereOp is the WHERE clause fragment for comparison operators.
const WhereOp = ` WHERE %s %s $1`
