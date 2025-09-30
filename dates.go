package pgorm

import (
	"time"
)

type Date struct {
}

type FormatString string

const (
	PERU                FormatString = "02/01/2006"
	YYYY_MM_DD_HH_MM_SS FormatString = "2006-01-02 15:04:05"
)

/**
 * Retorna la fecha en tipo time
 * @param {string} date: fecha en foprmato DD/MM/YYYY
 * @return {time.Time} fecha en tipo time
 */
func (d Date) GetDate(date string, format FormatString) time.Time {
	t, _ := time.Parse(string(format), date)
	return t
}

/**
 * Retorna Fecha en formato yyyy-mm-dd hh:mm:ss en zona horaria  (America/Bogota)
 * Return [string] : fecha  formato string yyyy-mm-dd hh:mm:ss (America/Bogota)
 */
func (d Date) GetDateLocationString() string {
	loc, _ := time.LoadLocation("America/Bogota")
	return time.Now().In(loc).Format(string(YYYY_MM_DD_HH_MM_SS))
}

func (d Date) GetDateLocation() time.Time {
	loc, _ := time.LoadLocation("America/Bogota")
	return time.Now().In(loc)
}

// retorna la fecha segun la locacion del usurio es este caso esta configurado para (America/Bogota) 2021-01-01 12:00:00.000
func (d Date) GetYear() int64 {
	loc, _ := time.LoadLocation("America/Bogota")
	t := time.Now().In(loc)
	return int64(t.Year())
}

func (d Date) GetMonth() int64 {
	loc, _ := time.LoadLocation("America/Bogota")
	t := time.Now().In(loc)
	return int64(t.Month())
}

// suma días a una fecha
func (d Date) Add(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}
