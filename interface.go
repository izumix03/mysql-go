package mysqlgo

import "time"

// Basic means you tend to use orm framework.
// Struct should have name tag (ex. `name:"column_name_dummy"`)
type Basic interface {
	// From returns table name
	Table() string
}

// Insertable means struct can be used when insert.
// Struct should have name tag (ex. `name:"column_name_dummy"`)
// and should have insert tag when exclude when insert (ex. `name:"id" insert:"false"`)
type Insertable interface {
	// SetCreatedAt sets created at by orm
	SetCreatedAt(time time.Time)
}

// Updatable means struct can be used when update.
// and should have update tag when duplicate key exists (ex. `name:"updated_at" upsert:"true"`)
type Updatable interface {
	// SetUpdatedAt sets updated at by orm
	SetUpdatedAt(time time.Time)
}
