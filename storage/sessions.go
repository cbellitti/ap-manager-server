package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"gitlab.com/astronomy/ap-manager/helpers"
)

type Session struct {
	ID                   int    `json:"id"`
	TelescopeFocalLength int    `json:"focalLength"`
	Mount                string `json:"mount"`
	Focuser              bool   `json:"focuser"`
	Guider               bool   `json:"guider"`
	Rotator              bool   `json:"rotator"`
	Camera               string `json:"camera"`
	Target               string `json:"target"`
	SubExposureTime      string `json:"subexposureTime"`
	SubExposureQty       int    `json:"subExposureQty"`
	Filter               string `json:"filter"`
	Gain                 string `json:"gain"`
	Offset               string `json:"offset"`
	Binning              string `json:"binning"`
	TimeStart            string `json:"timeStart"`
	TimeEnd              string `json:"timeEnd"`
}

func CreateSessions(db *sql.DB, s helpers.Session) []int64 {
	//This is a test

	focuser := "FALSE"
	if s.Equipment.Focuser {
		focuser = "TRUE"
	}

	guider := "FALSE"
	if s.Equipment.Guider {
		guider = "TRUE"
	}

	rotator := "FALSE"
	if s.Equipment.Rotator {
		rotator = "TRUE"
	}

	var ids []int64

	for _, t := range s.Targets {
		sql := "INSERT INTO sessions(starttime, endtime, focuser, guider, rotator, camera, mount, focallength, target,subexposuretime, subexposureqty, filter,gain, offset, binning) VALUES ('" + t.TimeStampStart + "','" + t.TimeStampEnd + "'," + focuser + "," + guider + "," + rotator + ",'" + s.Equipment.Camera + "','" + s.Equipment.Mount + "'," + strconv.Itoa(s.Equipment.TelescopeFocalLength) + ",'" + t.Name + "','" + t.SubExposureTime + "'," + strconv.Itoa(t.SubExposureQty) + ",'" + t.Filter + "','" + t.Gain + "','" + t.Offset + "','" + t.Binning + "')"

		fmt.Printf("%s", sql)

		res, err := db.Exec(sql)
		if err != nil {
			panic(err.Error())
		}

		lastId, err := res.LastInsertId()

		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, lastId)
	}

	return ids
}

func ReadSessions(db *sql.DB) []Session {
	// Execute the query fgfgf
	results, err := db.Query("SELECT * from sessions")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var ss []Session
	for results.Next() {
		var s Session

		// for each row, scan the result into our tag composite object
		err = results.Scan(&s.ID, &s.TimeStart, &s.TimeEnd, &s.Focuser, &s.Guider, &s.Rotator, &s.Camera, &s.Mount, &s.TelescopeFocalLength, &s.Target, &s.SubExposureTime, &s.SubExposureQty, &s.Filter, &s.Gain, &s.Offset, &s.Binning)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		ss = append(ss, s)
	}
	return ss
}
