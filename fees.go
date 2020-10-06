package main

type Fees struct{
	Id int
	Lvl string
	Fee int
}

// Display Fees
func FeesDisp() []Fees{
	// Start DB Connection
	db := DbConn()
	fee_row, err := db.Query("SELECT * FROM fees")
	CheckErr(err)
	fee := Fees{}
	fees := []Fees{}
	for fee_row.Next(){
		err = fee_row.Scan(&fee.Id, &fee.Lvl, &fee.Fee) 
		CheckErr(err)
		fees = append(fees, fee)
	}
	defer db.Close()
	return fees
}

func GetFee(course string) Fees{
	// Start DB Connection
	db := DbConn()
	var Fee Fees
	row, err := db.Query(`
		SELECT fee_id, 
		lvl, 
		fee
		FROM edc_db.fees WHERE lvl=?`, course)
		CheckErr(err)
	for row.Next(){
		err = row.Scan(&Fee.Id, &Fee.Lvl, &Fee.Fee)
		CheckErr(err)
	}
	defer db.Close()
	return Fee

}

func EditFees(fees Fees){
	// Start DB Connection
	db := DbConn()
	row, err := db.Query("UPDATE `edc_db`.`fees` SET `fee` = ? WHERE `fee_id` = ?",fees.Fee, fees.Id)
	CheckErr(err)
	row.Close()
	defer db.Close()
}