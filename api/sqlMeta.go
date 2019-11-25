package api

// 存储过程名称
var (
	billList   = "Mob_GetHisBill"
	billDetail = "Mob_GetHisBillDetail"
)

var mssqlQueryAccountCmd = "SELECT C_hh, C_dz, C_hm, C_sj FROM YHDA WHERE C_hh = @account"

var mssqlQueryCheckAccountCmd = "SELECT C_hh FROM YHDA WHERE C_hh = @account AND C_hm = @name"

var mysqlQueryAccountCmd = "SELECT account FROM user_data WHERE id_card_number=?"

var mysqlQueryAccountCheckCmd = "SELECT account FROM user_data WHERE id_card_number=? AND account=?"

var mysqlQueryDefaultAccountCmd = "SELECT default_account FROM wechat_profile WHERE id_card_number=?"

var mysqlQueryAccountDataCmd = "SELECT * FROM account_data WHERE account=?"

var mysqlInsertAccountCmd = `
INSERT INTO account_data
(account, address, name, phone, charge, current_meter, meter, paid, unpaid_count)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?)
`

var mysqlUpdateAccountCmd = `
UPDATE account_data SET
phone=?,               
charge=?,
current_meter=?,
meter=?,
paid=?,
unpaid_count=?
WHERE account=?
`

var mysqlInsertFeeDetailCmd = `
INSERT INTO fee_detail
(yszbh, account, name, current_period, charge, current_meter, previous_meter, paid, wsf, xfft, ljf, ecjydf, szyf, cjhys, wyj, wswyj, water_charge)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

var mysqlUpdateFeeDetailCmd = `
UPDATE fee_detail SET
paid = ?,
charge = ?
WHERE yszbh = ?
`

var mysqlQueryFeeDetailCmd = "SELECT * FROM fee_detail WHERE account = ? ORDER BY id DESC LIMIT 1"