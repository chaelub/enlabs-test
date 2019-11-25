package transaction

const (
	lastNByStatus = `SELECT id from transaction where state=%d order by tms desc limit 10`
)
