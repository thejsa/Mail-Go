package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func Account(w http.ResponseWriter, r *http.Request) {
	wiiID := r.Form.Get("mlid")
	if !friendCodeIsValid(wiiID) {
		fmt.Fprint(w, GenNormalErrorCode(610, "Invalid Wii Friend Code."))
		return
	}

	w.Header().Add("Content-Type", "text/plain;charset=utf-8")

	stmt, err := db.Prepare("INSERT IGNORE INTO `accounts` (`mlid`,`passwd`, `mlchkid` ) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Fprint(w, GenNormalErrorCode(410, "Database error."))
		log.Fatal(err)
		return
	}

	passwd := RandStringBytesMaskImprSrc(16)
	passwdByte := sha512.Sum512(append(salt, []byte(passwd)...))
	passwdHash := hex.EncodeToString(passwdByte[:])

	mlchkid := RandStringBytesMaskImprSrc(32)
	mlchkidByte := sha512.Sum512(append(salt, []byte(mlchkid)...))
	mlchkidHash := hex.EncodeToString(mlchkidByte[:])

	result, err := stmt.Exec(wiiID, passwdHash, mlchkidHash)
	if err != nil {
		fmt.Fprint(w, GenNormalErrorCode(410, "Database error."))
		log.Println(err)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Fprint(w, GenNormalErrorCode(410, "Database error."))
		log.Println(err)
		return
	}

	if affected == 0 {
		fmt.Fprint(w, "\n",
			GenNormalErrorCode(211, "Duplicate registration."))
		return
	}

	fmt.Fprint(w, "\n",
		GenNormalErrorCode(100, "Success."),
		"mlid=", wiiID, "\n",
		"passwd=", passwd, "\n",
		"mlchkid=", mlchkid, "\n")
}
