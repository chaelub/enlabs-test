package transaction

const (
	lastNByStatus      = `SELECT id FROM transactions WHERE state=%d ORDER BY tms DESC LIMIT 10`
	transactionByExtId = `SELECT * FROM transactions WHERE extid='%s'`
	checkExistsByExtId = `SELECT count(extid) from transactions where extid='%s'`
)
