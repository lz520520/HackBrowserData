package creditcard

import (
	"database/sql"
	"log/slog"
	"os"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	"tests/crypto"
	"tests/item"
)

type ChromiumCreditCard []card

type card struct {
	GUID            string
	Name            string
	ExpirationYear  string
	ExpirationMonth string
	CardNumber      string
	Address         string
	NickName        string
}

const (
	queryChromiumCredit = `SELECT guid, name_on_card, expiration_month, expiration_year, card_number_encrypted, billing_address_id, nickname FROM credit_cards`
)

func (c *ChromiumCreditCard) Parse(masterKey []byte) error {
	db, err := sql.Open("sqlite3", item.ChromiumCreditCard.TempFilename())
	if err != nil {
		return err
	}
	defer os.Remove(item.ChromiumCreditCard.TempFilename())
	defer db.Close()

	rows, err := db.Query(queryChromiumCredit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, month, year, guid, address, nickname string
			value, encryptValue                        []byte
		)
		if err := rows.Scan(&guid, &name, &month, &year, &encryptValue, &address, &nickname); err != nil {
			slog.Error("scan chromium credit card error", "err", err)
		}
		ccInfo := card{
			GUID:            guid,
			Name:            name,
			ExpirationMonth: month,
			ExpirationYear:  year,
			Address:         address,
			NickName:        nickname,
		}
		if len(encryptValue) > 0 {
			if len(masterKey) == 0 {
				value, err = crypto.DPAPI(encryptValue)
			} else {
				value, err = crypto.DecryptPass(masterKey, encryptValue)
			}
			if err != nil {
				slog.Error("decrypt chromium credit card error", "err", err)
			}
		}

		ccInfo.CardNumber = string(value)
		*c = append(*c, ccInfo)
	}
	return nil
}

func (c *ChromiumCreditCard) Name() string {
	return "creditcard"
}

func (c *ChromiumCreditCard) Len() int {
	return len(*c)
}

type YandexCreditCard []card

func (c *YandexCreditCard) Parse(masterKey []byte) error {
	db, err := sql.Open("sqlite3", item.YandexCreditCard.TempFilename())
	if err != nil {
		return err
	}
	defer os.Remove(item.YandexCreditCard.TempFilename())
	defer db.Close()
	rows, err := db.Query(queryChromiumCredit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, month, year, guid, address, nickname string
			value, encryptValue                        []byte
		)
		if err := rows.Scan(&guid, &name, &month, &year, &encryptValue, &address, &nickname); err != nil {
			slog.Error("scan chromium credit card error", "err", err)
		}
		ccInfo := card{
			GUID:            guid,
			Name:            name,
			ExpirationMonth: month,
			ExpirationYear:  year,
			Address:         address,
			NickName:        nickname,
		}
		if len(encryptValue) > 0 {
			if len(masterKey) == 0 {
				value, err = crypto.DPAPI(encryptValue)
			} else {
				value, err = crypto.DecryptPass(masterKey, encryptValue)
			}
			if err != nil {
				slog.Error("decrypt chromium credit card error", "err", err)
			}
		}
		ccInfo.CardNumber = string(value)
		*c = append(*c, ccInfo)
	}
	return nil
}

func (c *YandexCreditCard) Name() string {
	return "creditcard"
}

func (c *YandexCreditCard) Len() int {
	return len(*c)
}
