package sql

// OrderBy is the ORDER BY clause fragment.
const OrderBy = ` ORDER BY %s %s`

// OrderByIDFallback is a deterministic ORDER BY used when no sort is requested,
// so paginated results are stable across pages.
const OrderByIDFallback = ` ORDER BY id ASC`

// OrderByIDTiebreak appends id as a secondary sort key to keep pagination
// stable when the primary sort column has duplicate values.
const OrderByIDTiebreak = `, id ASC`

// LimitOffset is the pagination clause fragment; args bind LIMIT then OFFSET.
const LimitOffset = ` LIMIT $%d OFFSET $%d`

// WhereLike is the WHERE ILIKE clause fragment for partial string matching.
const WhereLike = ` WHERE %s ILIKE $1`

// WhereOp is the WHERE clause fragment for comparison operators.
const WhereOp = ` WHERE %s %s $1`
